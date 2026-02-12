package service

import (
	"context"
	"errors"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/modules/token/repository"
	myerror "rextra-backend/internal/pkg/error"

	"gorm.io/gorm"
)

type (
	WalletService interface {
		// USER
		GetMyBalance(ctx context.Context, userID string) (dto_response.TokenWalletDTOResponse, error)
		GetMyHistory(ctx context.Context, userID string, page int, size int) ([]dto_response.TokenLedgerDTOResponse, dto_response.PaginationMeta, error)
		GetMyHistoryDetail(ctx context.Context, userID string, id string) (dto_response.TokenLedgerDetailDTOResponse, error)
	}

	walletService struct {
		tokenWalletRepository repository.TokenWalletRepository
		tokenLedgerRepository repository.TokenLedgerRepository
		db                    *gorm.DB
	}
)

func NewWalletService(tokenWalletRepository repository.TokenWalletRepository, tokenLedgerRepository repository.TokenLedgerRepository, db *gorm.DB) WalletService {
	return &walletService{
		tokenWalletRepository: tokenWalletRepository,
		tokenLedgerRepository: tokenLedgerRepository,
		db:                    db,
	}
}

func (s *walletService) GetMyBalance(ctx context.Context, userID string) (dto_response.TokenWalletDTOResponse, error) {
	tokenWallet, err := s.tokenWalletRepository.GetByUserID(ctx, nil, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.TokenWalletDTOResponse{
				UserId:  userID,
				Balance: 0,
			}, nil
		}
		return dto_response.TokenWalletDTOResponse{}, err
	}

	return dto_response.TokenWalletDTOResponse{
		UserId:  userID,
		Balance: tokenWallet.Balance,
	}, nil
}

func (s *walletService) GetMyHistory(ctx context.Context, userID string, page int, size int) ([]dto_response.TokenLedgerDTOResponse, dto_response.PaginationMeta, error) {
	tokenWallet, err := s.tokenWalletRepository.GetByUserID(ctx, nil, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []dto_response.TokenLedgerDTOResponse{}, dto_response.PaginationMeta{}, nil
		}
		return nil, dto_response.PaginationMeta{}, err
	}

	offset := (page - 1) * size
	tokenLedgers, err := s.tokenLedgerRepository.GetWalletHistory(ctx, nil, tokenWallet.ID.String(), size, offset)

	if err != nil {
		return nil, dto_response.PaginationMeta{}, err
	}

	count, err := s.tokenLedgerRepository.CountByWalletID(ctx, nil, tokenWallet.ID.String())
	if err != nil {
		return nil, dto_response.PaginationMeta{}, err
	}

	var tokenLedgerDTOResponses []dto_response.TokenLedgerDTOResponse

	for _, tokenLedger := range tokenLedgers {
		tokenLedgerDTOResponses = append(tokenLedgerDTOResponses, dto_response.TokenLedgerDTOResponse{
			ID:         tokenLedger.ID.String(),
			OccurredAt: tokenLedger.OccurredAt,
			Amount:     tokenLedger.Amount,
			Direction:  string(tokenLedger.Direction),
			SourceType: string(tokenLedger.SourceType),
		})
	}

	return tokenLedgerDTOResponses, dto_response.PaginationMeta{
		CurrentPage:  page,
		PerPage:      size,
		TotalRecords: int64(count),
		TotalPages:   (count + int(size) - 1) / int(size),
	}, nil
}

func (s *walletService) GetMyHistoryDetail(ctx context.Context, userID string, id string) (dto_response.TokenLedgerDetailDTOResponse, error) {
	tokenWallet, err := s.tokenWalletRepository.GetByUserID(ctx, nil, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.TokenLedgerDetailDTOResponse{}, myerror.RecordNotFound("wallet")
		}
		return dto_response.TokenLedgerDetailDTOResponse{}, err
	}

	tokenLedger, err := s.tokenLedgerRepository.GetDetail(ctx, nil, tokenWallet.ID.String(), id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.TokenLedgerDetailDTOResponse{}, myerror.RecordNotFound("history/ledger")
		}
		return dto_response.TokenLedgerDetailDTOResponse{}, err
	}

	return dto_response.TokenLedgerDetailDTOResponse{
		ID:            tokenLedger.ID.String(),
		OccurredAt:    tokenLedger.OccurredAt,
		Amount:        tokenLedger.Amount,
		Direction:     string(tokenLedger.Direction),
		BalanceBefore: tokenLedger.BalanceBefore,
		BalanceAfter:  tokenLedger.BalanceAfter,
		SourceType:    string(tokenLedger.SourceType),
		SourceID:      tokenLedger.SourceID,
		// ReferenceID:   tokenLedger.ReferenceID,
		Description: tokenLedger.Description,
		Metadata:    tokenLedger.Metadata,
		OperatorID:  tokenLedger.OperatorID,
		CreatedAt:   tokenLedger.CreatedAt,
	}, nil
}
