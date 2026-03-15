package access_check

import (
	"rextra-backend/internal/modules/access_check/service"
	tokenRepository "rextra-backend/internal/modules/token/repository"

	"gorm.io/gorm"
)

func NewAccessCheckService(
	db *gorm.DB,
	tokenWalletRepo tokenRepository.TokenWalletRepository,
	tokenLedgerRepo tokenRepository.TokenLedgerRepository,
) service.AccessCheckService {
	repo := service.NewAccessCheckRepository(db)
	tokenRepo := service.NewTokenRepositoryAdapter(tokenWalletRepo, tokenLedgerRepo)
	return service.NewAccessCheckService(repo, tokenRepo, db)
}
