package service

import (
	"context"
	"errors"
	"math"
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

// paymentFee represents a payment method's fee structure
type paymentFee struct {
	FlatFee    int     // flat fee in Rupiah
	PercentFee float64 // percentage fee (e.g., 0.03 = 3%, 0.007 = 0.7%)
}

// paymentFeeMap contains Tripay fee data for each payment method
var paymentFeeMap = map[string]paymentFee{
	// Virtual Account - Rp 4.250
	"PERMATAVA":   {FlatFee: 4250},
	"BNIVA":       {FlatFee: 4250},
	"BRIVA":       {FlatFee: 4250},
	"MANDIRIVA":   {FlatFee: 4250},
	"MUAMALATVA":  {FlatFee: 4250},
	"CIMBVA":      {FlatFee: 4250},
	"BSIVA":       {FlatFee: 4250},
	"OCBCVA":      {FlatFee: 4250},
	"DANAMONVA":   {FlatFee: 4250},
	"OTHERBANKVA": {FlatFee: 4250},
	// Virtual Account - Rp 5.500
	"BCAVA": {FlatFee: 5500},
	// Convenience Store - Rp 3.500
	"ALFAMART":  {FlatFee: 3500},
	"INDOMARET": {FlatFee: 3500},
	"ALFAMIDI":  {FlatFee: 3500},
	// E-Wallet - 3%
	"OVO":       {PercentFee: 0.03},
	"DANA":      {PercentFee: 0.03},
	"SHOPEEPAY": {PercentFee: 0.03},
	// QRIS - Rp 750 + 0.7%
	"QRIS":           {FlatFee: 750, PercentFee: 0.007},
	"QRISC":          {FlatFee: 750, PercentFee: 0.007},
	"QRIS2":          {FlatFee: 750, PercentFee: 0.007},
	"QRIS_SHOPEEPAY": {FlatFee: 750, PercentFee: 0.007},
}

type (
	PaymentService interface {
		GetInstructions(ctx context.Context, code string) ([]tripay.InstructionStep, error)
		CreateTransaction(ctx context.Context, req dto_request.CreateTokenTransactionDTORequest, userId string) (dto_response.CreateTransactionDTOResponse, error)
		CountPaymentPrice(ctx context.Context, req dto_request.CreateTokenTransactionDTORequest) (dto_response.CountPaymentPriceDTOResponse, error)
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
		return s.createBundleTransaction(ctx, user, req)
	case "CUSTOM":
		return s.createCustomTransaction(ctx, user, req)
	default:
		return dto_response.CreateTransactionDTOResponse{}, myerror.New("Invalid topup type", myerror.Error_InvalidRequest)
	}
}

func (s *paymentService) createBundleTransaction(ctx context.Context, user entity.User, req dto_request.CreateTokenTransactionDTORequest) (dto_response.CreateTransactionDTOResponse, error) {
	bundle, err := s.tokenBundleRepository.GetByID(ctx, nil, req.BundleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.CreateTransactionDTOResponse{}, myerror.New("Bundle not found", myerror.Error_InvalidRequest)
		}
		return dto_response.CreateTransactionDTOResponse{}, err
	}

	payload := tripay.TransactionPayload{
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
	unitPrice, totalPrice, err := s.calculateCustomPrice(ctx, req)
	if err != nil {
		return dto_response.CreateTransactionDTOResponse{}, err
	}

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

// CountPaymentPrice calculates the total price including fee for a transaction
func (s *paymentService) CountPaymentPrice(ctx context.Context, req dto_request.CreateTokenTransactionDTORequest) (dto_response.CountPaymentPriceDTOResponse, error) {
	var basePrice int
	var err error

	switch req.TopupType {
	case "BUNDLE":
		basePrice, err = s.countBundlePrice(ctx, req)
	case "CUSTOM":
		_, basePrice, err = s.calculateCustomPrice(ctx, req)
	default:
		return dto_response.CountPaymentPriceDTOResponse{}, myerror.New("Invalid topup type", myerror.Error_InvalidRequest)
	}

	if err != nil {
		return dto_response.CountPaymentPriceDTOResponse{}, err
	}

	fee := CalculateFee(req.PaymentMethod, basePrice)

	return dto_response.CountPaymentPriceDTOResponse{
		TotalPrice: basePrice + fee,
		Fee:        fee,
		BasePrice:  basePrice,
	}, nil
}

func (s *paymentService) countBundlePrice(ctx context.Context, req dto_request.CreateTokenTransactionDTORequest) (int, error) {
	bundle, err := s.tokenBundleRepository.GetByID(ctx, nil, req.BundleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, myerror.New("Bundle not found", myerror.Error_InvalidRequest)
		}
		return 0, err
	}

	return int(bundle.PriceRp), nil
}

// calculateCustomPrice returns (unitPrice, totalPrice, error) for custom topup
func (s *paymentService) calculateCustomPrice(ctx context.Context, req dto_request.CreateTokenTransactionDTORequest) (int, int, error) {
	pricingConfig, err := s.customPricingRepository.GetCurrentConfig(ctx, nil)
	if err != nil {
		return 0, 0, err
	}

	if req.Amount < int(pricingConfig.MinToken) || req.Amount > int(pricingConfig.MaxToken) {
		return 0, 0, myerror.New("Invalid amount", myerror.Error_InvalidRequest)
	}

	unitPrice := int(pricingConfig.RecommendedPricePerToken)
	for _, tier := range pricingConfig.Tiers {
		if req.Amount >= int(tier.FromToken) && req.Amount <= int(tier.ToToken) {
			unitPrice = unitPrice * (100 - int(tier.DiscountPct)) / 100
			break
		}
	}
	totalPrice := req.Amount * unitPrice

	return unitPrice, totalPrice, nil
}

// CalculateFee returns the transaction fee for a given payment method and amount
func CalculateFee(paymentMethod string, amount int) int {
	fee, ok := paymentFeeMap[paymentMethod]
	if !ok {
		return 0
	}

	totalFee := fee.FlatFee
	if fee.PercentFee > 0 {
		totalFee += int(math.Ceil(float64(amount) * fee.PercentFee))
	}

	return totalFee
}
