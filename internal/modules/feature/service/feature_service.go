package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	"rextra-backend/internal/modules/feature/repository"
	myerror "rextra-backend/internal/pkg/error"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	FeatureService interface {
		Create(ctx context.Context, req dto_request.CreateFeatureRequest) (dto_response.GetFeatureResponse, error)
		GetAll(ctx context.Context) ([]dto_response.GetFeatureResponse, error)
		GetById(ctx context.Context, id string) (dto_response.GetFeatureResponse, error)
		Update(ctx context.Context, id string, req dto_request.UpdateFeatureRequest) (dto_response.GetFeatureResponse, error)
		Delete(ctx context.Context, id string) error
	}

	featureService struct {
		featureRepository repository.FeatureRepository
	}
)

func NewFeatureService(repo repository.FeatureRepository) FeatureService {
	return &featureService{featureRepository: repo}
}

func (s *featureService) Create(ctx context.Context, req dto_request.CreateFeatureRequest) (dto_response.GetFeatureResponse, error) {
	_, err := s.featureRepository.GetBySlug(ctx, nil, req.Slug)
	if err == nil {
		return dto_response.GetFeatureResponse{}, myerror.RecordAlreadyExist("feature with slug '" + req.Slug + "'")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto_response.GetFeatureResponse{}, myerror.DatabaseError(err)
	}

	_, err = s.featureRepository.GetByPrefix(ctx, nil, req.Prefix)
	if err == nil {
		return dto_response.GetFeatureResponse{}, myerror.RecordAlreadyExist("feature with prefix '" + req.Prefix + "'")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto_response.GetFeatureResponse{}, myerror.DatabaseError(err)
	}

	newFeature := entity.NewFeature(req.Name, req.Slug, req.Prefix, req.Description, entity.FeatureType(req.Type))
	result, err := s.featureRepository.Create(ctx, nil, newFeature)
	if err != nil {
		return dto_response.GetFeatureResponse{}, myerror.DatabaseError(err)
	}

	return toFeatureResponse(result, 0), nil
}

func (s *featureService) GetAll(ctx context.Context) ([]dto_response.GetFeatureResponse, error) {
	results, featureCountMap, subFeatureCountMap, err := s.featureRepository.GetAllWithCounts(ctx, nil)
	if err != nil {
		return nil, myerror.DatabaseError(err)
	}

	var res []dto_response.GetFeatureResponse
	for _, r := range results {
		resp := toFeatureResponse(r, int(featureCountMap[r.ID]))
		for _, sf := range r.SubFeatures {
			resp.SubFeatures = append(resp.SubFeatures, toSubFeatureResponse(sf, int(subFeatureCountMap[sf.ID]), r.Name))
		}
		res = append(res, resp)
	}
	return res, nil
}

func (s *featureService) GetById(ctx context.Context, id string) (dto_response.GetFeatureResponse, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return dto_response.GetFeatureResponse{}, myerror.InvalidRequest(err)
	}

	result, err := s.featureRepository.GetById(ctx, nil, parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.GetFeatureResponse{}, myerror.RecordNotFound("feature")
		}
		return dto_response.GetFeatureResponse{}, myerror.DatabaseError(err)
	}

	count, _ := s.featureRepository.CountEntitlements(ctx, nil, result.ID)
	resp := toFeatureResponse(result, int(count))
	for _, sf := range result.SubFeatures {
		resp.SubFeatures = append(resp.SubFeatures, toSubFeatureResponse(sf, 0, result.Name))
	}
	return resp, nil
}

func (s *featureService) Update(ctx context.Context, id string, req dto_request.UpdateFeatureRequest) (dto_response.GetFeatureResponse, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return dto_response.GetFeatureResponse{}, myerror.InvalidRequest(err)
	}

	existing, err := s.featureRepository.GetById(ctx, nil, parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.GetFeatureResponse{}, myerror.RecordNotFound("feature")
		}
		return dto_response.GetFeatureResponse{}, myerror.DatabaseError(err)
	}

	slugChanged := req.Slug != existing.Slug
	prefixChanged := req.Prefix != existing.Prefix

	if slugChanged || prefixChanged {
		count, err := s.featureRepository.CountEntitlements(ctx, nil, parsedID)
		if err != nil {
			return dto_response.GetFeatureResponse{}, myerror.DatabaseError(err)
		}
		if count > 0 {
			if slugChanged && prefixChanged {
				return dto_response.GetFeatureResponse{}, myerror.New(
					fmt.Sprintf("slug and prefix cannot be changed because this feature has %d entitlement(s)", count),
					myerror.Error_InvalidRequest,
				)
			}
			if slugChanged {
				return dto_response.GetFeatureResponse{}, myerror.New(
					fmt.Sprintf("slug cannot be changed because this feature has %d entitlement(s)", count),
					myerror.Error_InvalidRequest,
				)
			}
			return dto_response.GetFeatureResponse{}, myerror.New(
				fmt.Sprintf("prefix cannot be changed because this feature has %d entitlement(s)", count),
				myerror.Error_InvalidRequest,
			)
		}

		if slugChanged {
			_, err = s.featureRepository.GetBySlug(ctx, nil, req.Slug)
			if err == nil {
				return dto_response.GetFeatureResponse{}, myerror.RecordAlreadyExist("feature with slug '" + req.Slug + "'")
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return dto_response.GetFeatureResponse{}, myerror.DatabaseError(err)
			}
		}
		if prefixChanged {
			_, err = s.featureRepository.GetByPrefix(ctx, nil, req.Prefix)
			if err == nil {
				return dto_response.GetFeatureResponse{}, myerror.RecordAlreadyExist("feature with prefix '" + req.Prefix + "'")
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return dto_response.GetFeatureResponse{}, myerror.DatabaseError(err)
			}
		}
	}

	existing.Name = req.Name
	existing.Slug = req.Slug
	existing.Prefix = req.Prefix
	existing.Description = req.Description
	existing.Status = entity.FeatureStatus(req.Status)
	existing.UpdatedAt = time.Now().UTC()

	result, err := s.featureRepository.Update(ctx, nil, existing)
	if err != nil {
		return dto_response.GetFeatureResponse{}, myerror.DatabaseError(err)
	}

	count, _ := s.featureRepository.CountEntitlements(ctx, nil, result.ID)
	return toFeatureResponse(result, int(count)), nil
}

func (s *featureService) Delete(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return myerror.InvalidRequest(err)
	}

	_, err = s.featureRepository.GetById(ctx, nil, parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return myerror.RecordNotFound("feature")
		}
		return myerror.DatabaseError(err)
	}

	entCount, err := s.featureRepository.CountEntitlements(ctx, nil, parsedID)
	if err != nil {
		return myerror.DatabaseError(err)
	}
	if entCount > 0 {
		return myerror.New(
			fmt.Sprintf("feature cannot be deleted because it has %d entitlement(s)", entCount),
			myerror.Error_InvalidRequest,
		)
	}

	sfCount, err := s.featureRepository.CountSubFeatures(ctx, nil, parsedID)
	if err != nil {
		return myerror.DatabaseError(err)
	}
	if sfCount > 0 {
		return myerror.New(
			fmt.Sprintf("feature cannot be deleted because it has %d sub feature(s)", sfCount),
			myerror.Error_InvalidRequest,
		)
	}

	if err := s.featureRepository.Delete(ctx, nil, parsedID); err != nil {
		return myerror.DatabaseError(err)
	}
	return nil
}

func toFeatureResponse(f entity.Feature, entitlementCount int) dto_response.GetFeatureResponse {
	createdAt, updatedAt := "", ""
	if !f.CreatedAt.IsZero() {
		createdAt = f.CreatedAt.Format(time.RFC3339)
	}
	if !f.UpdatedAt.IsZero() {
		updatedAt = f.UpdatedAt.Format(time.RFC3339)
	}
	return dto_response.GetFeatureResponse{
		Id:               f.ID.String(),
		Name:             f.Name,
		Slug:             f.Slug,
		Prefix:           f.Prefix,
		Description:      f.Description,
		Type:             string(f.Type),
		Status:           string(f.Status),
		EntitlementCount: entitlementCount,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}
}

func toSubFeatureResponse(sf entity.SubFeature, entitlementCount int, featureName string) dto_response.GetSubFeatureResponse {
	createdAt, updatedAt := "", ""
	if !sf.CreatedAt.IsZero() {
		createdAt = sf.CreatedAt.Format(time.RFC3339)
	}
	if !sf.UpdatedAt.IsZero() {
		updatedAt = sf.UpdatedAt.Format(time.RFC3339)
	}
	return dto_response.GetSubFeatureResponse{
		Id:               sf.ID.String(),
		FeatureId:        sf.FeatureID.String(),
		FeatureName:      featureName,
		Name:             sf.Name,
		Slug:             sf.Slug,
		Description:      sf.Description,
		Status:           string(sf.Status),
		EntitlementCount: entitlementCount,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}
}
