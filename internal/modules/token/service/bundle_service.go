package service

import (
	"context"
	"errors"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/modules/token/repository"
	myerror "rextra-backend/internal/pkg/error"

	"gorm.io/gorm"
)

type (
	BundleService interface {
		// FOR ALL ROLE
		GetAll(ctx context.Context) ([]dto_response.TokenBundleDTOResponse, error)
		GetByID(ctx context.Context, id string) (dto_response.TokenBundleDTOResponse, error)
	}

	bundleService struct {
		tokenBundleRepository repository.TokenBundlePackageRepository
		db                    *gorm.DB
	}
)

func NewBundleService(tokenBundleRepository repository.TokenBundlePackageRepository, db *gorm.DB) BundleService {
	return &bundleService{
		tokenBundleRepository: tokenBundleRepository,
		db:                    db,
	}
}

func (s *bundleService) GetAll(ctx context.Context) ([]dto_response.TokenBundleDTOResponse, error) {
	var tokenBundles []dto_response.TokenBundleDTOResponse

	bundles, err := s.tokenBundleRepository.GetAll(ctx, nil)
	if err != nil {
		return nil, err
	}

	for _, bundle := range bundles {
		var label string
		if bundle.Label != nil {
			label = *bundle.Label
		}
		tokenBundle := dto_response.TokenBundleDTOResponse{
			ID:          bundle.ID.String(),
			Name:        bundle.Name,
			TokenAmount: bundle.TokenAmount,
			PriceRp:     bundle.PriceRp,
			Label:       label,
			IsActive:    bundle.IsActive,
		}
		tokenBundles = append(tokenBundles, tokenBundle)
	}

	return tokenBundles, nil
}

func (s *bundleService) GetByID(ctx context.Context, id string) (dto_response.TokenBundleDTOResponse, error) {
	var tokenBundle dto_response.TokenBundleDTOResponse

	bundle, err := s.tokenBundleRepository.GetByID(ctx, nil, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tokenBundle, myerror.RecordNotFound("bundle")
		}
		return tokenBundle, err
	}

	var label string
	if bundle.Label != nil {
		label = *bundle.Label
	}

	tokenBundle = dto_response.TokenBundleDTOResponse{
		ID:          bundle.ID.String(),
		Name:        bundle.Name,
		TokenAmount: bundle.TokenAmount,
		PriceRp:     bundle.PriceRp,
		Label:       label,
		IsActive:    bundle.IsActive,
	}

	return tokenBundle, nil
}
