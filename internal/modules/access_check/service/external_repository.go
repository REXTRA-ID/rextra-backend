package service

import (
	"context"
	"fmt"

	entity "rextra-backend/internal/entity"
	tokenRepo "rextra-backend/internal/modules/token/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TokenRepository adalah interface adapter ke modul token.
// Implementasinya ada di modul token — di-inject saat InitModule.
type TokenRepository interface {
	// GetWalletByUserID mengambil token wallet user.
	// Return gorm.ErrRecordNotFound jika wallet belum ada.
	GetWalletByUserID(ctx context.Context, tx *gorm.DB, userID uuid.UUID) (TokenWalletSnapshot, error)

	// DeductToken mengurangi saldo wallet dan insert TokenLedger dalam satu transaksi.
	// Wajib dipanggil dalam DB transaction yang sama dengan InsertUsageLog.
	DeductToken(ctx context.Context, tx *gorm.DB, userID uuid.UUID, amount int, description string) error
}

// TokenWalletSnapshot adalah data minimal wallet yang dibutuhkan access-check.
// Tidak import entity token langsung untuk menghindari circular dependency.
type TokenWalletSnapshot struct {
	UserID  uuid.UUID
	Balance int64
}

// tokenRepositoryAdapter mengadaptasi token repository ke interface yang dibutuhkan access-check.
type tokenRepositoryAdapter struct {
	walletRepo tokenRepo.TokenWalletRepository
	ledgerRepo tokenRepo.TokenLedgerRepository
}

func NewTokenRepositoryAdapter(
	walletRepo tokenRepo.TokenWalletRepository,
	ledgerRepo tokenRepo.TokenLedgerRepository,
) TokenRepository {
	return &tokenRepositoryAdapter{
		walletRepo: walletRepo,
		ledgerRepo: ledgerRepo,
	}
}

func (a *tokenRepositoryAdapter) GetWalletByUserID(ctx context.Context, tx *gorm.DB, userID uuid.UUID) (TokenWalletSnapshot, error) {
	wallet, err := a.walletRepo.GetByUserID(ctx, tx, userID.String())
	if err != nil {
		return TokenWalletSnapshot{}, err
	}
	return TokenWalletSnapshot{
		UserID:  wallet.UserID,
		Balance: wallet.Balance,
	}, nil
}

func (a *tokenRepositoryAdapter) DeductToken(ctx context.Context, tx *gorm.DB, userID uuid.UUID, amount int, description string) error {
	// Ambil wallet untuk baca balance sekarang
	wallet, err := a.walletRepo.GetByUserID(ctx, tx, userID.String())
	if err != nil {
		return err
	}

	newBalance := wallet.Balance - int64(amount)
	if newBalance < 0 {
		return fmt.Errorf("saldo tidak cukup")
	}

	// Insert TokenLedger
	ledger := entity.TokenLedger{
		WalletID:      wallet.ID,
		Direction:     entity.DirectionOUT,
		Amount:        int64(amount),
		BalanceBefore: wallet.Balance,
		BalanceAfter:  newBalance,
		SourceType:    entity.SourceTypeUsage,
		Description:   description,
		Metadata:      entity.TokenLedgerMetadata{},
	}
	createdLedger, err := a.ledgerRepo.Create(ctx, tx, ledger)
	if err != nil {
		return err
	}

	// Update wallet balance
	// [FIX v2] Hilangkan float64() cast — newBalance adalah int64, TokenWallet.Balance adalah int64.
	// Jika signature UpdateBalance di Hamzah menerima float64, koordinasikan untuk diubah ke int64
	// agar konsisten dengan tipe storage. Jangan cast int64 → float64 karena bisa kehilangan presisi.
	return a.walletRepo.UpdateBalance(ctx, tx, userID.String(), newBalance, createdLedger.ID)
}
