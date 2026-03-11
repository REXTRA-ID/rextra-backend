package service

import (
	"context"
	"errors"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	myerror "rextra-backend/internal/pkg/error"
	"time"

	"rextra-backend/internal/modules/token/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	TopupTransactionService interface {
		GetAll(ctx context.Context, filters repository.TopupTransactionFilters, limit, offset int) ([]dto_response.TopupTransactionDTOResponse, int64, error)
		GetByID(ctx context.Context, id string) (*dto_response.TopupTransactionDTOResponse, error)
	}

	topupTransactionService struct {
		topupTransactionRepository repository.TopupTransactionRepository
		db                         *gorm.DB
	}
)

func NewTopupTransactionService(topupTransactionRepository repository.TopupTransactionRepository, db *gorm.DB) TopupTransactionService {
	return &topupTransactionService{
		topupTransactionRepository: topupTransactionRepository,
		db:                         db,
	}
}

func (s *topupTransactionService) GetAll(ctx context.Context, filters repository.TopupTransactionFilters, limit, offset int) ([]dto_response.TopupTransactionDTOResponse, int64, error) {
	transactions, total, err := s.topupTransactionRepository.GetAll(ctx, nil, filters, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	var transactionDTOs []dto_response.TopupTransactionDTOResponse
	for i := range transactions {
		transactionDTOs = append(transactionDTOs, *s.mapToTopupTransactionDTOResponse(&transactions[i]))
	}
	return transactionDTOs, total, nil
}

func (s *topupTransactionService) GetByID(ctx context.Context, id string) (*dto_response.TopupTransactionDTOResponse, error) {
	parsedUserID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	transaction, err := s.topupTransactionRepository.GetByID(ctx, nil, parsedUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, myerror.RecordNotFound("topup transaction")
		}
		return nil, err
	}

	return s.mapToTopupTransactionDTOResponse(transaction), nil
}

func (s *topupTransactionService) mapToTopupTransactionDTOResponse(transaction *entity.TopupTransaction) *dto_response.TopupTransactionDTOResponse {
	var bundlePackageID string
	if transaction.BundlePackageID != nil {
		bundlePackageID = transaction.BundlePackageID.String()
	}

	var ledgerID string
	if transaction.LedgerID != nil {
		ledgerID = transaction.LedgerID.String()
	}

	var provider string
	if transaction.Provider != nil {
		provider = *transaction.Provider
	}

	var paidAt time.Time
	if transaction.PaidAt != nil {
		paidAt = *transaction.PaidAt
	}

	var expiredAt time.Time
	if transaction.ExpiredAt != nil {
		expiredAt = *transaction.ExpiredAt
	}

	return &dto_response.TopupTransactionDTOResponse{
		ID:              transaction.ID.String(),
		UserID:          transaction.UserID.String(),
		Type:            string(transaction.Type),
		BundlePackageID: bundlePackageID,
		TokenAmount:     transaction.TokenAmount,
		TotalPriceRp:    transaction.TotalPriceRp,
		Status:          string(transaction.Status),
		InvoiceID:       transaction.InvoiceID,
		Provider:        provider,
		PaidAt:          paidAt,
		ExpiredAt:       expiredAt,
		LedgerID:        ledgerID,
	}
}
