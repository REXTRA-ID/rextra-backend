package service

import (
	"context"
	"errors"
	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	"rextra-backend/internal/modules/token/repository"
	myerror "rextra-backend/internal/pkg/error"

	"gorm.io/gorm"
)

type (
	BundleService interface {
		// FOR ALL ROLE
		GetAll(ctx context.Context) ([]dto_response.TokenBundleDTOResponse, error)
		GetByID(ctx context.Context, id string) (dto_response.TokenBundleDTOResponse, error)

		// FOR ADMIN ROLE
		Create(ctx context.Context, data dto_request.TokenBundleDTORequest) (dto_response.TokenBundleDTOResponse, error)
		Update(ctx context.Context, id string, data dto_request.TokenBundleDTORequest) (dto_response.TokenBundleDTOResponse, error)
		Delete(ctx context.Context, id string) error
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

func (s *bundleService) Create(ctx context.Context, data dto_request.TokenBundleDTORequest) (dto_response.TokenBundleDTOResponse, error) {
	var tokenBundle dto_response.TokenBundleDTOResponse

	_, isExist, err := s.tokenBundleRepository.GetByName(ctx, nil, data.Name)
	if isExist {
		return tokenBundle, myerror.New("Bundle Name Already Exist", myerror.Error_InvalidRequest)
	}
	if err != nil {
		return tokenBundle, err
	}

	bundle := entity.TokenBundlePackage{
		Name:        data.Name,
		TokenAmount: int64(data.TokenAmount),
		PriceRp:     int64(data.PriceRp),
		Label:       &data.Label,
		IsActive:    *data.IsActive,
	}

	createBundle, err := s.tokenBundleRepository.Create(ctx, nil, bundle)
	if err != nil {
		return tokenBundle, err
	}

	var label string
	if createBundle.Label != nil {
		label = *createBundle.Label
	}

	tokenBundle = dto_response.TokenBundleDTOResponse{
		ID:          createBundle.ID.String(),
		Name:        createBundle.Name,
		TokenAmount: createBundle.TokenAmount,
		PriceRp:     createBundle.PriceRp,
		Label:       label,
		IsActive:    createBundle.IsActive,
	}

	return tokenBundle, nil
}

func (s *bundleService) Update(ctx context.Context, id string, data dto_request.TokenBundleDTORequest) (dto_response.TokenBundleDTOResponse, error) {

	existingBundle, err := s.tokenBundleRepository.GetByID(ctx, nil, id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.TokenBundleDTOResponse{}, myerror.RecordNotFound(id)
		}
		return dto_response.TokenBundleDTOResponse{}, err
	}

	if existingBundle.Name != data.Name {
		duplicateCheck, isExist, err := s.tokenBundleRepository.GetByName(ctx, nil, data.Name)
		if err != nil {
			return dto_response.TokenBundleDTOResponse{}, err
		}
		if isExist && duplicateCheck.ID != existingBundle.ID {
			return dto_response.TokenBundleDTOResponse{}, myerror.New("Bundle Name Already Exist", myerror.Error_InvalidRequest)
		}

	}

	existingBundle.Name = data.Name
	existingBundle.TokenAmount = int64(data.TokenAmount)
	existingBundle.PriceRp = int64(data.PriceRp)
	existingBundle.Label = &data.Label
	existingBundle.IsActive = *data.IsActive

	updatedBundle, err := s.tokenBundleRepository.Update(ctx, nil, existingBundle)
	if err != nil {
		return dto_response.TokenBundleDTOResponse{}, err
	}

	var label string
	if updatedBundle.Label != nil {
		label = *updatedBundle.Label
	}

	response := dto_response.TokenBundleDTOResponse{
		ID:          updatedBundle.ID.String(),
		Name:        updatedBundle.Name,
		TokenAmount: updatedBundle.TokenAmount,
		PriceRp:     updatedBundle.PriceRp,
		Label:       label,
		IsActive:    updatedBundle.IsActive,
	}

	return response, nil
}

func (s *bundleService) Delete(ctx context.Context, id string) error {
	_, err := s.tokenBundleRepository.GetByID(ctx, nil, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return myerror.RecordNotFound("bundle")
		}
		return err
	}

	return s.tokenBundleRepository.Delete(ctx, nil, id)
}
