package service

import (
	"context"
	"errors"
	"rextra-backend/internal/api/repository"
	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	myerror "rextra-backend/internal/pkg/error"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	TokenTransactionService interface {
		UseToken(ctx context.Context, req dto_request.UseTokenRequest, userId string) (dto_response.UseTokenResponse, error)
		RefillToken(ctx context.Context, userId uuid.UUID) (dto_response.UseTokenResponse, error)
	}

	tokenTransactionService struct {
		membershipRepository        repository.MembershipRepository
		tokenTransactionRepository  repository.TokenTransactionRepository
		tokenUsageHistoryRepository repository.TokenUsageHistoryRepository
		db                          *gorm.DB
	}
)

func NewTokenTransactionService(
	membershipRepository repository.MembershipRepository,
	tokenTransactionReposiroty repository.TokenTransactionRepository,
	tokenUsageHistoryRepository repository.TokenUsageHistoryRepository,
	db *gorm.DB,
) TokenTransactionService {
	return &tokenTransactionService{
		membershipRepository:        membershipRepository,
		tokenTransactionRepository:  tokenTransactionReposiroty,
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
		if err != nil {
			return err
		}

		if req.TokenRequired > userMembership.GetTokenBalance() {
			return myerror.New("not enough token", myerror.Error_InvalidRequest)
		}

		metadata, err := collectMetaData(req)
		if err != nil {
			return err
		}

		totalTokenLeft := userMembership.CurrentTokenBalance - req.TokenRequired

		userMembership.UseMembershipToken(req.TokenRequired)

		newTokenTransaction := entity.NewTokenTransaction(
			uuidUserID,
			userMembership.ID,
			string(entity.TOKENUSAGE),
			userMembership.CurrentTokenBalance,
			req.TokenRequired,
			"",
		)

		newTokenUsageHistory := entity.NewTokenUsageHistory(
			uuidUserID,
			newTokenTransaction.ID,
			string(req.FeatureName),
			metadata,
		)

		_, err = s.membershipRepository.Update(ctx, tx, userMembership)
		if err != nil {
			return err
		}

		_, err = s.tokenTransactionRepository.Create(ctx, tx, newTokenTransaction)
		if err != nil {
			return err
		}

		_, err = s.tokenUsageHistoryRepository.Create(ctx, tx, newTokenUsageHistory)
		if err != nil {
			return err
		}

		returnValue.TotalToken = totalTokenLeft
		return nil
	})

	if err != nil {
		return dto_response.UseTokenResponse{}, err
	}

	return returnValue, nil
}

func (s *tokenTransactionService) RefillToken(ctx context.Context, userId uuid.UUID) (dto_response.UseTokenResponse, error) {
	return dto_response.UseTokenResponse{}, nil
}

func collectMetaData(req dto_request.UseTokenRequest) ([]byte, error) {
	if req.FeatureName == entity.KENALIDIRI && req.UsageMetaData.AssessmentMetadata != nil {
		return dto_request.ConvertAssessmentToBytes(req.UsageMetaData.AssessmentMetadata)
	} else if req.FeatureName == entity.CVGENERATOR && req.UsageMetaData.CVGenerator != nil {
		return dto_request.ConvertCVGeneratorToBytes(req.UsageMetaData.CVGenerator)
	} else if req.FeatureName == entity.AIINTERVIEWER && req.UsageMetaData.AiInterviewer != nil {
		return dto_request.ConvertAiInterviewerToBytes(req.UsageMetaData.AiInterviewer)
	} else {
		return []byte{}, myerror.InvalidRequest(errors.New("request invalid"))
	}
}
