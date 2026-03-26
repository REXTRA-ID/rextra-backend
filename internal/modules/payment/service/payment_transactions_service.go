package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	membershipRepo "rextra-backend/internal/modules/membership/repository"
	"rextra-backend/internal/modules/payment/repository"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/tripay"
	"rextra-backend/internal/utils"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type (
	PaymentTransactionService interface {
		Prepare(ctx context.Context, userID uuid.UUID, req dto_request.CheckoutPrepareRequest) (dto_response.CheckoutPrepareResponse, error)
		Calculate(ctx context.Context, userID uuid.UUID, req dto_request.CheckoutCalculateRequest) (dto_response.CheckoutCalculateResponse, error)
		Initiate(ctx context.Context, userID uuid.UUID, req dto_request.CheckoutInitiateRequest) (dto_response.CheckoutInitiateResponse, error)
		Repeat(ctx context.Context, userID uuid.UUID, transactionID string) (dto_response.CheckoutRepeatResponse, error)
		Cancel(ctx context.Context, userID uuid.UUID, transactionID string, req dto_request.CheckoutCancelRequest) error
		GetTransaction(ctx context.Context, userID uuid.UUID, transactionID string) (dto_response.MyTransactionDetailResponse, error)
		GetTransactions(ctx context.Context, userID uuid.UUID, req dto_request.MyTransactionFilterRequest) ([]dto_response.MyTransactionListResponse, int64, error)
		GetPaymentChannels(ctx context.Context) ([]dto_response.PaymentChannelResponse, error)
	}

	paymentTransactionService struct {
		repo          repository.PaymentTransactionsRepository
		planRepo      membershipRepo.MembershipPlanRepository
		durationRepo  membershipRepo.MembershipDurationRepository
		promoRepo     PromoRepository
		userRepo      UserRepository
		tripayClient  *tripay.TripayClient
		db            *gorm.DB
	}
)

func NewPaymentTransactionService(
	repo repository.PaymentTransactionsRepository,
	planRepo membershipRepo.MembershipPlanRepository,
	durationRepo membershipRepo.MembershipDurationRepository,
	promoRepo PromoRepository,
	userRepo UserRepository,
	tripayClient *tripay.TripayClient,
	db *gorm.DB,
) PaymentTransactionService {
	return &paymentTransactionService{
		repo:         repo,
		planRepo:     planRepo,
		durationRepo: durationRepo,
		promoRepo:    promoRepo,
		userRepo:     userRepo,
		tripayClient: tripayClient,
		db:           db,
	}
}

var tierRank = map[entity.EnumPlanName]int{
	entity.PLANSTARTER:  0,
	entity.PLANSTANDARD: 0,
	entity.PLANBASIC:    1,
	entity.PLANPRO:      2,
	entity.PLANMAX:      3,
}

func calculateCredit(plan *entity.MembershipPlans, duration *entity.MembershipDuration, remainingDays int) int64 {
	durationPrice := utils.CalculateMembershipPrice(plan.BaseMonthlyPrice, duration.DurationMonth)
	totalDays := duration.DurationMonth * 30
	if totalDays == 0 {
		return 0
	}
	return int64((durationPrice / float64(totalDays)) * float64(remainingDays))
}

// ─── Prepare ─────────────────────────────────────────────────────────────────

func (s *paymentTransactionService) Prepare(ctx context.Context, userID uuid.UUID, req dto_request.CheckoutPrepareRequest) (dto_response.CheckoutPrepareResponse, error) {
	parsedPlanID, err := uuid.Parse(req.PlanID)
	if err != nil {
		return dto_response.CheckoutPrepareResponse{}, myerror.InvalidRequest(err)
	}

	plan, err := s.planRepo.GetByID(ctx, nil, parsedPlanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.CheckoutPrepareResponse{}, myerror.RecordNotFound("membership plan")
		}
		return dto_response.CheckoutPrepareResponse{}, myerror.DatabaseError(err)
	}

	durations, err := s.durationRepo.GetAllByPlanID(ctx, nil, parsedPlanID)
	if err != nil {
		return dto_response.CheckoutPrepareResponse{}, myerror.DatabaseError(err)
	}

	var membershipInfo *dto_response.CheckoutMembershipInfo
	membership, err := s.repo.GetMembershipByUserID(ctx, nil, userID)
	if err == nil {
		changeType := entity.ChangeType(req.ChangeType)
		if changeType == entity.ChangeTypeUpgrade || changeType == entity.ChangeTypeDowngrade {
			currentRank := tierRank[membership.PlanName]
			targetRank := tierRank[plan.PlanName]
			if changeType == entity.ChangeTypeUpgrade && targetRank <= currentRank {
				return dto_response.CheckoutPrepareResponse{}, myerror.New(
					fmt.Sprintf("UPGRADE membutuhkan plan tujuan lebih tinggi dari %s", membership.PlanName),
					myerror.Error_InvalidRequest,
				)
			}
			if changeType == entity.ChangeTypeDowngrade && targetRank >= currentRank {
				return dto_response.CheckoutPrepareResponse{}, myerror.New(
					fmt.Sprintf("DOWNGRADE membutuhkan plan tujuan lebih rendah dari %s", membership.PlanName),
					myerror.Error_InvalidRequest,
				)
			}
		}

		if membership.IsActive && membership.Duration != nil && membership.Plan != nil {
			remainingDays := membership.RemainingDays()
			estimatedCredit := calculateCredit(membership.Plan, membership.Duration, remainingDays)
			durationMonths := 0
			if membership.DurationMonths != nil {
				durationMonths = *membership.DurationMonths
			}
			membershipInfo = &dto_response.CheckoutMembershipInfo{
				PlanName:        string(membership.PlanName),
				DurationMonths:  durationMonths,
				RemainingDays:   remainingDays,
				EstimatedCredit: estimatedCredit,
			}
		}
	}

	voucherCount, _ := s.promoRepo.CountEligibleVouchers(ctx, req.PlanID)

	bundles, err := s.repo.GetActiveTokenBundles(ctx, nil)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto_response.CheckoutPrepareResponse{}, myerror.DatabaseError(err)
	}

	planInfo := dto_response.CheckoutPlanInfo{
		ID:           plan.ID.String(),
		PlanName:     string(plan.PlanName),
		TierLabel:    string(plan.PlanName), // tier label fallback to plan name
		EmblemKey:    "",                    // emblem key fallback config
		DurationMode: "MONTHLY",
	}

	for _, d := range durations {
		if !d.IsActive {
			continue
		}
		price := utils.CalculateMembershipPrice(plan.BaseMonthlyPrice, d.DurationMonth)
		finalPrice := price // For now, we assume price is final unless we have a separate discount config
		pricePerMonth := int64(finalPrice)
		if d.DurationMonth > 1 {
			pricePerMonth = int64(finalPrice / float64(d.DurationMonth))
		}
		planInfo.Durations = append(planInfo.Durations, dto_response.CheckoutDurationInfo{
			ID:             d.ID.String(),
			DurationMonths: d.DurationMonth,
			Price:          int64(price),
			FinalPrice:     int64(finalPrice),
			DiscountPct:    0,
			PricePerMonth:  pricePerMonth,
		})
	}

	var tokenBundleInfos []dto_response.CheckoutTokenBundleInfo
	for _, b := range bundles {
		label := ""
		if b.Label != nil {
			label = *b.Label
		}
		tokenBundleInfos = append(tokenBundleInfos, dto_response.CheckoutTokenBundleInfo{
			ID:          b.ID.String(),
			Name:        b.Name,
			TokenAmount: b.TokenAmount,
			PriceRp:     b.PriceRp,
			Label:       label,
		})
	}

	return dto_response.CheckoutPrepareResponse{
		Plan:                 planInfo,
		CurrentMembership:    membershipInfo,
		EligibleVoucherCount: voucherCount,
		TokenBundles:         tokenBundleInfos,
	}, nil
}

// ─── Calculate ────────────────────────────────────────────────────────────────

func (s *paymentTransactionService) Calculate(ctx context.Context, userID uuid.UUID, req dto_request.CheckoutCalculateRequest) (dto_response.CheckoutCalculateResponse, error) {
	parsedDurationID, err := uuid.Parse(req.DurationID)
	if err != nil {
		return dto_response.CheckoutCalculateResponse{}, myerror.InvalidRequest(err)
	}

	parsedPlanID, err := uuid.Parse(req.PlanID)
	if err != nil {
		return dto_response.CheckoutCalculateResponse{}, myerror.InvalidRequest(err)
	}

	plan, err := s.planRepo.GetByID(ctx, nil, parsedPlanID)
	if err != nil {
		return dto_response.CheckoutCalculateResponse{}, myerror.RecordNotFound("plan")
	}

	planDuration, err := s.repo.GetPlanDurationByID(ctx, nil, parsedDurationID)
	if err != nil {
		return dto_response.CheckoutCalculateResponse{}, myerror.RecordNotFound("plan duration")
	}

	membershipPrice := int64(utils.CalculateMembershipPrice(plan.BaseMonthlyPrice, planDuration.DurationMonth))

	var tokenPrice int64
	if req.TokenBundlePackageID != nil {
		parsedBundleID, err := uuid.Parse(*req.TokenBundlePackageID)
		if err != nil {
			return dto_response.CheckoutCalculateResponse{}, myerror.InvalidRequest(err)
		}
		bundle, err := s.repo.GetTokenBundleByID(ctx, nil, parsedBundleID)
		if err != nil {
			return dto_response.CheckoutCalculateResponse{}, myerror.RecordNotFound("token bundle")
		}
		tokenPrice = bundle.PriceRp
	}

	subtotal := membershipPrice + tokenPrice

	var discountAmount int64
	promoApplied := false
	promoErrMsg := ""
	promoCode := ""
	if req.PromoCode != nil && *req.PromoCode != "" {
		discount, err := s.promoRepo.ValidateAndCalculateDiscount(ctx, *req.PromoCode, req.PlanID, subtotal)
		if err != nil {
			promoErrMsg = err.Error()
		} else {
			discountAmount = discount
			promoApplied = true
			promoCode = *req.PromoCode
		}
	}

	var durationCredit int64
	var remainingDays int
	var creditAvailable int64
	changeType := entity.ChangeType(req.ChangeType)
	if (changeType == entity.ChangeTypeUpgrade || changeType == entity.ChangeTypeDowngrade) && req.UseCredit {
		membership, err := s.repo.GetMembershipByUserID(ctx, nil, userID)
		if err == nil && membership.IsActive && membership.Duration != nil && membership.Plan != nil {
			remainingDays = membership.RemainingDays()
			creditAvailable = calculateCredit(membership.Plan, membership.Duration, remainingDays)
			durationCredit = creditAvailable
		}
	}

	total := subtotal - discountAmount - durationCredit
	if total < 0 {
		total = 0
	}

	return dto_response.CheckoutCalculateResponse{
		MembershipPrice: membershipPrice,
		TokenPrice:      tokenPrice,
		SubtotalAmount:  subtotal,
		DiscountAmount:  discountAmount,
		DurationCredit:  durationCredit,
		TotalAmount:     total,
		PromoApplied:    promoApplied,
		PromoCode:       promoCode,
		PromoErrMsg:     promoErrMsg,
		RemainingDays:   remainingDays,
		CreditAvailable: creditAvailable,
	}, nil
}

// ─── Initiate ─────────────────────────────────────────────────────────────────

func (s *paymentTransactionService) Initiate(ctx context.Context, userID uuid.UUID, req dto_request.CheckoutInitiateRequest) (dto_response.CheckoutInitiateResponse, error) {
	parsedDurationID, err := uuid.Parse(req.DurationID)
	if err != nil {
		return dto_response.CheckoutInitiateResponse{}, myerror.InvalidRequest(err)
	}
	parsedPlanID, err := uuid.Parse(req.PlanID)
	if err != nil {
		return dto_response.CheckoutInitiateResponse{}, myerror.InvalidRequest(err)
	}

	existingPending, err := s.repo.GetPendingTransactionByUserID(ctx, nil, userID)
	if err == nil {
		return dto_response.CheckoutInitiateResponse{}, &ErrPendingTransactionExists{
			TransactionID: existingPending.TransactionID,
		}
	}

	plan, err := s.planRepo.GetByID(ctx, nil, parsedPlanID)
	if err != nil {
		return dto_response.CheckoutInitiateResponse{}, myerror.RecordNotFound("plan")
	}

	planDuration, err := s.repo.GetPlanDurationByID(ctx, nil, parsedDurationID)
	if err != nil || !planDuration.IsActive {
		return dto_response.CheckoutInitiateResponse{}, myerror.RecordNotFound("plan duration")
	}

	membership, err := s.repo.GetMembershipByUserID(ctx, nil, userID)
	if err != nil {
		return dto_response.CheckoutInitiateResponse{}, myerror.DatabaseError(err)
	}

	changeType := entity.ChangeType(req.ChangeType)
	if err := validateChangeType(changeType, membership); err != nil {
		return dto_response.CheckoutInitiateResponse{}, myerror.New(err.Error(), myerror.Error_InvalidRequest)
	}

	membershipPrice := int64(utils.CalculateMembershipPrice(plan.BaseMonthlyPrice, planDuration.DurationMonth))

	var tokenBundlePackageID *uuid.UUID
	var tokenTopupAmount int64
	var tokenTopupPrice int64
	if req.TokenBundlePackageID != nil {
		parsedBundleID, err := uuid.Parse(*req.TokenBundlePackageID)
		if err != nil {
			return dto_response.CheckoutInitiateResponse{}, myerror.InvalidRequest(err)
		}
		bundle, err := s.repo.GetTokenBundleByID(ctx, nil, parsedBundleID)
		if err != nil {
			return dto_response.CheckoutInitiateResponse{}, myerror.RecordNotFound("token bundle")
		}
		tokenBundlePackageID = &parsedBundleID
		tokenTopupAmount = bundle.TokenAmount
		tokenTopupPrice = bundle.PriceRp
	}

	subtotal := membershipPrice + tokenTopupPrice

	var discountAmount int64
	var appliedPromoCode *string
	if req.PromoCode != nil && *req.PromoCode != "" {
		discount, err := s.promoRepo.ValidateAndCalculateDiscount(ctx, *req.PromoCode, req.PlanID, subtotal)
		if err != nil {
			return dto_response.CheckoutInitiateResponse{}, myerror.New(
				fmt.Sprintf("kode promo tidak valid: %s", err.Error()),
				myerror.Error_InvalidRequest,
			)
		}
		discountAmount = discount
		appliedPromoCode = req.PromoCode
	}

	var durationCredit int64
	if (changeType == entity.ChangeTypeUpgrade || changeType == entity.ChangeTypeDowngrade) && req.UseCredit {
		if membership.IsActive && membership.Duration != nil && membership.Plan != nil {
			remainingDays := membership.RemainingDays()
			durationCredit = calculateCredit(membership.Plan, membership.Duration, remainingDays)
		}
	}

	total := subtotal - discountAmount - durationCredit
	if total < 0 {
		total = 0
	}

	items := buildItemsJSON(plan, planDuration, tokenBundlePackageID, tokenTopupAmount, tokenTopupPrice, changeType, membership)
	tripayOrderItems := toTripayOrderItems(items)

	transactionID := "TRX-" + strings.ToUpper(uuid.NewString()[:8])

	user, err := s.userRepo.GetById(ctx, nil, userID.String())
	if err != nil {
		return dto_response.CheckoutInitiateResponse{}, myerror.DatabaseError(err)
	}

	var adminFee int64
	channels, chErr := s.tripayClient.GetPaymentChannel()
	if chErr == nil {
		for _, ch := range channels {
			if ch.Code == req.PaymentMethod {
				adminFee = ch.FeeFlat
				break
			}
		}
	}
	total += adminFee

	merchantRef := transactionID
	expiresAt := time.Now().UTC().Add(6 * time.Hour)
	tripayResp, err := s.tripayClient.CreatePaymentTransaction(tripay.CreatePaymentRequest{
		Method:        req.PaymentMethod,
		MerchantRef:   merchantRef,
		Amount:        total,
		CustomerName:  user.Fullname,
		CustomerEmail: user.Email,
		OrderItems:    tripayOrderItems,
		ExpiredHours:  6,
	})
	if err != nil {
		return dto_response.CheckoutInitiateResponse{}, myerror.New(
			"gagal membuat invoice pembayaran, coba lagi",
			myerror.SystemError,
		)
	}

	var fromPlan *string
	var fromDurationMonths *int
	if membership.IsActive && membership.PlanName != entity.PLANSTARTER {
		fp := string(membership.PlanName)
		fromPlan = &fp
		fromDurationMonths = membership.DurationMonths
	}

	newTrx := entity.NewPaymentTransaction(
		transactionID,
		userID,
		&membership.ID,
		changeType,
		fromPlan,
		req.PlanID,
		fromDurationMonths,
		planDuration.DurationMonth,
		&parsedPlanID,
		&parsedDurationID,
		tokenBundlePackageID,
		tokenTopupAmount,
		tokenTopupPrice,
		subtotal, discountAmount, durationCredit, adminFee, total,
		appliedPromoCode,
		req.PaymentMethod,
		merchantRef,
		user.Fullname, user.Email,
		items,
		&expiresAt,
	)
	newTrx.PaymentExternalID = tripayResp.Reference
	newTrx.PaymentURL = &tripayResp.CheckoutUrl

	if tripayResp.PayCode != "" {
		newTrx.PayCode = &tripayResp.PayCode
	}

	if _, err := s.repo.Create(ctx, nil, newTrx); err != nil {
		return dto_response.CheckoutInitiateResponse{}, myerror.DatabaseError(err)
	}

	if appliedPromoCode != nil {
		_ = s.promoRepo.RecordRedemption(ctx, *appliedPromoCode, userID, newTrx.ID)
	}

	return dto_response.CheckoutInitiateResponse{
		TransactionID: transactionID,
		PaymentURL:    tripayResp.CheckoutUrl,
		PayCode:       tripayResp.PayCode,
		ExpiresAt:     expiresAt.Format(time.RFC3339),
		TotalAmount:   total,
	}, nil
}

// ─── Repeat ───────────────────────────────────────────────────────────────────

func (s *paymentTransactionService) Repeat(ctx context.Context, userID uuid.UUID, transactionID string) (dto_response.CheckoutRepeatResponse, error) {
	trx, err := s.repo.GetTransactionByTransactionID(ctx, nil, transactionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.CheckoutRepeatResponse{}, myerror.RecordNotFound("transaksi")
		}
		return dto_response.CheckoutRepeatResponse{}, myerror.DatabaseError(err)
	}

	if trx.UserID != userID {
		return dto_response.CheckoutRepeatResponse{}, myerror.RecordNotFound("transaksi")
	}

	if trx.PaymentStatus != entity.PaymentStatusCanceled &&
		trx.PaymentStatus != entity.PaymentStatusFailed &&
		trx.PaymentStatus != entity.PaymentStatusExpired {
		return dto_response.CheckoutRepeatResponse{}, myerror.New(
			"hanya transaksi yang dibatalkan atau gagal yang bisa diulangi",
			myerror.Error_InvalidRequest,
		)
	}

	res := dto_response.CheckoutRepeatResponse{
		ChangeType: string(trx.ChangeType),
	}

	if trx.PlanID != nil {
		planIDStr := trx.PlanID.String()
		res.PlanID = &planIDStr
	}
	if trx.DurationID != nil {
		durationIDStr := trx.DurationID.String()
		res.DurationID = &durationIDStr
	}
	if trx.TokenBundlePackageID != nil {
		id := trx.TokenBundlePackageID.String()
		res.TokenBundlePackageID = &id
	}

	if trx.PromoCode != nil {
		planID := ""
		if trx.PlanID != nil {
			planID = trx.PlanID.String()
		}
		_, err := s.promoRepo.ValidateAndCalculateDiscount(ctx, *trx.PromoCode, planID, 0)
		if err != nil {
			res.PromoExpired = true
		} else {
			res.PromoCode = trx.PromoCode
		}
	}

	return res, nil
}

// ─── Cancel ───────────────────────────────────────────────────────────────────

func (s *paymentTransactionService) Cancel(ctx context.Context, userID uuid.UUID, transactionID string, req dto_request.CheckoutCancelRequest) error {
	if req.CancelReason == string(entity.CancelReasonOther) {
		if req.CancelNote == nil || len(*req.CancelNote) < 5 {
			return myerror.New("cancel_note wajib diisi minimal 5 karakter jika alasan adalah OTHER", myerror.Error_InvalidRequest)
		}
	}

	trx, err := s.repo.GetTransactionByTransactionID(ctx, nil, transactionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return myerror.RecordNotFound("transaksi")
		}
		return myerror.DatabaseError(err)
	}

	if trx.UserID != userID {
		return myerror.RecordNotFound("transaksi")
	}

	reason := entity.CancelReason(req.CancelReason)
	if err := trx.Cancel(reason, req.CancelNote); err != nil {
		return myerror.New(err.Error(), myerror.Error_InvalidRequest)
	}

	if _, err := s.repo.Update(ctx, nil, trx); err != nil {
		return myerror.DatabaseError(err)
	}

	return nil
}

// ─── GetTransaction ───────────────────────────────────────────────────────────

func (s *paymentTransactionService) GetTransaction(ctx context.Context, userID uuid.UUID, transactionID string) (dto_response.MyTransactionDetailResponse, error) {
	trx, err := s.repo.GetTransactionByTransactionID(ctx, nil, transactionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.MyTransactionDetailResponse{}, myerror.RecordNotFound("transaksi")
		}
		return dto_response.MyTransactionDetailResponse{}, myerror.DatabaseError(err)
	}

	if trx.UserID != userID {
		return dto_response.MyTransactionDetailResponse{}, myerror.RecordNotFound("transaksi")
	}

	return toTransactionDetailResponse(trx), nil
}

// ─── GetTransactions ──────────────────────────────────────────────────────────

func (s *paymentTransactionService) GetTransactions(ctx context.Context, userID uuid.UUID, req dto_request.MyTransactionFilterRequest) ([]dto_response.MyTransactionListResponse, int64, error) {
	pageSize := req.PageSize
	if pageSize == 0 {
		pageSize = 20
	}

	filter := repository.GetTransactionsFilter{
		Page:     req.Page,
		PageSize: pageSize,
		Status:   req.Status,
	}

	trxs, total, err := s.repo.GetTransactionsByUserIDWithFilter(ctx, nil, userID, filter)
	if err != nil {
		return nil, 0, myerror.DatabaseError(err)
	}

	var res []dto_response.MyTransactionListResponse
	for _, t := range trxs {
		res = append(res, toTransactionListResponse(t))
	}
	return res, total, nil
}

// ─── GetPaymentChannels ───────────────────────────────────────────────────────

func (s *paymentTransactionService) GetPaymentChannels(ctx context.Context) ([]dto_response.PaymentChannelResponse, error) {
	channels, err := s.tripayClient.GetPaymentChannel()
	if err != nil {
		return nil, myerror.New("gagal mengambil daftar metode pembayaran", myerror.SystemError)
	}
	var res []dto_response.PaymentChannelResponse
	for _, ch := range channels {
		res = append(res, dto_response.PaymentChannelResponse{
			Code:     ch.Code,
			Name:     ch.Name,
			Group:    ch.Group,
			IconURL:  ch.IconUrl,
			AdminFee: ch.AdminFee,
		})
	}
	return res, nil
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func validateChangeType(changeType entity.ChangeType, membership entity.Memberships) error {
	isFreePlan := membership.PlanName == entity.PLANSTARTER || !membership.IsActive
	isExpired := membership.ExpiredAt != nil && time.Now().After(*membership.ExpiredAt)

	switch changeType {
	case entity.ChangeTypeNewPurchase:
		if membership.IsActive && !isFreePlan && !isExpired {
			return fmt.Errorf("kamu masih memiliki membership aktif, gunakan RENEWAL/UPGRADE/DOWNGRADE")
		}
	case entity.ChangeTypeRenewal, entity.ChangeTypeUpgrade, entity.ChangeTypeDowngrade:
		if isFreePlan || isExpired {
			return fmt.Errorf("kamu belum memiliki membership aktif, gunakan PEMBELIAN_BARU")
		}
	}
	return nil
}

func buildItemsJSON(
	plan entity.MembershipPlans,
	planDuration entity.MembershipDuration,
	tokenBundleID *uuid.UUID,
	tokenAmount int64,
	tokenPrice int64,
	changeType entity.ChangeType,
	membership entity.Memberships,
) datatypes.JSON {
	type item struct {
		Type            string `json:"type"`
		Name            string `json:"name"`
		Description     string `json:"description"`
		Price           int64  `json:"price"`
		Quantity        int    `json:"quantity"`
		Subtotal        int64  `json:"subtotal"`
		TokenAmount     int64  `json:"token_amount,omitempty"`
		BundlePackageID string `json:"bundle_package_id,omitempty"`
	}

	planName := string(plan.PlanName)

	membershipName := fmt.Sprintf("REXTRA Club %s", planName)
	membershipDesc := fmt.Sprintf("Durasi %d bulan", planDuration.DurationMonth)

	switch changeType {
	case entity.ChangeTypeUpgrade:
		membershipName = fmt.Sprintf("Upgrade Membership — %s → %s", membership.PlanName, planName)
	case entity.ChangeTypeDowngrade:
		membershipName = fmt.Sprintf("Downgrade Membership — %s → %s", membership.PlanName, planName)
	case entity.ChangeTypeRenewal:
		membershipName = fmt.Sprintf("Perpanjang Membership — REXTRA Club %s", planName)
	}

	finalPrice := int64(utils.CalculateMembershipPrice(plan.BaseMonthlyPrice, planDuration.DurationMonth))

	items := []item{
		{
			Type:        "membership",
			Name:        membershipName,
			Description: membershipDesc,
			Price:       finalPrice,
			Quantity:    1,
			Subtotal:    finalPrice,
		},
	}

	if tokenBundleID != nil && tokenAmount > 0 {
		items = append(items, item{
			Type:            "token_addon",
			Name:            "Top up REXTRA Token",
			Description:     fmt.Sprintf("%d token", tokenAmount),
			Price:           tokenPrice,
			Quantity:        1,
			Subtotal:        tokenPrice,
			TokenAmount:     tokenAmount,
			BundlePackageID: tokenBundleID.String(),
		})
	}

	b, _ := json.Marshal(items)
	return datatypes.JSON(b)
}

func toTripayOrderItems(itemsJSON datatypes.JSON) []tripay.TripayOrderItem {
	type rawItem struct {
		Type        string `json:"type"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Price       int64  `json:"price"`
		Quantity    int    `json:"quantity"`
		Subtotal    int64  `json:"subtotal"`
	}
	var items []rawItem
	if err := json.Unmarshal(itemsJSON, &items); err != nil {
		return nil
	}
	var result []tripay.TripayOrderItem
	for _, item := range items {
		result = append(result, tripay.TripayOrderItem{
			SKU:      item.Type,
			Name:     fmt.Sprintf("%s — %s", item.Name, item.Description),
			Price:    int(item.Price),
			Quantity: item.Quantity,
		})
	}
	return result
}

func toTransactionDetailResponse(t entity.PaymentTransactions) dto_response.MyTransactionDetailResponse {
	res := dto_response.MyTransactionDetailResponse{
		ID:             t.ID.String(),
		TransactionID:  t.TransactionID,
		ChangeType:     string(t.ChangeType),
		ToPlan:         t.ToPlan,
		ToDuration:     t.ToDurationMonths,
		PrimaryStatus:  t.PrimaryStatus,
		PaymentStatus:  string(t.PaymentStatus),
		SubtotalAmount: t.SubtotalAmount,
		DiscountAmount: t.DiscountAmount,
		DurationCredit: t.DurationCredit,
		AdminFee:       t.AdminFee,
		TotalAmount:    t.TotalAmount,
		IsCancellable:  t.IsCancellableByUser(),
		IsRepeatable:   t.PaymentStatus == entity.PaymentStatusCanceled || t.PaymentStatus == entity.PaymentStatusFailed || t.PaymentStatus == entity.PaymentStatusExpired,
		CreatedAt:      t.CreatedAt.Format(time.RFC3339),
	}
	if t.FromPlan != nil {
		res.FromPlan = *t.FromPlan
	}
	if t.FromDurationMonths != nil {
		res.FromDuration = *t.FromDurationMonths
	}
	if t.PaymentMethod != nil {
		res.PaymentMethod = *t.PaymentMethod
	}
	if t.PaymentURL != nil {
		res.PaymentURL = *t.PaymentURL
	}
	if t.PayCode != nil {
		res.PayCode = *t.PayCode
	}
	if t.PromoCode != nil {
		res.PromoCode = *t.PromoCode
	}
	if t.CancelReason != nil {
		res.CancelReason = string(*t.CancelReason)
	}
	if t.CancelNote != nil {
		res.CancelNote = *t.CancelNote
	}
	if t.PaidAt != nil {
		res.PaidAt = t.PaidAt.Format(time.RFC3339)
	}
	if t.ExpiresAt != nil {
		res.ExpiresAt = t.ExpiresAt.Format(time.RFC3339)
	}
	if t.CanceledAt != nil {
		res.CanceledAt = t.CanceledAt.Format(time.RFC3339)
	}

	var rawItems []dto_response.MyTransactionItemDetail
	if err := json.Unmarshal(t.Items, &rawItems); err == nil {
		res.Items = rawItems
	}

	return res
}

func toTransactionListResponse(t entity.PaymentTransactions) dto_response.MyTransactionListResponse {
	res := dto_response.MyTransactionListResponse{
		ID:            t.ID.String(),
		TransactionID: t.TransactionID,
		ChangeType:    string(t.ChangeType),
		PrimaryStatus: t.PrimaryStatus,
		PaymentStatus: string(t.PaymentStatus),
		TotalAmount:   t.TotalAmount,
		CreatedAt:     t.CreatedAt.Format(time.RFC3339),
	}
	if t.PaidAt != nil {
		res.PaidAt = t.PaidAt.Format(time.RFC3339)
	}

	type rawItem struct {
		Type        string `json:"type"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	var items []rawItem
	if err := json.Unmarshal(t.Items, &items); err == nil {
		max := 2
		if len(items) < max {
			max = len(items)
		}
		for i := 0; i < max; i++ {
			res.Items = append(res.Items, dto_response.MyTransactionItemPreview{
				Type:        items[i].Type,
				Name:        items[i].Name,
				Description: items[i].Description,
			})
		}
	}

	return res
}
