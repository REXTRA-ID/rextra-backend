package service

import (
	"context"
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
		UseTokenAssessment(ctx context.Context, req dto_request.UseTokenAssessmentRequest, userId string) (dto_response.UseTokenResponse, error)
		UseTokenCVGenerator(ctx context.Context, req dto_request.UseTokenCVGeneratorRequest, userId string) (dto_response.UseTokenResponse, error)
		UseTokenAiInterviewer(ctx context.Context, req dto_request.UseTokenAiInteweviewerRequest, userId string) (dto_response.UseTokenResponse, error)
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

func (s *tokenTransactionService) UseTokenAssessment(
	ctx context.Context,
	req dto_request.UseTokenAssessmentRequest,
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

		metadata, err := dto_request.ConvertAssessmentToBytes(&req.UsageMetaData)
		if err != nil {
			return myerror.ProcessingError(err)
		}

		totalTokenLeft := userMembership.CurrentTokenBalance - req.TokenRequired

		newTokenTransaction := entity.NewTokenTransaction(
			uuidUserID,
			&userMembership,
			string(entity.TOKENUSAGE),
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

func (s *tokenTransactionService) UseTokenCVGenerator(
	ctx context.Context,
	req dto_request.UseTokenCVGeneratorRequest,
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

		metadata, err := dto_request.ConvertCVGeneratorToBytes(&req.UsageMetaData)
		if err != nil {
			return myerror.ProcessingError(err)
		}

		totalTokenLeft := userMembership.CurrentTokenBalance - req.TokenRequired

		newTokenTransaction := entity.NewTokenTransaction(
			uuidUserID,
			&userMembership,
			string(entity.TOKENUSAGE),
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

func (s *tokenTransactionService) UseTokenAiInterviewer(
	ctx context.Context,
	req dto_request.UseTokenAiInteweviewerRequest,
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

		metadata, err := dto_request.ConvertAiInterviewerToBytes(&req.UsageMetaData)
		if err != nil {
			return myerror.ProcessingError(err)
		}

		totalTokenLeft := userMembership.CurrentTokenBalance - req.TokenRequired

		newTokenTransaction := entity.NewTokenTransaction(
			uuidUserID,
			&userMembership,
			string(entity.TOKENUSAGE),
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
