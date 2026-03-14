package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	"rextra-backend/internal/modules/entitlement/repository"
	myerror "rextra-backend/internal/pkg/error"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EntitlementService interface {
	Create(ctx context.Context, req dto_request.CreateEntitlementRequest) (dto_response.GetEntitlementResponse, error)
	GetAll(ctx context.Context, featureID *string) ([]dto_response.GetEntitlementResponse, error)
	GetById(ctx context.Context, id string) (dto_response.GetEntitlementResponse, error)
	UpdateRestriction(ctx context.Context, id string, req dto_request.UpdateEntitlementRestrictionRequest) (dto_response.GetEntitlementResponse, error)
	Delete(ctx context.Context, id string) error
}

type entitlementService struct {
	entitlementRepository repository.EntitlementRepository
	featureRepository     repository.FeatureRepository
	subFeatureRepository  repository.SubFeatureRepository
	actionCategoryRepo    repository.ActionCategoryRepository
}

func NewEntitlementService(
	entitlementRepo repository.EntitlementRepository,
	featureRepo repository.FeatureRepository,
	subFeatureRepo repository.SubFeatureRepository,
	actionCategoryRepo repository.ActionCategoryRepository,
) EntitlementService {
	return &entitlementService{
		entitlementRepository: entitlementRepo,
		featureRepository:     featureRepo,
		subFeatureRepository:  subFeatureRepo,
		actionCategoryRepo:    actionCategoryRepo,
	}
}

func (s *entitlementService) Create(ctx context.Context, req dto_request.CreateEntitlementRequest) (dto_response.GetEntitlementResponse, error) {
	parsedFeatureID, err := uuid.Parse(req.FeatureID)
	if err != nil {
		return dto_response.GetEntitlementResponse{}, myerror.InvalidRequest(err)
	}

	feature, err := s.featureRepository.GetById(ctx, nil, parsedFeatureID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.GetEntitlementResponse{}, myerror.RecordNotFound("feature")
		}
		return dto_response.GetEntitlementResponse{}, myerror.DatabaseError(err)
	}

	parsedActionCategoryID, err := uuid.Parse(req.ActionCategoryID)
	if err != nil {
		return dto_response.GetEntitlementResponse{}, myerror.InvalidRequest(err)
	}

	actionCategory, err := s.actionCategoryRepo.GetById(ctx, nil, parsedActionCategoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.GetEntitlementResponse{}, myerror.RecordNotFound("action category")
		}
		return dto_response.GetEntitlementResponse{}, myerror.DatabaseError(err)
	}

	var parsedSubFeatureID *uuid.UUID
	var subFeatureSlug string
	var subFeature *entity.SubFeature

	if entity.EntitlementLevel(req.Level) == entity.EntitlementLevelSubFeature {
		if req.SubFeatureID == nil || *req.SubFeatureID == "" {
			return dto_response.GetEntitlementResponse{}, myerror.New(
				"sub_feature_id is required when level is 'sub_fitur'",
				myerror.Error_InvalidRequest,
			)
		}
		parsedSFID, err := uuid.Parse(*req.SubFeatureID)
		if err != nil {
			return dto_response.GetEntitlementResponse{}, myerror.InvalidRequest(err)
		}
		parsedSubFeatureID = &parsedSFID

		sf, err := s.subFeatureRepository.GetById(ctx, nil, parsedSFID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return dto_response.GetEntitlementResponse{}, myerror.RecordNotFound("sub feature")
			}
			return dto_response.GetEntitlementResponse{}, myerror.DatabaseError(err)
		}
		if sf.FeatureID != parsedFeatureID {
			return dto_response.GetEntitlementResponse{}, myerror.New(
				"sub feature does not belong to the specified feature",
				myerror.Error_InvalidRequest,
			)
		}
		subFeatureSlug = sf.Slug
		subFeature = &sf

	} else {
		if req.SubFeatureID != nil && *req.SubFeatureID != "" {
			return dto_response.GetEntitlementResponse{}, myerror.New(
				"sub_feature_id must be empty when level is 'fitur'",
				myerror.Error_InvalidRequest,
			)
		}
	}

	if err := validateRestriction(req.RestrictionType, req.TokenCost); err != nil {
		return dto_response.GetEntitlementResponse{}, err
	}

	key := entity.BuildEntitlementKey(feature.Prefix, subFeatureSlug, actionCategory.Slug)

	_, err = s.entitlementRepository.GetByKey(ctx, nil, key)
	if err == nil {
		return dto_response.GetEntitlementResponse{}, myerror.RecordAlreadyExist(
			fmt.Sprintf("entitlement with key '%s'", key),
		)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto_response.GetEntitlementResponse{}, myerror.DatabaseError(err)
	}

	_, err = s.entitlementRepository.GetByComponents(ctx, nil, parsedFeatureID, parsedSubFeatureID, parsedActionCategoryID)
	if err == nil {
		return dto_response.GetEntitlementResponse{}, myerror.RecordAlreadyExist(
			"entitlement with same feature, sub feature, and action combination",
		)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto_response.GetEntitlementResponse{}, myerror.DatabaseError(err)
	}

	newEntitlement := entity.NewEntitlement(
		key,
		req.Name,
		req.Description,
		string(entity.RestrictionType(req.RestrictionType)),
		"",
		req.TokenCost,
		entity.EntitlementLevel(req.Level),
		parsedFeatureID,
		parsedSubFeatureID,
		parsedActionCategoryID,
	)

	result, err := s.entitlementRepository.Create(ctx, nil, newEntitlement)
	if err != nil {
		return dto_response.GetEntitlementResponse{}, myerror.DatabaseError(err)
	}

	result.Feature = feature
	result.SubFeature = subFeature
	result.ActionCategory = actionCategory

	return toEntitlementResponse(result, 0), nil
}

func (s *entitlementService) GetAll(ctx context.Context, featureID *string) ([]dto_response.GetEntitlementResponse, error) {
	var parsedFeatureID *uuid.UUID
	if featureID != nil && *featureID != "" {
		id, err := uuid.Parse(*featureID)
		if err != nil {
			return nil, myerror.InvalidRequest(err)
		}
		parsedFeatureID = &id
	}

	results, countMap, err := s.entitlementRepository.GetAllWithMappingCounts(ctx, nil, parsedFeatureID)
	if err != nil {
		return nil, myerror.DatabaseError(err)
	}

	res := make([]dto_response.GetEntitlementResponse, 0, len(results))
	for _, r := range results {
		res = append(res, toEntitlementResponse(r, int(countMap[r.ID])))
	}
	return res, nil
}

func (s *entitlementService) GetById(ctx context.Context, id string) (dto_response.GetEntitlementResponse, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return dto_response.GetEntitlementResponse{}, myerror.InvalidRequest(err)
	}

	result, err := s.entitlementRepository.GetById(ctx, nil, parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.GetEntitlementResponse{}, myerror.RecordNotFound("entitlement")
		}
		return dto_response.GetEntitlementResponse{}, myerror.DatabaseError(err)
	}

	count, _ := s.entitlementRepository.CountMappings(ctx, nil, result.ID)
	return toEntitlementResponse(result, int(count)), nil
}

func (s *entitlementService) UpdateRestriction(ctx context.Context, id string, req dto_request.UpdateEntitlementRestrictionRequest) (dto_response.GetEntitlementResponse, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return dto_response.GetEntitlementResponse{}, myerror.InvalidRequest(err)
	}

	existing, err := s.entitlementRepository.GetById(ctx, nil, parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.GetEntitlementResponse{}, myerror.RecordNotFound("entitlement")
		}
		return dto_response.GetEntitlementResponse{}, myerror.DatabaseError(err)
	}

	if err := validateRestriction(req.RestrictionType, req.TokenCost); err != nil {
		return dto_response.GetEntitlementResponse{}, err
	}

	existing.RestrictionType = entity.RestrictionType(req.RestrictionType)
	existing.TokenCost = req.TokenCost
	existing.UpdatedAt = time.Now().UTC()

	result, err := s.entitlementRepository.UpdateRestriction(ctx, nil, existing)
	if err != nil {
		return dto_response.GetEntitlementResponse{}, myerror.DatabaseError(err)
	}

	count, _ := s.entitlementRepository.CountMappings(ctx, nil, result.ID)
	return toEntitlementResponse(result, int(count)), nil
}

func (s *entitlementService) Delete(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return myerror.InvalidRequest(err)
	}

	_, err = s.entitlementRepository.GetById(ctx, nil, parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return myerror.RecordNotFound("entitlement")
		}
		return myerror.DatabaseError(err)
	}

	count, err := s.entitlementRepository.CountMappings(ctx, nil, parsedID)
	if err != nil {
		return myerror.DatabaseError(err)
	}
	if count > 0 {
		return myerror.New(
			fmt.Sprintf("entitlement cannot be deleted because it is used in %d plan mapping(s) — remove all mappings first", count),
			myerror.Error_InvalidRequest,
		)
	}

	if err := s.entitlementRepository.Delete(ctx, nil, parsedID); err != nil {
		return myerror.DatabaseError(err)
	}
	return nil
}

func validateRestriction(restrictionType string, tokenCost int) error {
	if restrictionType == string(entity.RestrictionTokenGated) && tokenCost <= 0 {
		return myerror.New(
			"token_cost must be > 0 when restriction_type is 'token_gated'",
			myerror.Error_InvalidRequest,
		)
	}
	return nil
}

func toEntitlementResponse(e entity.Entitlement, mappingCount int) dto_response.GetEntitlementResponse {
	createdAt, updatedAt := "", ""
	if !e.CreatedAt.IsZero() {
		createdAt = e.CreatedAt.Format(time.RFC3339)
	}
	if !e.UpdatedAt.IsZero() {
		updatedAt = e.UpdatedAt.Format(time.RFC3339)
	}

	res := dto_response.GetEntitlementResponse{
		Id:                 e.ID.String(),
		Key:                e.Key,
		Name:               e.Name,
		Description:        e.Description,
		Level:              string(e.Level),
		Status:             string(e.Status),
		MappingCount:       mappingCount,
		FeatureId:          e.FeatureID.String(),
		FeatureName:        e.Feature.Name,
		FeaturePrefix:      e.Feature.Prefix,
		ActionCategoryId:   e.ActionCategoryID.String(),
		ActionCategoryName: e.ActionCategory.Name,
		ActionCategorySlug: e.ActionCategory.Slug,
		RestrictionType:    string(e.RestrictionType),
		TokenCost:          e.TokenCost,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
	}

	if e.SubFeature != nil {
		sfID := e.SubFeature.ID.String()
		sfName := e.SubFeature.Name
		res.SubFeatureId = &sfID
		res.SubFeatureName = &sfName
	}

	return res
}
