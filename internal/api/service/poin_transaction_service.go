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

func NewPoinTransactionService()
