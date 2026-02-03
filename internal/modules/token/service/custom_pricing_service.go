package service

import (
	"context"
	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"

	"rextra-backend/internal/modules/token/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	CustomPricingService interface {
		GetCurrent(ctx context.Context) (dto_response.CustomPricingDTOResponse, error)
		CreateNewConfig(ctx context.Context, data dto_request.NewCustomPricingDTORequest, userID string) error
		GetHistory(ctx context.Context, limit int, tiers bool) ([]dto_response.CustomPricingDTOResponse, error)
		ToggleActive(ctx context.Context, active bool) error
	}
	customPricingService struct {
		customPricingRepository repository.CustomPricingRepository
		db                      *gorm.DB
	}
)

func NewCustomPricingService(customPricingRepository repository.CustomPricingRepository, db *gorm.DB) CustomPricingService {
	return &customPricingService{
		customPricingRepository: customPricingRepository,
		db:                      db,
	}
}

func (s *customPricingService) GetCurrent(ctx context.Context) (dto_response.CustomPricingDTOResponse, error) {
	customPricing, err := s.customPricingRepository.GetCurrentConfig(ctx, nil)
	if err != nil {
		return dto_response.CustomPricingDTOResponse{}, err
	}

	var tiers []dto_response.CustomPricingTier

	if customPricing.Tiers != nil {
		tiers = make([]dto_response.CustomPricingTier, len(customPricing.Tiers))
		for i, tier := range customPricing.Tiers {
			tiers[i] = dto_response.CustomPricingTier{
				ID:          tier.ID.String(),
				FromToken:   tier.FromToken,
				ToToken:     tier.ToToken,
				DiscountPct: tier.DiscountPct,
			}
		}
	}

	return dto_response.CustomPricingDTOResponse{
		ID:                       customPricing.ID.String(),
		IsEnabled:                customPricing.IsEnabled,
		Mintoken:                 customPricing.MinToken,
		Maxtoken:                 customPricing.MaxToken,
		RecommendedPricePerToken: customPricing.RecommendedPricePerToken,
		Tiers:                    tiers,
	}, nil
}

func (s *customPricingService) CreateNewConfig(ctx context.Context, data dto_request.NewCustomPricingDTORequest, userID string) error {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	customPricing := &entity.CustomPricingConfig{
		IsEnabled:                *data.IsEnable,
		MinToken:                 int64(data.MinToken),
		MaxToken:                 int64(data.MaxToken),
		RecommendedPricePerToken: int64(data.RecommendedPricePerToken),
		UpdatedBy:                &parsedUserID,
	}

	tiers := make([]entity.CustomPricingTier, len(data.Tiers))
	for i, tier := range data.Tiers {
		tiers[i] = entity.CustomPricingTier{
			FromToken:   int64(tier.FromToken),
			ToToken:     int64(tier.ToToken),
			DiscountPct: tier.DiscountPct,
		}
	}

	err = s.customPricingRepository.CreateNewVersion(ctx, nil, customPricing, tiers)

	if err != nil {
		return err
	}

	return nil
}

func (s *customPricingService) GetHistory(ctx context.Context, limit int, tiers bool) ([]dto_response.CustomPricingDTOResponse, error) {
	var preloads []string
	if tiers {
		preloads = append(preloads, "Tiers")
	}

	customPricings, err := s.customPricingRepository.GetConfigHistory(ctx, nil, limit, preloads...)

	if err != nil {
		return nil, err
	}

	dtos := make([]dto_response.CustomPricingDTOResponse, len(customPricings))
	for i, customPricing := range customPricings {
		var tiersDto []dto_response.CustomPricingTier

		if customPricing.Tiers != nil {
			tiersDto = make([]dto_response.CustomPricingTier, len(customPricing.Tiers))
			for j, tier := range customPricing.Tiers {
				tiersDto[j] = dto_response.CustomPricingTier{
					ID:          tier.ID.String(),
					FromToken:   tier.FromToken,
					ToToken:     tier.ToToken,
					DiscountPct: tier.DiscountPct,
				}
			}
		}

		dtos[i] = dto_response.CustomPricingDTOResponse{
			ID:                       customPricing.ID.String(),
			IsEnabled:                customPricing.IsEnabled,
			Mintoken:                 customPricing.MinToken,
			Maxtoken:                 customPricing.MaxToken,
			RecommendedPricePerToken: customPricing.RecommendedPricePerToken,
			Tiers:                    tiersDto,
		}
	}

	return dtos, nil
}

func (s *customPricingService) ToggleActive(ctx context.Context, active bool) error {
	return s.customPricingRepository.ToggleActive(ctx, nil, active)
}
