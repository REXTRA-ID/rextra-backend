package service

import (
	"context"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/modules/token/repository"

	"gorm.io/gorm"
)

type (
	TokenLedgerService interface {
		GetActivity(ctx context.Context, filters repository.TokenLedgerFilter, limit int, offset int) ([]dto_response.TokenLedgerActivityDTOResponse, int64, error)
		// GetKPISummary(ctx context.Context, filters repository.TokenLedgerFilter) (dto_response.TokenLedgerSummaryDTOResponse, error)
	}

	tokenLedgerService struct {
		tokenLedgerRepository repository.TokenLedgerRepository
		db                    *gorm.DB
	}
)

func NewTokenLedgerService(tokenLedgerRepository repository.TokenLedgerRepository, db *gorm.DB) TokenLedgerService {
	return &tokenLedgerService{
		tokenLedgerRepository: tokenLedgerRepository,
		db:                    db,
	}
}

func (s *tokenLedgerService) GetActivity(ctx context.Context, filters repository.TokenLedgerFilter, limit int, offset int) ([]dto_response.TokenLedgerActivityDTOResponse, int64, error) {
	ledgers, total, err := s.tokenLedgerRepository.GetAllActivity(ctx, nil, filters, limit, offset, "Wallet.User")

	if err != nil {
		return nil, 0, err
	}

	var ledgerResponse []dto_response.TokenLedgerActivityDTOResponse

	for _, ledger := range ledgers {
		ledgerResponse = append(ledgerResponse, dto_response.TokenLedgerActivityDTOResponse{
			ID:            ledger.ID.String(),
			Username:      ledger.Wallet.User.Fullname,
			SourceType:    string(ledger.SourceType),
			BalanceBefore: ledger.BalanceBefore,
			BalanceAfter:  ledger.BalanceAfter,
			CreatedAt:     ledger.CreatedAt,
		})
	}

	return ledgerResponse, total, nil
}
