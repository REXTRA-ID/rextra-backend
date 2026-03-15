package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	payrepo "rextra-backend/internal/modules/payment/repository"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/tripay"
	"rextra-backend/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	PaymentTransactionService interface {
		CreateTransaction(ctx context.Context, req dto_request.MakeNewTransactionMembershipRequest, userID, userName, userEmail string) (dto_response.MakeNewTransactionResponse, error)
		GetMyTransactions(ctx context.Context, userID string, req dto_request.MyTransactionFilterRequest) (dto_response.PaginatedPaymentTransactionResponse, error)
		GetTransactionDetail(ctx context.Context, transactionID, ownerUserID string) (dto_response.GetPaymentTransactionDetailResponse, error)
		GetAllTransactions(ctx context.Context, req dto_request.PaymentTransactionFilterRequest) (dto_response.PaginatedPaymentTransactionResponse, error)
		CancelTransaction(ctx context.Context, transactionID string, req dto_request.CancelTransactionRequest) error
		HandleCallback(ctx context.Context, incomingSignature string, rawBody []byte, payload dto_request.TripayCallbackRequest) error
		CalculatePrice(ctx context.Context, req dto_request.CalculatePriceRequest, userID string) (dto_response.PriceCalculationResponse, error)
		GetPaymentChannels(ctx context.Context, amount int64) ([]dto_response.GetPaymentChannelResponse, error)
		SimulateTripayPayment(ctx context.Context, transactionID string) error
	}

	paymentTransactionService struct {
		paymentRepo        payrepo.PaymentTransactionsRepository
		membershipRepo     payrepo.MembershipWriter
		planRepo           payrepo.MembershipPlanReader
		durationRepo       payrepo.PlanDurationReader
		cycleRepo          payrepo.SubscriptionCycleWriter
		discountRepo       payrepo.DiscountReader
		redemptionRepo     payrepo.DiscountRedemptionWriter
		redemptionReverser payrepo.DiscountRedemptionReverser
		tokenLedgerRepo    payrepo.TokenLedgerWriter
		topupTrxRepo       payrepo.TopupTransactionWriter
		tokenWalletRepo    payrepo.TokenWalletUpdater
		poinRepo           payrepo.PoinTransactionWriter
		accessMappingRepo  payrepo.DurationAccessMappingReader
		quotaRepo          payrepo.UserEntitlementQuotaWriter
		tripayClient       tripay.TripayClient
		db                 *gorm.DB
	}

	priceCalculationResult struct {
		SubtotalAmount int64
		DiscountAmount int64
		DurationCredit int64
		TotalAmount    int64
		DiscountID     *uuid.UUID
	}
)

func NewPaymentTransactionService(
	paymentRepo payrepo.PaymentTransactionsRepository,
	membershipRepo payrepo.MembershipWriter,
	planRepo payrepo.MembershipPlanReader,
	durationRepo payrepo.PlanDurationReader,
	cycleRepo payrepo.SubscriptionCycleWriter,
	discountRepo payrepo.DiscountReader,
	redemptionRepo payrepo.DiscountRedemptionWriter,
	redemptionReverser payrepo.DiscountRedemptionReverser,
	tokenLedgerRepo payrepo.TokenLedgerWriter,
	topupTrxRepo payrepo.TopupTransactionWriter,
	tokenWalletRepo payrepo.TokenWalletUpdater,
	poinRepo payrepo.PoinTransactionWriter,
	accessMappingRepo payrepo.DurationAccessMappingReader,
	quotaRepo payrepo.UserEntitlementQuotaWriter,
	tripayClient tripay.TripayClient,
	db *gorm.DB,
) PaymentTransactionService {
	return &paymentTransactionService{
		paymentRepo:        paymentRepo,
		membershipRepo:     membershipRepo,
		planRepo:           planRepo,
		durationRepo:       durationRepo,
		cycleRepo:          cycleRepo,
		discountRepo:       discountRepo,
		redemptionRepo:     redemptionRepo,
		redemptionReverser: redemptionReverser,
		tokenLedgerRepo:    tokenLedgerRepo,
		topupTrxRepo:       topupTrxRepo,
		tokenWalletRepo:    tokenWalletRepo,
		poinRepo:           poinRepo,
		accessMappingRepo:  accessMappingRepo,
		quotaRepo:          quotaRepo,
		tripayClient:       tripayClient,
		db:                 db,
	}
}

func (s *paymentTransactionService) CreateTransaction(ctx context.Context, req dto_request.MakeNewTransactionMembershipRequest, userID, userName, userEmail string) (dto_response.MakeNewTransactionResponse, error) {
	parsedUserID, _ := uuid.Parse(userID)
	parsedPlanID, _ := uuid.Parse(req.PlanID)
	parsedDurationID, _ := uuid.Parse(req.DurationID)

	plan, err := s.planRepo.GetByID(ctx, nil, parsedPlanID)
	if err != nil { return dto_response.MakeNewTransactionResponse{}, myerror.RecordNotFound("plan") }
	duration, err := s.durationRepo.GetByID(ctx, nil, parsedDurationID)
	if err != nil { return dto_response.MakeNewTransactionResponse{}, myerror.RecordNotFound("duration") }

	membership, err := s.membershipRepo.GetByUserID(ctx, nil, parsedUserID)
	isNew := errors.Is(err, gorm.ErrRecordNotFound)

	changeType := entity.ChangeType(req.ChangeType)
	priceCalc, err := s.calculatePrice(ctx, plan, duration, &membership, changeType, req.PromoCode)
	if err != nil { return dto_response.MakeNewTransactionResponse{}, err }

	if priceCalc.TotalAmount <= 0 {
		transactionID := utils.GenerateTransactionID()
		if err := s.activateMembershipDirect(ctx, transactionID, parsedUserID, &membership, isNew, changeType, plan, duration, priceCalc); err != nil { return dto_response.MakeNewTransactionResponse{}, err }
		return dto_response.MakeNewTransactionResponse{TransactionID: transactionID, TotalAmount: 0, Message: "Membership diaktifkan"}, nil
	}

	transactionID := utils.GenerateTransactionID()
	merchantRef := utils.GenerateMerchantReff()
	orderItems := []tripay.TripayOrderItem{{SKU: "membership", Name: fmt.Sprintf("Membership %s", plan.PlanName), Price: int64(priceCalc.SubtotalAmount), Quantity: 1}}

	tripayResp, err := s.tripayClient.CreatePaymentTransaction(tripay.CreatePaymentRequest{
		Method: req.PaymentMethod, MerchantRef: merchantRef, Amount: priceCalc.TotalAmount,
		CustomerName: userName, CustomerEmail: userEmail, OrderItems: orderItems,
	})
	if err != nil { return dto_response.MakeNewTransactionResponse{}, myerror.ProcessingError(err) }

	fromPlan := (*string)(nil)
	if !isNew { fp := string(membership.PlanName); fromPlan = &fp }
	membershipID := (*uuid.UUID)(nil)
	if !isNew { membershipID = &membership.ID }

	newTrx := entity.NewPaymentTransaction(transactionID, parsedUserID, membershipID, changeType, fromPlan, string(plan.PlanName), nil, duration.DurationMonths, &plan.ID, &duration.ID, nil, 0, 0, priceCalc.SubtotalAmount, priceCalc.DiscountAmount, priceCalc.DurationCredit, 0, priceCalc.TotalAmount, req.PromoCode, req.PaymentMethod, merchantRef, userName, userEmail, nil, nil)
	newTrx.SetTripayResponse(tripayResp.Reference, tripayResp.CheckoutURL, tripayResp.PayCode, tripayResp.ExpiredTime)

	if _, err := s.paymentRepo.Create(ctx, nil, newTrx); err != nil { return dto_response.MakeNewTransactionResponse{}, myerror.DatabaseError(err) }
	if req.PromoCode != nil && priceCalc.DiscountAmount > 0 { _ = s.recordDiscountRedemption(ctx, parsedUserID, transactionID, *req.PromoCode, string(plan.PlanName), priceCalc) }

	return dto_response.MakeNewTransactionResponse{TransactionID: transactionID, TotalAmount: priceCalc.TotalAmount, PaymentURL: &tripayResp.CheckoutURL, PayCode: &tripayResp.PayCode}, nil
}

func (s *paymentTransactionService) HandleCallback(ctx context.Context, incomingSignature string, rawBody []byte, payload dto_request.TripayCallbackRequest) error {
	if !s.tripayClient.VerifyCallbackSignature(rawBody, incomingSignature) { return myerror.New("invalid signature", myerror.Error_Unauthorized) }
	trx, err := s.paymentRepo.GetByMerchantRef(ctx, nil, payload.MerchantRef)
	if err != nil { return myerror.RecordNotFound("transaction") }
	if trx.PaymentStatus != entity.PaymentStatusPending { return nil }

	if payload.Status == "PAID" { return s.handlePaidCallback(ctx, trx, payload, rawBody) }
	if payload.Status == "FAILED" { trx.MarkFailed(rawBody) }
	if payload.Status == "EXPIRED" { trx.MarkExpired(rawBody) }
	_, err = s.paymentRepo.Update(ctx, nil, trx)
	return err
}

func (s *paymentTransactionService) handlePaidCallback(ctx context.Context, trx entity.PaymentTransactions, payload dto_request.TripayCallbackRequest, rawBody []byte) error {
	paidAt := time.Now().UTC()
	if payload.PaidAt != nil { paidAt = time.Unix(*payload.PaidAt, 0).UTC() }

	return s.db.Transaction(func(tx *gorm.DB) error {
		trx.MarkPaid(paidAt, rawBody)
		if _, err := s.paymentRepo.Update(ctx, tx, trx); err != nil { return err }

		plan, err := s.planRepo.GetByID(ctx, tx, *trx.PlanID)
		if err != nil { return fmt.Errorf("plan not found: %v", err) }

		duration, err := s.durationRepo.GetByID(ctx, tx, *trx.DurationID)
		if err != nil { return fmt.Errorf("duration not found: %v", err) }

		membership, err := s.membershipRepo.GetByUserID(ctx, tx, trx.UserID)
		isNew := errors.Is(err, gorm.ErrRecordNotFound)
		if isNew {
			membership = entity.NewMembership(trx.UserID, entity.PlanName(plan.PlanName))
			membership.ID = uuid.New()
		}

		membership.Activate(&plan, &duration)

		if duration.TokenAmount > 0 {
			wallet, err := s.tokenWalletRepo.GetOrCreateByUserID(ctx, tx, trx.UserID)
			if err != nil { return fmt.Errorf("failed get/create wallet: %v", err) }

			ledger := entity.TokenLedger{WalletID: wallet.ID, Direction: entity.DirectionIN, Amount: int64(duration.TokenAmount), BalanceBefore: wallet.Balance, BalanceAfter: wallet.Balance + int64(duration.TokenAmount), SourceType: entity.SourceTypeMembership, Description: "Membership tokens"}
			if _, err := s.tokenLedgerRepo.Create(ctx, tx, ledger); err != nil { return fmt.Errorf("failed create ledger: %v", err) }

			if _, err := s.tokenWalletRepo.AddBalance(ctx, tx, trx.UserID, int64(duration.TokenAmount)); err != nil { return fmt.Errorf("failed add balance: %v", err) }
			membership.CurrentTokenBalance += duration.TokenAmount
		}

		if duration.PointsActive && duration.PointsValue > 0 {
			poinBefore := membership.CurrentPoinBalance
			membership.CurrentPoinBalance += duration.PointsValue
			poinTrx := entity.NewPoinTransaction(trx.UserID, membership.ID, entity.PoinTypeEarn, entity.PoinSourceMembershipPurchase, poinBefore, duration.PointsValue, &trx.ID, "Poin membership")
			if _, err := s.poinRepo.Create(ctx, tx, poinTrx); err != nil { return fmt.Errorf("failed create poin trx: %v", err) }
		}

		s.provisionEntitlementQuotas(ctx, tx, membership, duration)

		if isNew { 
			if _, err := s.membershipRepo.Create(ctx, tx, membership); err != nil { return fmt.Errorf("failed create membership: %v", err) }
		} else { 
			if _, err := s.membershipRepo.Update(ctx, tx, membership); err != nil { return fmt.Errorf("failed update membership: %v", err) }
		}

		lastCycle, _ := s.cycleRepo.GetLastCycleNumber(ctx, tx, membership.ID)
		cycle := entity.NewSubscriptionCycle(membership.ID, trx.DurationID, string(plan.PlanName), plan.TierLabel, duration.DurationMonths, paidAt, trx.TotalAmount, *trx.PaymentMethod, trx.TransactionID, lastCycle+1)
		if _, err := s.cycleRepo.Create(ctx, tx, cycle); err != nil { return fmt.Errorf("failed create cycle: %v", err) }

		return nil
	})
}

func (s *paymentTransactionService) provisionEntitlementQuotas(ctx context.Context, tx *gorm.DB, membership entity.Memberships, duration entity.PlanDuration) error {
	if err := s.quotaRepo.InvalidateByMembershipID(ctx, tx, membership.ID); err != nil { return err }
	mappings, _ := s.accessMappingRepo.GetByDurationIDWithEntitlement(ctx, tx, duration.ID)
	for _, m := range mappings {
		if m.Entitlement.RestrictionType == entity.RestrictionFrequencyLimited {
			quota := entity.NewUserEntitlementQuota(membership.ID, m.EntitlementID, m.EntitlementKey, m.UsageLimit, *membership.StartedAt, membership.ExpiredAt)
			if err := s.quotaRepo.Create(ctx, tx, quota); err != nil { return err }
		}
	}
	return nil
}

func (s *paymentTransactionService) CalculatePrice(ctx context.Context, req dto_request.CalculatePriceRequest, userID string) (dto_response.PriceCalculationResponse, error) {
	parsedUserID, _ := uuid.Parse(userID)
	plan, _ := s.planRepo.GetByID(ctx, nil, uuid.MustParse(req.PlanID))
	duration, _ := s.durationRepo.GetByID(ctx, nil, uuid.MustParse(req.DurationID))
	membership, _ := s.membershipRepo.GetByUserID(ctx, nil, parsedUserID)
	priceCalc, err := s.calculatePrice(ctx, plan, duration, &membership, entity.ChangeType(req.ChangeType), req.PromoCode)
	if err != nil { return dto_response.PriceCalculationResponse{}, err }
	return dto_response.PriceCalculationResponse{SubtotalAmount: priceCalc.SubtotalAmount, DiscountAmount: priceCalc.DiscountAmount, DurationCredit: priceCalc.DurationCredit, TotalAmount: priceCalc.TotalAmount, PlanName: string(plan.PlanName), DurationMonths: duration.DurationMonths, TokenAmount: duration.TokenAmount}, nil
}

func (s *paymentTransactionService) calculatePrice(ctx context.Context, plan entity.MembershipPlans, duration entity.PlanDuration, membership *entity.Memberships, changeType entity.ChangeType, promoCode *string) (priceCalculationResult, error) {
	res := priceCalculationResult{SubtotalAmount: duration.FinalPrice}
	if promoCode != nil && *promoCode != "" {
		discount, err := s.discountRepo.GetByCode(ctx, nil, *promoCode)
		if err == nil && discount.IsValid() { res.DiscountAmount = discount.CalculateDiscount(res.SubtotalAmount); res.DiscountID = &discount.ID }
	}
	if (changeType == entity.ChangeTypeUpgrade || changeType == entity.ChangeTypeDowngrade) && membership.DurationID != nil {
		currentDur, err := s.durationRepo.GetByID(ctx, nil, *membership.DurationID)
		if err == nil { res.DurationCredit = currentDur.CalculateCredit(membership.RemainingDays()) }
	}
	res.TotalAmount = res.SubtotalAmount - res.DiscountAmount - res.DurationCredit
	if res.TotalAmount < 0 { res.TotalAmount = 0 }
	return res, nil
}

func (s *paymentTransactionService) activateMembershipDirect(ctx context.Context, transactionID string, userID uuid.UUID, membership *entity.Memberships, isNew bool, changeType entity.ChangeType, plan entity.MembershipPlans, duration entity.PlanDuration, priceCalc priceCalculationResult) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if isNew { membership.ID = uuid.New() }
		membership.Activate(&plan, &duration)
		s.provisionEntitlementQuotas(ctx, tx, *membership, duration)
		if isNew { s.membershipRepo.Create(ctx, tx, *membership) } else { s.membershipRepo.Update(ctx, tx, *membership) }
		return nil
	})
}

func (s *paymentTransactionService) GetMyTransactions(ctx context.Context, userID string, req dto_request.MyTransactionFilterRequest) (dto_response.PaginatedPaymentTransactionResponse, error) {
	page, pageSize := normalizePagination(req.Page, req.PageSize)
	results, total, _ := s.paymentRepo.GetByUserIDPaginated(ctx, nil, uuid.MustParse(userID), page*pageSize, pageSize, req.Status)
	var data []dto_response.GetPaymentTransactionListResponse
	for _, t := range results { data = append(data, toListResponse(t)) }
	return dto_response.PaginatedPaymentTransactionResponse{Data: data, Total: int(total), Page: page, PageSize: pageSize, TotalPages: int(math.Ceil(float64(total) / float64(pageSize)))}, nil
}

func (s *paymentTransactionService) GetTransactionDetail(ctx context.Context, transactionID, ownerUserID string) (dto_response.GetPaymentTransactionDetailResponse, error) {
	trx, _ := s.paymentRepo.GetByTransactionID(ctx, nil, transactionID)
	return toDetailResponse(trx), nil
}

func (s *paymentTransactionService) GetAllTransactions(ctx context.Context, req dto_request.PaymentTransactionFilterRequest) (dto_response.PaginatedPaymentTransactionResponse, error) {
	page, pageSize := normalizePagination(req.Page, req.PageSize)
	results, total, _ := s.paymentRepo.GetAllPaginatedWithFilter(ctx, nil, page*pageSize, pageSize, uuid.Nil, req.Status, req.PlanName, req.Search, time.Time{}, time.Time{})
	var data []dto_response.GetPaymentTransactionListResponse
	for _, t := range results { data = append(data, toListResponse(t)) }
	return dto_response.PaginatedPaymentTransactionResponse{Data: data, Total: int(total), Page: page, PageSize: pageSize, TotalPages: int(math.Ceil(float64(total) / float64(pageSize)))}, nil
}

func (s *paymentTransactionService) CancelTransaction(ctx context.Context, transactionID string, req dto_request.CancelTransactionRequest) error {
	trx, _ := s.paymentRepo.GetByTransactionID(ctx, nil, transactionID)
	trx.Cancel(entity.CancelReasonOther, &req.Note)
	s.paymentRepo.Update(ctx, nil, trx)
	return nil
}

func (s *paymentTransactionService) GetPaymentChannels(ctx context.Context, amount int64) ([]dto_response.GetPaymentChannelResponse, error) {
	var channels []tripay.PaymentChannel
	var err error
	if amount > 0 {
		channels, err = s.tripayClient.GetFeeCalculation(amount, "")
	} else {
		channels, err = s.tripayClient.GetPaymentChannels()
	}
	if err != nil { return nil, err }

	var res []dto_response.GetPaymentChannelResponse
	for _, ch := range channels {
		res = append(res, dto_response.GetPaymentChannelResponse{
			Code:       ch.Code,
			Name:       ch.Name,
			Group:      ch.Group,
			FeeFlat:    ch.TotalFee.Flat,
			FeePercent: ch.TotalFee.Percent,
		})
	}
	return res, nil
}

func (s *paymentTransactionService) SimulateTripayPayment(ctx context.Context, transactionID string) error {
	trx, err := s.paymentRepo.GetByTransactionID(ctx, nil, transactionID)
	if err != nil { return myerror.RecordNotFound("transaction") }
	if trx.PaymentStatus != entity.PaymentStatusPending { return myerror.New("transaction not pending", myerror.Error_InvalidRequest) }

	now := time.Now().Unix()
	payload := dto_request.TripayCallbackRequest{
		Reference:   trx.PaymentExternalID,
		MerchantRef: trx.MerchantRef,
		Status:      "PAID",
		PaidAt:      &now,
	}

	return s.handlePaidCallback(ctx, trx, payload, []byte("{\"simulation\": true}"))
}

func (s *paymentTransactionService) recordDiscountRedemption(ctx context.Context, userID uuid.UUID, transactionID, code, planName string, priceCalc priceCalculationResult) error {
	redemption := entity.NewDiscountRedemption(*priceCalc.DiscountID, userID, transactionID, code, &planName, "User", "MEMBERSHIP", priceCalc.SubtotalAmount, priceCalc.DiscountAmount, priceCalc.TotalAmount)
	s.redemptionRepo.Create(ctx, nil, redemption)
	s.discountRepo.IncrementRedemption(ctx, nil, *priceCalc.DiscountID)
	return nil
}

func toDetailResponse(t entity.PaymentTransactions) dto_response.GetPaymentTransactionDetailResponse {
	createdAt, paidAt, expiresAt, canceledAt := "", "", "", ""
	createdAt = t.CreatedAt.Format(time.RFC3339)
	if t.PaidAt != nil { paidAt = t.PaidAt.Format(time.RFC3339) }
	if t.ExpiresAt != nil { expiresAt = t.ExpiresAt.Format(time.RFC3339) }
	if t.CanceledAt != nil { canceledAt = t.CanceledAt.Format(time.RFC3339) }

	res := dto_response.GetPaymentTransactionDetailResponse{
		ID: t.ID.String(), TransactionID: t.TransactionID, UserID: t.UserID.String(), UserName: t.UserName, UserEmail: t.UserEmail,
		ChangeType: string(t.ChangeType), ToPlan: t.ToPlan, ToDurationMonths: t.ToDurationMonths,
		SubtotalAmount: t.SubtotalAmount, DiscountAmount: t.DiscountAmount, DurationCredit: t.DurationCredit, TotalAmount: t.TotalAmount,
		PromoCode: t.PromoCode, PaymentProvider: string(t.PaymentProvider), PaymentMethod: t.PaymentMethod,
		PaymentStatus: string(t.PaymentStatus), PrimaryStatus: t.PrimaryStatus, PaymentURL: t.PaymentURL, PayCode: t.PayCode,
		PaymentExternalID: t.PaymentExternalID, CreatedAt: createdAt, PaidAt: &paidAt, ExpiresAt: &expiresAt, CanceledAt: &canceledAt,
	}
	if t.FromPlan != nil { res.FromPlan = t.FromPlan }
	if t.FromDurationMonths != nil { res.FromDurationMonths = t.FromDurationMonths }
	if t.CancelReason != nil { cr := string(*t.CancelReason); res.CancelReason = &cr }
	if t.CancelNote != nil { res.CancelNote = t.CancelNote }
	return res
}

func toListResponse(t entity.PaymentTransactions) dto_response.GetPaymentTransactionListResponse {
	createdAt, paidAt, expiresAt := "", "", ""
	createdAt = t.CreatedAt.Format(time.RFC3339)
	if t.PaidAt != nil { paidAt = t.PaidAt.Format(time.RFC3339) }
	if t.ExpiresAt != nil { expiresAt = t.ExpiresAt.Format(time.RFC3339) }

	return dto_response.GetPaymentTransactionListResponse{
		ID: t.ID.String(), TransactionID: t.TransactionID, UserID: t.UserID.String(), UserName: t.UserName, UserEmail: t.UserEmail,
		ToPlan: t.ToPlan, ToDurationMonths: t.ToDurationMonths, ChangeType: string(t.ChangeType),
		TotalAmount: t.TotalAmount, PaymentMethod: t.PaymentMethod, PaymentStatus: string(t.PaymentStatus),
		PrimaryStatus: t.PrimaryStatus, CreatedAt: createdAt, PaidAt: &paidAt, ExpiresAt: &expiresAt,
		PaymentURL: t.PaymentURL, PayCode: t.PayCode,
	}
}

func normalizePagination(page, pageSize int) (int, int) {
	if page < 0 { page = 0 }
	if pageSize <= 0 { pageSize = 20 }
	return page, pageSize
}
