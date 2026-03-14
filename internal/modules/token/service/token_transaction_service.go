package service

import (
	"context"
	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	membershipRepo "rextra-backend/internal/modules/membership/repository"
	"rextra-backend/internal/modules/token/repository"
	myerror "rextra-backend/internal/pkg/error"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	TokenTransactionService interface {
		UseToken(ctx context.Context, req dto_request.UseTokenRequest, userId string) (dto_response.UseTokenResponse, error)
		RefillToken() error
	}

	tokenTransactionService struct {
		membershipRepository        membershipRepo.UserMembershipRepository
		tokenTransactionRepository  repository.TokenTransactionRepository
		tokenUsageHistoryRepository repository.TokenUsageHistoryRepository
		db                          *gorm.DB
	}
)

func NewTokenTransactionService(
	membershipRepository membershipRepo.UserMembershipRepository,
	tokenTransactionRepository repository.TokenTransactionRepository,
	tokenUsageHistoryRepository repository.TokenUsageHistoryRepository,
	db *gorm.DB,
) TokenTransactionService {
	return &tokenTransactionService{
		membershipRepository:        membershipRepository,
		tokenTransactionRepository:  tokenTransactionRepository,
		tokenUsageHistoryRepository: tokenUsageHistoryRepository,
		db:                          db,
	}
}

func (s *tokenTransactionService) UseToken(
	ctx context.Context,
	req dto_request.UseTokenRequest,
	userId string,
) (dto_response.UseTokenResponse, error) {
	uuidUserID := uuid.MustParse(userId)
	returnValue := dto_response.UseTokenResponse{}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		userMembership, err := s.membershipRepository.GetByUserID(ctx, tx, uuidUserID)
		if err != nil { return err }

		if req.TokenRequired > userMembership.CurrentTokenBalance {
			return myerror.New("not enough token", myerror.Error_InvalidRequest)
		}

		totalTokenLeft := userMembership.CurrentTokenBalance - req.TokenRequired
		userMembership.UseMembershipToken(req.TokenRequired)

		newTokenTransaction := entity.NewTokenTransaction(
			uuidUserID, userMembership.ID, string(entity.TOKENUSAGE),
			userMembership.CurrentTokenBalance, req.TokenRequired, "",
		)

		if err := tx.Create(&newTokenTransaction).Error; err != nil { return err }

		balanceBefore := int64(userMembership.CurrentTokenBalance + req.TokenRequired)
		balanceAfter := int64(userMembership.CurrentTokenBalance)
		
		newTokenUsageHistory := entity.NewTokenUsageHistory(
			uuidUserID, &userMembership.ID, req.FeatureName, req.FeatureName,
			req.TokenRequired, balanceBefore, balanceAfter, nil, nil,
		)

		if _, err := s.membershipRepository.Update(ctx, tx, userMembership); err != nil { return err }
		if _, err := s.tokenUsageHistoryRepository.Create(ctx, tx, newTokenUsageHistory); err != nil { return err }

		returnValue.TotalToken = totalTokenLeft
		return nil
	})

	if err != nil { return dto_response.UseTokenResponse{}, err }
	return returnValue, nil
}

func (s *tokenTransactionService) RefillToken() error {
	return nil
}
