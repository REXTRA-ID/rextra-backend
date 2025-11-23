package service

import (
	"context"
	"rextra-backend/internal/api/repository"

	"gorm.io/gorm"
)

type (
	PoinTransactionService interface {
		UsePoin(ctx context.Context)
		EarnPoin(ctx context.Context)
		RedeemPoin(ctx context.Context)
	}

	poinTransactionService struct {
		membershipRepository      repository.MembershipRepository
		poinTransactionRepository repository.PoinTransactionsRepository
		db                        *gorm.DB
	}
)

func NewPoinTransactionService(
	membershipRepository repository.MembershipRepository,
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
