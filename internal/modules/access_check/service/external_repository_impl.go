package service

import (
	"context"
	"fmt"

	"rextra-backend/internal/entity"
	tokenRepo "rextra-backend/internal/modules/token/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type tokenRepositoryAdapter struct {
	walletRepo tokenRepo.TokenWalletRepository
	ledgerRepo tokenRepo.TokenLedgerRepository
}

func NewTokenRepositoryAdapter(walletRepo tokenRepo.TokenWalletRepository, ledgerRepo tokenRepo.TokenLedgerRepository) TokenRepository {
	return &tokenRepositoryAdapter{walletRepo: walletRepo, ledgerRepo: ledgerRepo}
}

func (a *tokenRepositoryAdapter) GetWalletByUserID(ctx context.Context, tx *gorm.DB, userID uuid.UUID) (TokenWalletSnapshot, error) {
	wallet, err := a.walletRepo.GetByUserID(ctx, tx, userID.String())
	if err != nil { return TokenWalletSnapshot{}, err }
	return TokenWalletSnapshot{UserID: wallet.UserID, Balance: wallet.Balance}, nil
}

func (a *tokenRepositoryAdapter) DeductToken(ctx context.Context, tx *gorm.DB, userID uuid.UUID, amount int, description string) error {
	wallet, err := a.walletRepo.GetByUserID(ctx, tx, userID.String())
	if err != nil { return err }
	newBalance := wallet.Balance - int64(amount)
	if newBalance < 0 { return fmt.Errorf("saldo tidak cukup") }

	ledger := entity.TokenLedger{
		WalletID: wallet.ID, Direction: entity.DirectionOUT, Amount: int64(amount),
		BalanceBefore: wallet.Balance, BalanceAfter: newBalance,
		SourceType: entity.SourceTypeUsage, Description: description,
	}
	createdLedger, err := a.ledgerRepo.Create(ctx, tx, ledger)
	if err != nil { return err }

	return a.walletRepo.UpdateBalance(ctx, tx, userID.String(), float64(newBalance), createdLedger.ID)
}
