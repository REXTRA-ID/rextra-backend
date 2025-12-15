package service

import (
	"context"
	"errors"
	"rextra-backend/internal/api/repository"
	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/utils"
	"rextra-backend/payment_handler"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	PaymentTransactionService interface {
		MakeNewTransactionToken(ctx context.Context, req dto_request.MakeNewTransactionTokenRequest, userId, email string) (dto_response.MakeNewTransactionResponse, error)
		MakeNewTransactionMembership(ctx context.Context, req dto_request.MakeNewTransactionMembershipRequest, userId, email string) (dto_response.MakeNewTransactionResponse, error)
		GetAllTransaction(ctx context.Context) ([]dto_response.GetPaymentTransactionsResponse, error)
		GetAllPaginatedTransaction(ctx context.Context, page, totalPage int) ([]dto_response.GetPaymentTransactionsResponse, error)
		UpdateTransaction(ctx context.Context, req dto_request.PaymentWebhookRequest, webhookCallback []byte) (dto_response.GetPaymentTransactionsResponse, error)
	}

	paymentTransactionService struct {
		tokenTransactionRepository   repository.TokenTransactionRepository
		poinTransactionRepository    repository.PoinTransactionsRepository
		paymentService               payment_handler.PaymentService
		paymentTransactionRepository repository.PaymentTransactionsRepository
		promoCodeRepository          repository.PromoCodeRepository
		promoCodeUsageRepository     repository.PromoCodeUsageRepository
		membershipPlanRepository     repository.MembershipPlanRepository
		membershipDurationRepository repository.MembershipDurationRepository
		membershipRepository         repository.MembershipRepository
		db                           *gorm.DB
	}
)

func NewPaymentTransactionService(
	tokenTransactionRepository repository.TokenTransactionRepository,
	poinTransactionRepository repository.PoinTransactionsRepository,
	paymenService payment_handler.PaymentService,
	paymentTransactionRepository repository.PaymentTransactionsRepository,
	promoCodeRepsitory repository.PromoCodeRepository,
	prmoCodeUsageRepsitory repository.PromoCodeUsageRepository,
	membershipPlanRepository repository.MembershipPlanRepository,
	membershipDuration repository.MembershipDurationRepository,
	membershipRepository repository.MembershipRepository,
	db *gorm.DB) PaymentTransactionService {

	return &paymentTransactionService{
		tokenTransactionRepository:   tokenTransactionRepository,
		poinTransactionRepository:    poinTransactionRepository,
		paymentService:               paymenService,
		paymentTransactionRepository: paymentTransactionRepository,
		promoCodeRepository:          promoCodeRepsitory,
		promoCodeUsageRepository:     prmoCodeUsageRepsitory,
		membershipPlanRepository:     membershipPlanRepository,
		membershipDurationRepository: membershipDuration,
		membershipRepository:         membershipRepository,
		db:                           db,
	}
}

func (s *paymentTransactionService) MakeNewTransactionToken(ctx context.Context, req dto_request.MakeNewTransactionTokenRequest, userId, email string) (dto_response.MakeNewTransactionResponse, error) {
	uuidUserID := uuid.MustParse(userId)

	url, invoiceId, err := s.paymentService.CreateTokenPaymentRequest(req, email)
	if err != nil {
		return dto_response.MakeNewTransactionResponse{}, err
	}

	grossAmount := req.GrossAmount

	var promo *entity.PromoCodes
	if req.PromoCode != nil {
		promo, err = s.promoCodeRepository.GetPromoCodeByCodeName(ctx, nil, *req.PromoCode)
		if err != nil {
			return dto_response.MakeNewTransactionResponse{}, err
		}
		isApplicable := isDiscountApplicableToken(req.TokenQuantity, req.PaymentType, promo)
		if !isApplicable {
			return dto_response.MakeNewTransactionResponse{}, errors.New("discount is not applicable")
		}

		if promo.CurrentUsage > 0 {
			used, err := s.promoCodeUsageRepository.IsPromoCodeUsedByUserID(ctx, nil, uuidUserID)
			if err != nil {
				return dto_response.MakeNewTransactionResponse{}, err
			}
			if used {
				return dto_response.MakeNewTransactionResponse{}, errors.New("you already use this discount")
			}
		}
		tokenQuantity := req.TokenQuantity
		if promo.BonusToken != nil {
			tokenQuantity += *promo.BonusToken
		}

		if promo.DiscountType == entity.DiscountFixedAmount {
			grossAmount = utils.ApplyDiscountFixed(grossAmount, int(promo.DiscountValue))
		}

		if promo.DiscountType == entity.DiscountPercentage {
			grossAmount = utils.ApplyDiscountPercentage(grossAmount, int(promo.DiscountValue))
		}

		newTransaction := entity.NewPaymenTransaction(uuidUserID, string(entity.TOKENSTANDALONE), nil, nil, grossAmount, "", tokenQuantity)

		_, err = s.paymentTransactionRepository.Create(ctx, nil, newTransaction)
		if err != nil {
			return dto_response.MakeNewTransactionResponse{}, myerror.ProcessingError(err)
		}
	}

	newTransaction := entity.NewPaymenTransaction(uuidUserID, string(entity.TOKENSTANDALONE), nil, nil, req.GrossAmount, invoiceId, req.TokenQuantity)

	_, err = s.paymentTransactionRepository.Create(ctx, nil, newTransaction)
	if err != nil {
		return dto_response.MakeNewTransactionResponse{}, myerror.ProcessingError(err)
	}

	return dto_response.MakeNewTransactionResponse{
		RedirectURL: url,
	}, nil
}

func (s *paymentTransactionService) MakeNewTransactionMembership(ctx context.Context, req dto_request.MakeNewTransactionMembershipRequest, userId, email string) (dto_response.MakeNewTransactionResponse, error) {

	uuidUserID := uuid.MustParse(userId)

	uuidPlanId := uuid.MustParse(req.PlanID)
	uuidDurationId := uuid.MustParse(req.DurationID)

	plan, err := s.membershipPlanRepository.GetByID(ctx, nil, uuidPlanId)
	if err != nil {
		return dto_response.MakeNewTransactionResponse{}, err
	}

	duration, err := s.membershipDurationRepository.GetByID(ctx, nil, uuidDurationId)
	if err != nil {
		return dto_response.MakeNewTransactionResponse{}, err
	}

	grossAmount := utils.CalculateMembershipPrice(plan.BaseMonthlyPrice, duration.DurationMonth)

	var promo *entity.PromoCodes
	if req.PromoCode != nil {
		promo, err = s.promoCodeRepository.GetPromoCodeByCodeName(ctx, nil, *req.PromoCode)
		if err != nil {
			return dto_response.MakeNewTransactionResponse{}, err
		}
		isApplicable := isDiscountApplicableMembership(string(plan.PlanName), duration.DurationMonth, promo, string(entity.MEMBERSHIP))
		if !isApplicable {
			return dto_response.MakeNewTransactionResponse{}, errors.New("discount is not applicable")
		}

		if promo.CurrentUsage > 0 {
			used, err := s.promoCodeUsageRepository.IsPromoCodeUsedByUserID(ctx, nil, uuidUserID)
			if err != nil {
				return dto_response.MakeNewTransactionResponse{}, err
			}
			if used {
				return dto_response.MakeNewTransactionResponse{}, errors.New("you already use this discount")
			}
		}

		if promo.DiscountType == entity.DiscountFixedAmount {
			grossAmount = utils.ApplyDiscountFixed(grossAmount, int(promo.DiscountValue))
		}

		if promo.DiscountType == entity.DiscountPercentage {
			grossAmount = utils.ApplyDiscountPercentage(grossAmount, int(promo.DiscountValue))
		}
	}

	if grossAmount == 0 {
		newTransaction := entity.NewPaymenTransaction(uuidUserID, string(entity.MEMBERSHIP), &uuidPlanId, &uuidDurationId, grossAmount, "", *promo.BonusToken)
		_, err = s.paymentTransactionRepository.Create(ctx, nil, newTransaction)
		if err != nil {
			return dto_response.MakeNewTransactionResponse{}, err
		}

		userMembership, err := s.membershipRepository.GetByUserID(ctx, nil, uuidUserID)
		if err != nil {
			return dto_response.MakeNewTransactionResponse{}, err
		}

		if err := s.updateMembershipTransaction(ctx, &newTransaction, &userMembership, newTransaction.Plan, newTransaction.Duration, promo.BonusToken); err != nil {
			return dto_response.MakeNewTransactionResponse{}, err
		}

		return dto_response.MakeNewTransactionResponse{
			Message: "Selamat! Anda mendapat membership GRATIS",
		}, nil
	}

	url, invoiceId, err := s.paymentService.CreateMembershipPaymentRequest(grossAmount, string(plan.PlanName), email)
	if err != nil {
		return dto_response.MakeNewTransactionResponse{}, err
	}

	newTransaction := entity.NewPaymenTransaction(uuidUserID, string(entity.MEMBERSHIP), &uuidPlanId, &uuidDurationId, grossAmount, invoiceId, 0)

	_, err = s.paymentTransactionRepository.Create(ctx, nil, newTransaction)
	if err != nil {
		return dto_response.MakeNewTransactionResponse{}, err
	}

	return dto_response.MakeNewTransactionResponse{
		RedirectURL: url,
	}, nil
}

func (s *paymentTransactionService) GetAllTransaction(ctx context.Context) ([]dto_response.GetPaymentTransactionsResponse, error) {
	return []dto_response.GetPaymentTransactionsResponse{}, nil
}

func (s *paymentTransactionService) GetAllPaginatedTransaction(ctx context.Context, page, totalPage int) ([]dto_response.GetPaymentTransactionsResponse, error) {
	return []dto_response.GetPaymentTransactionsResponse{}, nil
}

func (s *paymentTransactionService) UpdateTransaction(ctx context.Context, req dto_request.PaymentWebhookRequest, webhookCallback []byte) (dto_response.GetPaymentTransactionsResponse, error) {

	transaction, err := s.paymentTransactionRepository.GetByPaymentInvoiceID(ctx, nil, req.OrderID)
	if err != nil {
		return dto_response.GetPaymentTransactionsResponse{}, err
	}

	transaction.UpdateMidtransTransaction(req.TransactionStatus, req.TransactionTime, req.PaymentType, webhookCallback)

	if transaction.PaymentType == entity.MEMBERSHIP && (transaction.PlanID != nil && transaction.DurationID != nil) {

		userMembership, err := s.membershipRepository.GetByUserID(ctx, nil, transaction.UserID)
		if err != nil {
			return dto_response.GetPaymentTransactionsResponse{}, err
		}

		if err := s.updateMembershipTransaction(ctx, &transaction, &userMembership, transaction.Plan, transaction.Duration, nil); err != nil {
			return dto_response.GetPaymentTransactionsResponse{}, err
		}
	} else if transaction.PaymentType == entity.TOKENSTANDALONE && transaction.TokenQuantity != 0 {
		userMembership, err := s.membershipRepository.GetByUserID(ctx, nil, transaction.UserID)
		if err != nil {
			return dto_response.GetPaymentTransactionsResponse{}, err
		}

		if err := s.updateTokenTransaction(ctx, &transaction, &userMembership); err != nil {
			return dto_response.GetPaymentTransactionsResponse{}, err
		}
	} else {
		return dto_response.GetPaymentTransactionsResponse{}, errors.New("invalid payment")
	}

	return dto_response.GetPaymentTransactionsResponse{
		ID:            transaction.ID.String(),
		UserID:        transaction.UserID.String(),
		PaymentType:   string(transaction.PaymentType),
		GrossAmount:   transaction.GrossAmount,
		FinalAmount:   transaction.FinalAmount,
		PaymentStatus: string(transaction.PaymentStatus),
	}, nil
}

func (s *paymentTransactionService) updateMembershipTransaction(ctx context.Context, transaction *entity.PaymentTransactions, userMembership *entity.Memberships, plan *entity.MembershipPlans, duration *entity.MembershipDuration, bonusToken *int) error {
	userMembership.UpdateMembership(plan, duration)
	tokenAmount := userMembership.CalculateTotalToken(bonusToken)
	poinAmount := userMembership.CalculateRextraPoin()

	curretnTokenBalance := userMembership.CurrentTokenBalance
	currentPoinBalance := userMembership.CurrentPoinBalance

	tokenTransaction := entity.NewTokenTransaction(transaction.UserID, userMembership.ID, string(entity.TOKENPURCHASEMEMBERSHIP), curretnTokenBalance, tokenAmount, "")

	poinTransaction := entity.NewPoinTransaction(transaction.UserID, string(entity.EARN), string(entity.MEMBERSHIPPURCHASE), currentPoinBalance, poinAmount)

	_, err := s.membershipRepository.Update(ctx, nil, *userMembership)
	if err != nil {
		return err
	}

	_, err = s.tokenTransactionRepository.Create(ctx, nil, tokenTransaction)
	if err != nil {
		return err
	}
	_, err = s.poinTransactionRepository.Create(ctx, nil, poinTransaction)
	if err != nil {
		return err
	}

	return nil
}

func (s *paymentTransactionService) updateTokenTransaction(ctx context.Context, transaction *entity.PaymentTransactions, userMembership *entity.Memberships) error {
	currentTokenBalance := userMembership.CurrentTokenBalance

	tokenQuntity := transaction.TokenQuantity

	if userMembership.MembershipStatus == entity.PLANNONMEMBER {
		plan, err := s.membershipPlanRepository.GetByPlanName(ctx, nil, string(entity.PLANSTARTER))
		if err != nil {
			return err
		}

		duration, err := s.membershipDurationRepository.GetByDurationMonth(ctx, nil, 30)
		if err != nil {
			return err
		}

		userMembership.UpdateMembership(&plan, &duration)
	}

	userMembership.AddMembershipToken(tokenQuntity)

	tokenTransaction := entity.NewTokenTransaction(transaction.UserID, userMembership.ID, string(entity.TOKENPURCHASEMEMBERSHIP), currentTokenBalance, tokenQuntity, "")

	_, err := s.tokenTransactionRepository.Create(ctx, nil, tokenTransaction)
	if err != nil {
		return err
	}

	_, err = s.membershipRepository.Update(ctx, nil, *userMembership)
	if err != nil {
		return err
	}

	return nil
}

func isDiscountApplicableMembership(planName string, durationMont int,
	promo *entity.PromoCodes, transactionType string,
) bool {

	if string(promo.PromoType) != transactionType {
		return false
	}

	if !promo.IsValid() {
		return false
	}

	if promo.UsageExceed() {
		return false
	}

	statePlan, _ := promo.ApplyPromoApplicablePlans(planName)
	stateDuration, _ := promo.ApplyPromoApplicaleDuration(durationMont)

	return statePlan && stateDuration

}

func isDiscountApplicableToken(tokenQuantity int, transactionType string, promo *entity.PromoCodes) bool {
	if string(promo.PromoType) != transactionType {
		return false
	}

	if !promo.IsValid() {
		return false
	}

	if promo.UsageExceed() {
		return false
	}

	if *promo.MinTokenPurchase > tokenQuantity {
		return false
	}

	return true
}
