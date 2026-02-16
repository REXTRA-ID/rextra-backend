package service

import (
	"context"
	membershipRepo "rextra-backend/internal/modules/membership/repository"
	"rextra-backend/internal/modules/poin/repository"

	"gorm.io/gorm"
)

type (
	PoinTransactionService interface {
		UsePoin(ctx context.Context)
		EarnPoin(ctx context.Context)
		RedeemPoin(ctx context.Context)
	}

	poinTransactionService struct {
		membershipRepository      membershipRepo.MembershipRepository
		poinTransactionRepository repository.PoinTransactionsRepository
		db                        *gorm.DB
	}
)

func NewPoinTransactionService(
	membershipRepository membershipRepo.MembershipRepository,
	poinTransactionRepository repository.PoinTransactionsRepository,
	db *gorm.DB,
) PoinTransactionService {
	return &poinTransactionService{
		poinTransactionRepository: poinTransactionRepository,
		membershipRepository:      membershipRepository,
		db:                        db,
	}
}

func (s *poinTransactionService) UsePoin(ctx context.Context) {

}

func (s *poinTransactionService) EarnPoin(ctx context.Context) {

}

func (s *poinTransactionService) RedeemPoin(ctx context.Context) {

}
