package service

import (
	"context"
	"rextra-backend/internal/api/repository"
	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/payment_handler"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	PaymentTransactionService interface {
		MakeNewTransactionToken(ctx context.Context, req dto_request.MakeNewTransactionTokenRequest, userId string) (dto_response.MakeNewTransactionResponse, error)
		MakeNewTransactionMembership(ctx context.Context, req dto_request.MakeNewTransactionMembershipRequest, userId string) (dto_response.MakeNewTransactionResponse, error)
		GetAllTransaction(ctx context.Context) ([]dto_response.GetPaymentTransactionsResponse, error)
		GetAllPaginatedTransaction(ctx context.Context, page, totalPage int) ([]dto_response.GetPaymentTransactionsResponse, error)
		UpdateTransaction(ctx context.Context, paymentStatus string, webhookCallback any) (dto_response.GetPaymentTransactionsResponse, error)
	}

	paymentTransactionService struct {
		tokenUsageHistoryRepository  repository.TokenUsageHistoryRepository
		tokenTransactionRepository   repository.TokenTransactionRepository
		paymentService               payment_handler.PaymentService
		paymentTransactionRepository repository.PaymentTransactionsRepository
		membershipPlanRepository     repository.MembershipPlanRepository
		membershipRepository         repository.MembershipRepository
		db                           *gorm.DB
	}
)

// maaf kalau kode ini kepanjangan
func NewPaymentTransactionService(tokenUsageHistoryRepository repository.TokenUsageHistoryRepository,
	tokenTransactionRepository repository.TokenTransactionRepository,
	paymenService payment_handler.PaymentService,
	paymentTransactionRepository repository.PaymentTransactionsRepository,
	membershipPlanRepository repository.MembershipPlanRepository,
	membershipRepository repository.MembershipRepository,
	db *gorm.DB) PaymentTransactionService {

	return &paymentTransactionService{
		tokenUsageHistoryRepository:  tokenUsageHistoryRepository,
		tokenTransactionRepository:   tokenTransactionRepository,
		paymentService:               paymenService,
		paymentTransactionRepository: paymentTransactionRepository,
		membershipPlanRepository:     membershipPlanRepository,
		membershipRepository:         membershipRepository,
		db:                           db,
	}
}

func (s *paymentTransactionService) MakeNewTransactionToken(ctx context.Context, req dto_request.MakeNewTransactionTokenRequest, userId string) (dto_response.MakeNewTransactionResponse, error) {
	uuidUserID := uuid.MustParse(userId)

	newTransaction := entity.NewPaymenTransaction(uuidUserID, string(entity.TOKENSTANDALONE), "", nil, nil, req.GrossAmount, "", req.TokenQuantity)

	_, err := s.paymentTransactionRepository.Create(ctx, nil, newTransaction)
	if err != nil {
		return dto_response.MakeNewTransactionResponse{}, myerror.ProcessingError(err)
	}

	return dto_response.MakeNewTransactionResponse{
		RedirectURL: "",
	}, nil
}

func (s *paymentTransactionService) MakeNewTransactionMembership(ctx context.Context, req dto_request.MakeNewTransactionMembershipRequest, userId string) (dto_response.MakeNewTransactionResponse, error) {

	uuidUserID := uuid.MustParse(userId)

	uuidPlanId := uuid.MustParse(req.PlanID)
	uuidDurationId := uuid.MustParse(req.DurationID)

	newTransaction := entity.NewPaymenTransaction(uuidUserID, string(entity.MEMBERSHIP), "", &uuidPlanId, &uuidDurationId, req.GrossAmount, "", 0)

	_, err := s.paymentTransactionRepository.Create(ctx, nil, newTransaction)
	if err != nil {
		return dto_response.MakeNewTransactionResponse{}, err
	}

	return dto_response.MakeNewTransactionResponse{
		RedirectURL: "",
	}, nil
}

func (s *paymentTransactionService) GetAllTransaction(ctx context.Context) ([]dto_response.GetPaymentTransactionsResponse, error) {
	return []dto_response.GetPaymentTransactionsResponse{}, nil
}

func (s *paymentTransactionService) GetAllPaginatedTransaction(ctx context.Context, page, totalPage int) ([]dto_response.GetPaymentTransactionsResponse, error) {
	return []dto_response.GetPaymentTransactionsResponse{}, nil
}

func (s *paymentTransactionService) UpdateTransaction(ctx context.Context, paymentStatus string, webhookCallback any) (dto_response.GetPaymentTransactionsResponse, error) {
	return dto_response.GetPaymentTransactionsResponse{}, nil
}
