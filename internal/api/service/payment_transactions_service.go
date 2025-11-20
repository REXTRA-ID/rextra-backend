package service

import (
	"context"
	"rextra-backend/internal/api/repository"
	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	"rextra-backend/payment_handler"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	PaymentTransactionService interface {
		MakeNewTransaction(ctx context.Context, req dto_request.MakeNewTransactionRequest, userId string) (dto_response.MakeNewTransactionResponse, error)
		GetAllTransaction(ctx context.Context) ([]dto_response.GetPaymentTransactionsResponse, error)
		GetAllPaginatedTransaction(ctx context.Context, page, totalPage int) ([]dto_response.GetPaymentTransactionsResponse, error)
		UpdateTransaction(ctx context.Context, paymentStatus string, webhookCallback any) (dto_response.GetPaymentTransactionsResponse, error)
	}

	paymentTransactionService struct {
		paymentService               payment_handler.PaymentService
		paymentTransactionRepository repository.PaymentTransactionsRepository
		membershipRepository         repository.MembershipRepository
		db                           *gorm.DB
	}
)

func NewPaymentTransactionService(paymenService payment_handler.PaymentService,
	paymentTransactionRepository repository.PaymentTransactionsRepository,
	membershipRepository repository.MembershipRepository,
	db *gorm.DB) PaymentTransactionService {

	return &paymentTransactionService{
		paymentService:               paymenService,
		paymentTransactionRepository: paymentTransactionRepository,
		membershipRepository:         membershipRepository,
		db:                           db,
	}
}

func (s *paymentTransactionService) MakeNewTransaction(ctx context.Context, req dto_request.MakeNewTransactionRequest, userId string) (dto_response.MakeNewTransactionResponse, error) {

	var uuidPlanId, uuidDurationId *uuid.UUID
	if req.PlanID != nil {
		strPlanId := *req.PlanID
		*uuidPlanId = uuid.MustParse(strPlanId)
	}

	if req.DurationID != nil {
		strDurationID := *req.DurationID
		*uuidDurationId = uuid.MustParse(strDurationID)
	}

	newTransaction := entity.NewPaymenTransaction(uuid.MustParse(userId), req.PaymentType, "", uuidPlanId, uuidDurationId, req.GrossAmount, "", 0)

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
