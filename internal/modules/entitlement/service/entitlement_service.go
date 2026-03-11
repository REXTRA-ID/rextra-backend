package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	ac "rextra-backend/internal/modules/action_category/repository"
	e "rextra-backend/internal/modules/entitlement/repository"
	f "rextra-backend/internal/modules/feature/repository"
	sf "rextra-backend/internal/modules/sub_feature/repository"
	myerror "rextra-backend/internal/pkg/error"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	EntitlementService interface {
		Create(ctx context.Context, req dto_request.CreateEntitlementRequest) (dto_response.GetEntitlementResponse, error)
		GetAll(ctx context.Context, featureID *string) ([]dto_response.GetEntitlementResponse, error)
		GetById(ctx context.Context, id string) (dto_response.GetEntitlementResponse, error)
		Delete(ctx context.Context, id string) error
	}

	entitlementService struct {
		entitlementRepository e.EntitlementRepository
		featureRepository     f.FeatureRepository
		subFeatureRepository  sf.SubFeatureRepository
		actionCategoryRepo    ac.ActionCategoryRepository
	}
)

func NewEntitlementService(
	entitlementRepo e.EntitlementRepository,
	featureRepo f.FeatureRepository,
	subFeatureRepo sf.SubFeatureRepository,
	actionCategoryRepo ac.ActionCategoryRepository,
) EntitlementService {
	return &entitlementService{
		entitlementRepository: entitlementRepo,
		featureRepository:     featureRepo,
		subFeatureRepository:  subFeatureRepo,
		actionCategoryRepo:    actionCategoryRepo,
	}
}

func (s *entitlementService) Create(ctx context.Context, req dto_request.CreateEntitlementRequest) (dto_response.GetEntitlementResponse, error) {
	// Parse dan validasi FeatureID
	parsedFeatureID, err := uuid.Parse(req.FeatureID)
	if err != nil {
		return dto_response.GetEntitlementResponse{}, myerror.InvalidRequest(err)
	}

	// Load feature untuk ambil prefix — dibutuhkan untuk build key
	feature, err := s.featureRepository.GetById(ctx, nil, parsedFeatureID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.GetEntitlementResponse{}, myerror.RecordNotFound("feature")
		}
		return dto_response.GetEntitlementResponse{}, myerror.DatabaseError(err)
	}

	// Parse ActionCategoryID
	parsedActionCategoryID, err := uuid.Parse(req.ActionCategoryID)
	if err != nil {
		return dto_response.GetEntitlementResponse{}, myerror.InvalidRequest(err)
	}

	// Load action category untuk ambil slug — dibutuhkan untuk build key
	actionCategory, err := s.actionCategoryRepo.GetById(ctx, nil, parsedActionCategoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.GetEntitlementResponse{}, myerror.RecordNotFound("action category")
		}
		return dto_response.GetEntitlementResponse{}, myerror.DatabaseError(err)
	}

	// Validasi dan proses SubFeatureID berdasarkan level
	var parsedSubFeatureID *uuid.UUID
	var subFeatureSlug string
	var subFeature *entity.SubFeature

	if entity.EntitlementLevel(req.Level) == entity.EntitlementLevelSubFeature {
		// Level sub_fitur — SubFeatureID wajib ada
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

		// Load sub feature untuk ambil slug dan validasi ownership ke feature induk
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
		// Level fitur — SubFeatureID harus kosong
		if req.SubFeatureID != nil && *req.SubFeatureID != "" {
			return dto_response.GetEntitlementResponse{}, myerror.New(
				"sub_feature_id must be empty when level is 'fitur'",
				myerror.Error_InvalidRequest,
			)
		}
	}

	// Build key otomatis dari komponen
	key := entity.BuildEntitlementKey(feature.Prefix, subFeatureSlug, actionCategory.Slug)

	// Cek duplikasi key — key harus unik global
	_, err = s.entitlementRepository.GetByKey(ctx, nil, key)
	if err == nil {
		return dto_response.GetEntitlementResponse{}, myerror.RecordAlreadyExist(
			fmt.Sprintf("entitlement with key '%s'", key),
		)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto_response.GetEntitlementResponse{}, myerror.DatabaseError(err)
	}

	// Cek duplikasi kombinasi komponen — sebagai double-check sebelum insert
	_, err = s.entitlementRepository.GetByComponents(ctx, nil, parsedFeatureID, parsedSubFeatureID, parsedActionCategoryID)
	if err == nil {
		return dto_response.GetEntitlementResponse{}, myerror.RecordAlreadyExist(
			"entitlement with same feature, sub feature, and action combination",
		)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto_response.GetEntitlementResponse{}, myerror.DatabaseError(err)
	}

	var resetPeriod *string
	if req.ResetPeriod != "" {
		resetPeriod = &req.ResetPeriod
	}

	newEntitlement := entity.NewEntitlement(
		key,
		req.Name,
		req.Description,
		req.RestrictionType,
		*resetPeriod,
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

	// Inject relasi yang sudah di-load ke result untuk mapping response tanpa query tambahan
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

	var res []dto_response.GetEntitlementResponse
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

func (s *entitlementService) Delete(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return myerror.InvalidRequest(err)
	}

	// Existence check — fail-fast
	_, err = s.entitlementRepository.GetById(ctx, nil, parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return myerror.RecordNotFound("entitlement")
		}
		return myerror.DatabaseError(err)
	}

	// Guard: tidak boleh dihapus jika masih dipakai oleh DurationAccessMapping
	count, err := s.entitlementRepository.CountMappings(ctx, nil, parsedID)
	if err != nil {
		return myerror.DatabaseError(err)
	}
	if count > 0 {
		return myerror.New(
			fmt.Sprintf("entitlement cannot be deleted because it is used in %d access mapping(s)", count),
			myerror.Error_InvalidRequest,
		)
	}

	if err := s.entitlementRepository.Delete(ctx, nil, parsedID); err != nil {
		return myerror.DatabaseError(err)
	}
	return nil
}

// toEntitlementResponse adalah helper mapping entity → response DTO.
// Asumsikan relasi Feature, SubFeature, ActionCategory sudah di-preload oleh repository.
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
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
	}

	// SubFeature hanya diisi jika level sub_fitur
	if e.SubFeature != nil {
		sfID := e.SubFeature.ID.String()
		sfName := e.SubFeature.Name
		res.SubFeatureId = &sfID
		res.SubFeatureName = &sfName
	}

	return res
}
