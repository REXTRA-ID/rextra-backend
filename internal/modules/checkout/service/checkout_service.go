package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	"rextra-backend/internal/modules/checkout/repository"
	"rextra-backend/internal/pkg/tripay"
	"rextra-backend/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	CheckoutService interface {
		PrepareCheckout(ctx context.Context, req dto_request.CheckoutPrepareRequest, userID string) (dto_response.CheckoutPrepareResponse, error)
		CalculatePrice(ctx context.Context, req dto_request.CheckoutCalculateRequest, userID string) (dto_response.CheckoutCalculateResponse, error)
		InitiateCheckout(ctx context.Context, req dto_request.CheckoutInitiateRequest, userID string) (dto_response.CheckoutInitiateResponse, error)
	}

	checkoutService struct {
		repo         repository.CheckoutRepository
		planRepo     PlanRepositoryAdapter
		promoRepo    PromoRepositoryAdapter
		userRepo     UserRepositoryAdapter
		tripayClient *tripay.TripayClient
		db           *gorm.DB
	}

	PlanRepositoryAdapter interface {
		GetByIDWithDurations(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.MembershipPlans, error)
	}

	PromoRepositoryAdapter interface {
		CountEligibleVouchers(ctx context.Context, planID string) (int, error)
		ValidateAndCalculateDiscount(ctx context.Context, code string, planID string, durationID string, subtotal int64) (int64, error)
		RecordRedemption(ctx context.Context, code string, userID uuid.UUID, transactionID uuid.UUID) error
	}

	UserRepositoryAdapter interface {
		GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.User, error)
	}
)

func NewCheckoutService(repo repository.CheckoutRepository, planRepo PlanRepositoryAdapter, promoRepo PromoRepositoryAdapter, userRepo UserRepositoryAdapter, tripayClient *tripay.TripayClient, db *gorm.DB) CheckoutService {
	return &checkoutService{repo, planRepo, promoRepo, userRepo, tripayClient, db}
}

func (s *checkoutService) PrepareCheckout(ctx context.Context, req dto_request.CheckoutPrepareRequest, userID string) (dto_response.CheckoutPrepareResponse, error) {
	planID, err := uuid.Parse(req.PlanID)
	if err != nil { return dto_response.CheckoutPrepareResponse{}, fmt.Errorf("invalid plan id") }

	plan, err := s.planRepo.GetByIDWithDurations(ctx, nil, planID)
	if err != nil { return dto_response.CheckoutPrepareResponse{}, err }

	voucherCount, _ := s.promoRepo.CountEligibleVouchers(ctx, req.PlanID)

	// Fetch Token Bundles
	var bundles []entity.TokenBundlePackage
	if err := s.db.Table("token_bundle_packages").Order("token_amount asc").Find(&bundles).Error; err != nil {
		fmt.Printf("Error fetching bundles: %v\n", err)
	}
	fmt.Printf("Fetched bundles count: %d\n", len(bundles))

	durations := make([]dto_response.CheckoutDurationInfo, len(plan.PlanDurations))
	for i, d := range plan.PlanDurations {
		durations[i] = dto_response.CheckoutDurationInfo{
			ID: d.ID.String(), DurationMonths: d.DurationMonths, Price: d.Price, FinalPrice: d.FinalPrice,
			DiscountPct: d.DiscountPct,
		}
	}

	tokenBundles := make([]dto_response.CheckoutTokenBundleInfo, len(bundles))
	for i, b := range bundles {
		label := ""
		if b.Label != nil {
			label = *b.Label
		}
		tokenBundles[i] = dto_response.CheckoutTokenBundleInfo{
			ID: b.ID.String(), Name: b.Name, TokenAmount: int64(b.TokenAmount), PriceRp: b.PriceRp, Label: label,
		}
	}

	return dto_response.CheckoutPrepareResponse{
		Plan: dto_response.CheckoutPlanInfo{
			ID: plan.ID.String(), PlanName: string(plan.PlanName), TierLabel: plan.TierLabel, Durations: durations,
		},
		EligibleVoucherCount: voucherCount,
		TokenBundles: tokenBundles,
	}, nil
}

func (s *checkoutService) CalculatePrice(ctx context.Context, req dto_request.CheckoutCalculateRequest, userID string) (dto_response.CheckoutCalculateResponse, error) {
	planID, err := uuid.Parse(req.PlanID)
	if err != nil { return dto_response.CheckoutCalculateResponse{}, fmt.Errorf("invalid plan id") }

	durationID, err := uuid.Parse(req.DurationID)
	if err != nil { return dto_response.CheckoutCalculateResponse{}, fmt.Errorf("invalid duration id") }

	plan, err := s.planRepo.GetByIDWithDurations(ctx, nil, planID)
	if err != nil { return dto_response.CheckoutCalculateResponse{}, fmt.Errorf("plan not found") }

	var duration entity.PlanDuration
	found := false
	for _, d := range plan.PlanDurations {
		if d.ID == durationID { duration = d; found = true; break }
	}
	if !found { return dto_response.CheckoutCalculateResponse{}, fmt.Errorf("duration not found") }

	membershipPrice := duration.FinalPrice
	tokenPrice := int64(0)

	if req.TokenBundlePackageID != nil && *req.TokenBundlePackageID != "" {
		var bundle entity.TokenBundlePackage
		if err := s.db.First(&bundle, "id = ?", *req.TokenBundlePackageID).Error; err == nil {
			tokenPrice = bundle.PriceRp
		}
	}

	subtotal := membershipPrice + tokenPrice
	discount := int64(0)
	promoApplied := false
	if req.PromoCode != nil && *req.PromoCode != "" {
		discount, _ = s.promoRepo.ValidateAndCalculateDiscount(ctx, *req.PromoCode, req.PlanID, durationID.String(), subtotal)
		if discount > 0 { promoApplied = true }
	}

	total := subtotal - discount
	if total < 0 { total = 0 }

	promoCode := ""
	if req.PromoCode != nil { promoCode = *req.PromoCode }

	return dto_response.CheckoutCalculateResponse{
		MembershipPrice: membershipPrice, TokenPrice: tokenPrice, SubtotalAmount: subtotal, DiscountAmount: discount, TotalAmount: total,
		PromoApplied: promoApplied, PromoCode: promoCode,
	}, nil
}

func (s *checkoutService) InitiateCheckout(ctx context.Context, req dto_request.CheckoutInitiateRequest, userID string) (dto_response.CheckoutInitiateResponse, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil { return dto_response.CheckoutInitiateResponse{}, fmt.Errorf("invalid user id") }

	user, err := s.userRepo.GetByID(ctx, nil, parsedUserID)
	if err != nil { return dto_response.CheckoutInitiateResponse{}, fmt.Errorf("user not found") }

	phone := user.PhoneNumber
	if phone == "" { phone = "08123456789" }

	planID, err := uuid.Parse(req.PlanID)
	if err != nil { return dto_response.CheckoutInitiateResponse{}, fmt.Errorf("invalid plan id") }

	durationID, err := uuid.Parse(req.DurationID)
	if err != nil { return dto_response.CheckoutInitiateResponse{}, fmt.Errorf("invalid duration id") }

	plan, err := s.planRepo.GetByIDWithDurations(ctx, nil, planID)
	if err != nil { return dto_response.CheckoutInitiateResponse{}, fmt.Errorf("plan not found") }

	var duration entity.PlanDuration
	found := false
	for _, d := range plan.PlanDurations {
		if d.ID == durationID { duration = d; found = true; break }
	}
	if !found { return dto_response.CheckoutInitiateResponse{}, fmt.Errorf("duration not found") }

	membershipPrice := duration.FinalPrice
	tokenPrice := int64(0)
	itemName := fmt.Sprintf("Membership %s %d Bln", plan.PlanName, duration.DurationMonths)

	if req.TokenBundlePackageID != nil && *req.TokenBundlePackageID != "" {
		var bundle entity.TokenBundlePackage
		if err := s.db.First(&bundle, "id = ?", *req.TokenBundlePackageID).Error; err == nil {
			tokenPrice = bundle.PriceRp
			itemName += fmt.Sprintf(" + Bundle %d Token", bundle.TokenAmount)
		}
	}

	subtotal := membershipPrice + tokenPrice
	discount := int64(0)
	if req.PromoCode != nil && *req.PromoCode != "" {
		discount, _ = s.promoRepo.ValidateAndCalculateDiscount(ctx, *req.PromoCode, req.PlanID, durationID.String(), subtotal)
	}
	total := subtotal - discount

	trxID := utils.GenerateTransactionID()
	merchantRef := utils.GenerateMerchantReff()

	orderItems := []tripay.TripayOrderItem{{SKU: "REXTRA", Name: itemName, Price: total, Quantity: 1}}
	
	// FIX: Set ExpiredTime to 24 hours (86400 seconds)
	expiredTime := time.Now().Unix() + 86400

	tripayResp, err := s.tripayClient.CreatePaymentTransaction(tripay.CreatePaymentRequest{
		Method: req.PaymentMethod, MerchantRef: merchantRef, Amount: total,
		CustomerName: user.Fullname, CustomerEmail: user.Email, CustomerPhone: phone, 
		OrderItems: orderItems, ExpiredTime: expiredTime,
	})
	if err != nil { return dto_response.CheckoutInitiateResponse{}, err }

	// SAVE TO DATABASE
	var bundleID *uuid.UUID
	tokenAmount := int64(0)
	if req.TokenBundlePackageID != nil && *req.TokenBundlePackageID != "" {
		id := uuid.MustParse(*req.TokenBundlePackageID)
		bundleID = &id
		var b entity.TokenBundlePackage
		s.db.First(&b, "id = ?", id)
		tokenAmount = int64(b.TokenAmount)
	}

	itemsJSON, _ := json.Marshal(orderItems)
	expiresAt := time.Unix(tripayResp.ExpiredTime, 0).UTC()

	transaction := entity.NewPaymentTransaction(
		trxID, parsedUserID, nil, entity.ChangeType(req.ChangeType),
		nil, string(plan.PlanName), nil, duration.DurationMonths,
		&planID, &durationID, bundleID, tokenAmount, tokenPrice,
		subtotal, discount, 0, 0, total,
		req.PromoCode, req.PaymentMethod, merchantRef,
		user.Fullname, user.Email, itemsJSON, &expiresAt,
	)
	transaction.SetTripayResponse(tripayResp.Reference, tripayResp.CheckoutURL, tripayResp.PayCode, tripayResp.ExpiredTime)

	_, err = s.repo.CreateTransaction(ctx, nil, transaction)
	if err != nil {
		fmt.Printf("Error saving transaction to DB: %v\n", err)
	}

	return dto_response.CheckoutInitiateResponse{
		TransactionID: trxID, 
		PaymentURL:    tripayResp.CheckoutURL, 
		PayCode:       tripayResp.PayCode, 
		TotalAmount:   total,
		ExpiresAt:     time.Unix(tripayResp.ExpiredTime, 0).Format(time.RFC3339),
	}, nil
}
