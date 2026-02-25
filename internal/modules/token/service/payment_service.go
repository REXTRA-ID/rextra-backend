package service

import (
	"context"
	"errors"
	"fmt"
	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	authrepo "rextra-backend/internal/modules/auth/repository"
	tokenrepo "rextra-backend/internal/modules/token/repository"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/tripay"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	PaymentService interface {
		GetInstructions(ctx context.Context, code string) ([]tripay.InstructionStep, error)
		CreateTransaction(ctx context.Context, req dto_request.CreateTokenTransactionDTORequest, userId string) (dto_response.CreateTransactionDTOResponse, error)
	}

	paymentService struct {
		topUpRepository         tokenrepo.TopupTransactionRepository
		tokenLedgerRepository   tokenrepo.TokenLedgerRepository
		tokenBundleRepository   tokenrepo.TokenBundlePackageRepository
		userRepository          authrepo.UserRepository
		customPricingRepository tokenrepo.CustomPricingRepository
	}
)

func NewPaymentService(topUpRepository tokenrepo.TopupTransactionRepository, tokenLedgerRepository tokenrepo.TokenLedgerRepository, userRepository authrepo.UserRepository, tokenBundleRepository tokenrepo.TokenBundlePackageRepository, customPricingRepository tokenrepo.CustomPricingRepository) PaymentService {
	return &paymentService{
		topUpRepository:         topUpRepository,
		tokenLedgerRepository:   tokenLedgerRepository,
		userRepository:          userRepository,
		tokenBundleRepository:   tokenBundleRepository,
		customPricingRepository: customPricingRepository,
	}
}

func (s *paymentService) GetInstructions(ctx context.Context, code string) ([]tripay.InstructionStep, error) {
	res, err := tripay.GetInstruction(code)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *paymentService) CreateTransaction(ctx context.Context, req dto_request.CreateTokenTransactionDTORequest, userId string) (dto_response.CreateTransactionDTOResponse, error) {
	user, err := s.userRepository.GetById(ctx, nil, userId)
	if err != nil {
		return dto_response.CreateTransactionDTOResponse{}, err
	}

	switch req.TopupType {
	case "BUNDLE":
		res, err := s.createBundleTransaction(ctx, user, req)
		return res, err
	case "CUSTOM":
		res, err := s.createCustomTransaction(ctx, user, req)
		return res, err
	default:
		return dto_response.CreateTransactionDTOResponse{}, myerror.New("Invalid topup type", myerror.Error_InvalidRequest)
	}
}

func (s *paymentService) createBundleTransaction(ctx context.Context, user entity.User, req dto_request.CreateTokenTransactionDTORequest) (dto_response.CreateTransactionDTOResponse, error) {
	var payload tripay.TransactionPayload
	bundle, err := s.tokenBundleRepository.GetByID(ctx, nil, req.BundleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.CreateTransactionDTOResponse{}, myerror.New("Bundle not found", myerror.Error_InvalidRequest)
		}
		return dto_response.CreateTransactionDTOResponse{}, err
	}

	payload = tripay.TransactionPayload{
		Method:        req.PaymentMethod,
		Amount:        int(bundle.PriceRp),
		CustomerName:  user.Fullname,
		CustomerEmail: user.Email,
		CustomerPhone: user.PhoneNumber,
	}

	orderItems := []tripay.OrderItem{
		{
			"name":      bundle.Name,
			"price":     bundle.PriceRp,
			"quantity":  1,
			"bundle_id": bundle.ID,
		},
	}

	paymentRes, err := tripay.CreateTransaction(payload, orderItems)
	if err != nil {
		return dto_response.CreateTransactionDTOResponse{}, err
	}

	expiredTime := time.Unix(paymentRes.ExpiredTime, 0)
	topup, err := s.topUpRepository.Create(ctx, nil, &entity.TopupTransaction{
		UserID:          user.ID,
		Type:            "BUNDLE",
		BundlePackageID: (*uuid.UUID)(&bundle.ID),
		TokenAmount:     bundle.TokenAmount,
		TotalPriceRp:    int64(paymentRes.Amount) + int64(paymentRes.TotalFee),
		Status:          "PENDING",
		InvoiceID:       paymentRes.Reference,
		Provider:        &paymentRes.PaymentMethod,
		ExpiredAt:       &expiredTime,
		Metadata: entity.TopupMetadata{
			"callback_url": paymentRes.CallbackURL,
			"return_url":   paymentRes.ReturnURL,
			"amount":       bundle.PriceRp,
			"fee_merchant": paymentRes.FeeMerchant,
			"fee_customer": paymentRes.FeeCustomer,
			"pay_code":     paymentRes.PayCode,
			"pay_url":      paymentRes.PayURL,
			"checkout_url": paymentRes.CheckoutURL,
		},
	})

	if err != nil {
		return dto_response.CreateTransactionDTOResponse{}, err
	}

	return dto_response.CreateTransactionDTOResponse{
		ID:        topup.ID.String(),
		Amount:    int(topup.TotalPriceRp),
		Invoice:   topup.InvoiceID,
		ExpiredAt: *topup.ExpiredAt,
		Status:    string(topup.Status),
		Metadata:  topup.Metadata,
	}, nil
}

func (s *paymentService) createCustomTransaction(ctx context.Context, user entity.User, req dto_request.CreateTokenTransactionDTORequest) (dto_response.CreateTransactionDTOResponse, error) {
	// Implement the logic to create a transaction using the provided DTO
	pricingConfig, err := s.customPricingRepository.GetCurrentConfig(ctx, nil)
	if err != nil {
		return dto_response.CreateTransactionDTOResponse{}, err
	}

	if req.Amount < int(pricingConfig.MinToken) || req.Amount > int(pricingConfig.MaxToken) {
		return dto_response.CreateTransactionDTOResponse{}, myerror.New("Invalid amount", myerror.Error_InvalidRequest)
	}

	unitPrice := int(pricingConfig.RecommendedPricePerToken)
	for _, tier := range pricingConfig.Tiers {
		if req.Amount >= int(tier.FromToken) && req.Amount <= int(tier.ToToken) {
			unitPrice = unitPrice * (100 - int(tier.DiscountPct)) / 100
			break
		}
	}
	totalPrice := req.Amount * unitPrice

	fmt.Println("unitPrice", unitPrice, "totalPrice", totalPrice)

	payload := tripay.TransactionPayload{
		Method:        req.PaymentMethod,
		Amount:        totalPrice,
		CustomerName:  user.Fullname,
		CustomerEmail: user.Email,
		CustomerPhone: user.PhoneNumber,
	}

	orderItems := []tripay.OrderItem{
		{
			"name":      "custom",
			"price":     unitPrice,
			"quantity":  req.Amount,
			"bundle_id": "",
		},
	}

	paymentRes, err := tripay.CreateTransaction(payload, orderItems)
	if err != nil {
		return dto_response.CreateTransactionDTOResponse{}, err
	}

	expiredTime := time.Unix(paymentRes.ExpiredTime, 0)
	topup, err := s.topUpRepository.Create(ctx, nil, &entity.TopupTransaction{
		UserID:       user.ID,
		Type:         "CUSTOM",
		TokenAmount:  int64(req.Amount),
		TotalPriceRp: int64(paymentRes.Amount) + int64(paymentRes.TotalFee),
		Status:       "PENDING",
		InvoiceID:    paymentRes.Reference,
		Provider:     &paymentRes.PaymentMethod,
		ExpiredAt:    &expiredTime,
		Metadata: entity.TopupMetadata{
			"callback_url": paymentRes.CallbackURL,
			"return_url":   paymentRes.ReturnURL,
			"amount":       totalPrice,
			"fee_merchant": paymentRes.FeeMerchant,
			"fee_customer": paymentRes.FeeCustomer,
			"pay_code":     paymentRes.PayCode,
			"pay_url":      paymentRes.PayURL,
			"checkout_url": paymentRes.CheckoutURL,
		},
	})

	if err != nil {
		return dto_response.CreateTransactionDTOResponse{}, err
	}

	return dto_response.CreateTransactionDTOResponse{
		ID:        topup.ID.String(),
		Amount:    int(topup.TotalPriceRp),
		Invoice:   topup.InvoiceID,
		ExpiredAt: *topup.ExpiredAt,
		Status:    string(topup.Status),
		Metadata:  topup.Metadata,
	}, nil
}
