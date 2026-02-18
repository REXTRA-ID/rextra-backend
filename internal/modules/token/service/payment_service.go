package service

import (
	"rextra-backend/internal/modules/token/repository"
	"rextra-backend/internal/pkg/tripay"
)

type (
	PaymentService interface {
		GetInstructions(string) ([]tripay.InstructionStep, error)
	}

	paymentService struct {
		topUpRepository       repository.TopupTransactionRepository
		tokenLedgerRepository repository.TokenLedgerRepository
	}
)

func NewPaymentService(topUpRepository repository.TopupTransactionRepository, tokenLedgerRepository repository.TokenLedgerRepository) PaymentService {
	return &paymentService{
		topUpRepository:       topUpRepository,
		tokenLedgerRepository: tokenLedgerRepository,
	}
}

func (s *paymentService) GetInstructions(code string) ([]tripay.InstructionStep, error) {
	res, err := tripay.GetInstruction(code)
	if err != nil {
		return nil, err
	}
	return res, nil
}
