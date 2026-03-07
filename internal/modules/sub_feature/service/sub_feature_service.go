package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	frepo "rextra-backend/internal/modules/feature/repository"
	sfrepo "rextra-backend/internal/modules/sub_feature/repository"
	myerror "rextra-backend/internal/pkg/error"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	SubFeatureService interface {
		Create(ctx context.Context, featureID string, req dto_request.CreateSubFeatureRequest) (dto_response.GetSubFeatureResponse, error)
		GetAllByFeatureID(ctx context.Context, featureID string) ([]dto_response.GetSubFeatureResponse, error)
		GetById(ctx context.Context, featureID, id string) (dto_response.GetSubFeatureResponse, error)
		Update(ctx context.Context, featureID, id string, req dto_request.UpdateSubFeatureRequest) (dto_response.GetSubFeatureResponse, error)
		Delete(ctx context.Context, featureID, id string) error
	}

	subFeatureService struct {
		featureRepository    frepo.FeatureRepository
		subFeatureRepository sfrepo.SubFeatureRepository
	}
)

func NewSubFeatureService(featureRepo frepo.FeatureRepository, subFeatureRepo sfrepo.SubFeatureRepository) SubFeatureService {
	return &subFeatureService{
		featureRepository:    featureRepo,
		subFeatureRepository: subFeatureRepo,
	}
}

func (s *subFeatureService) Create(ctx context.Context, featureID string, req dto_request.CreateSubFeatureRequest) (dto_response.GetSubFeatureResponse, error) {
	parsedFeatureID, err := uuid.Parse(featureID)
	if err != nil {
		return dto_response.GetSubFeatureResponse{}, myerror.InvalidRequest(err)
	}

	// Validasi feature induk exist dan bertipe bertingkat
	feature, err := s.featureRepository.GetById(ctx, nil, parsedFeatureID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.GetSubFeatureResponse{}, myerror.RecordNotFound("feature")
		}
		return dto_response.GetSubFeatureResponse{}, myerror.DatabaseError(err)
	}
	if feature.Type != entity.FeatureTypeHierarchic {
		return dto_response.GetSubFeatureResponse{}, myerror.New(
			"sub feature can only be added to feature with type 'bertingkat'",
			myerror.Error_InvalidRequest,
		)
	}

	// Cek duplikasi slug dalam scope feature yang sama
	_, err = s.subFeatureRepository.GetBySlugAndFeatureID(ctx, nil, req.Slug, parsedFeatureID)
	if err == nil {
		return dto_response.GetSubFeatureResponse{}, myerror.RecordAlreadyExist(
			"sub feature with slug '" + req.Slug + "' in this feature",
		)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto_response.GetSubFeatureResponse{}, myerror.DatabaseError(err)
	}

	newSubFeature := entity.NewSubFeature(parsedFeatureID, req.Name, req.Slug, req.Description)
	result, err := s.subFeatureRepository.Create(ctx, nil, newSubFeature)
	if err != nil {
		return dto_response.GetSubFeatureResponse{}, myerror.DatabaseError(err)
	}

	return toSubFeatureResponse(result, 0, feature.Name), nil
}

func (s *subFeatureService) GetAllByFeatureID(ctx context.Context, featureID string) ([]dto_response.GetSubFeatureResponse, error) {
	parsedFeatureID, err := uuid.Parse(featureID)
	if err != nil {
		return nil, myerror.InvalidRequest(err)
	}

	// Validasi feature exist
	feature, err := s.featureRepository.GetById(ctx, nil, parsedFeatureID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, myerror.RecordNotFound("feature")
		}
		return nil, myerror.DatabaseError(err)
	}

	results, countMap, err := s.subFeatureRepository.GetAllByFeatureIDWithCounts(ctx, nil, parsedFeatureID)
	if err != nil {
		return nil, myerror.DatabaseError(err)
	}

	var res []dto_response.GetSubFeatureResponse
	for _, r := range results {
		res = append(res, toSubFeatureResponse(r, int(countMap[r.ID]), feature.Name))
	}
	return res, nil
}

func (s *subFeatureService) GetById(ctx context.Context, featureID, id string) (dto_response.GetSubFeatureResponse, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return dto_response.GetSubFeatureResponse{}, myerror.InvalidRequest(err)
	}

	result, err := s.subFeatureRepository.GetById(ctx, nil, parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.GetSubFeatureResponse{}, myerror.RecordNotFound("sub feature")
		}
		return dto_response.GetSubFeatureResponse{}, myerror.DatabaseError(err)
	}

	// Validasi sub feature memang milik feature yang diminta (security check)
	if result.FeatureID.String() != featureID {
		return dto_response.GetSubFeatureResponse{}, myerror.RecordNotFound("sub feature")
	}

	count, _ := s.subFeatureRepository.CountEntitlements(ctx, nil, result.ID)
	return toSubFeatureResponse(result, int(count), result.Feature.Name), nil
}

func (s *subFeatureService) Update(ctx context.Context, featureID, id string, req dto_request.UpdateSubFeatureRequest) (dto_response.GetSubFeatureResponse, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return dto_response.GetSubFeatureResponse{}, myerror.InvalidRequest(err)
	}

	existing, err := s.subFeatureRepository.GetById(ctx, nil, parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.GetSubFeatureResponse{}, myerror.RecordNotFound("sub feature")
		}
		return dto_response.GetSubFeatureResponse{}, myerror.DatabaseError(err)
	}

	// Validasi ownership
	if existing.FeatureID.String() != featureID {
		return dto_response.GetSubFeatureResponse{}, myerror.RecordNotFound("sub feature")
	}

	// Cek immutability slug
	if req.Slug != existing.Slug {
		count, err := s.subFeatureRepository.CountEntitlements(ctx, nil, parsedID)
		if err != nil {
			return dto_response.GetSubFeatureResponse{}, myerror.DatabaseError(err)
		}
		if count > 0 {
			return dto_response.GetSubFeatureResponse{}, myerror.New(
				fmt.Sprintf("slug cannot be changed because this sub feature has %d entitlement(s)", count),
				myerror.Error_InvalidRequest,
			)
		}

		// Cek duplikasi slug baru dalam scope feature yang sama
		parsedFeatureID := existing.FeatureID
		_, err = s.subFeatureRepository.GetBySlugAndFeatureID(ctx, nil, req.Slug, parsedFeatureID)
		if err == nil {
			return dto_response.GetSubFeatureResponse{}, myerror.RecordAlreadyExist(
				"sub feature with slug '" + req.Slug + "' in this feature",
			)
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.GetSubFeatureResponse{}, myerror.DatabaseError(err)
		}
	}

	existing.Name = req.Name
	existing.Slug = req.Slug
	existing.Description = req.Description
	existing.Status = entity.SubFeatureStatus(req.Status)
	existing.UpdatedAt = time.Now().UTC()

	result, err := s.subFeatureRepository.Update(ctx, nil, existing)
	if err != nil {
		return dto_response.GetSubFeatureResponse{}, myerror.DatabaseError(err)
	}

	count, _ := s.subFeatureRepository.CountEntitlements(ctx, nil, result.ID)
	// Intentional: pakai existing.Feature.Name, bukan result.Feature.Name.
	// repository.Update() memanggil gorm Save() yang tidak preload relasi —
	// result.Feature akan kosong setelah Save. existing di-load via GetById
	// yang sudah Preload("Feature"), sehingga existing.Feature.Name aman dipakai.
	return toSubFeatureResponse(result, int(count), existing.Feature.Name), nil
}

func (s *subFeatureService) Delete(ctx context.Context, featureID, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return myerror.InvalidRequest(err)
	}

	// Existence check + ownership check
	existing, err := s.subFeatureRepository.GetById(ctx, nil, parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return myerror.RecordNotFound("sub feature")
		}
		return myerror.DatabaseError(err)
	}
	if existing.FeatureID.String() != featureID {
		return myerror.RecordNotFound("sub feature")
	}

	// Tidak boleh dihapus jika masih ada entitlement
	count, err := s.subFeatureRepository.CountEntitlements(ctx, nil, parsedID)
	if err != nil {
		return myerror.DatabaseError(err)
	}
	if count > 0 {
		return myerror.New(
			fmt.Sprintf("sub feature cannot be deleted because it has %d entitlement(s)", count),
			myerror.Error_InvalidRequest,
		)
	}

	if err := s.subFeatureRepository.Delete(ctx, nil, parsedID); err != nil {
		return myerror.DatabaseError(err)
	}
	return nil
}

func toSubFeatureResponse(subFeature entity.SubFeature, entitlementCount int, featureName string) dto_response.GetSubFeatureResponse {
	return dto_response.GetSubFeatureResponse{
		Id:               subFeature.ID.String(),
		FeatureId:        subFeature.FeatureID.String(),
		FeatureName:      featureName,
		Name:             subFeature.Name,
		Slug:             subFeature.Slug,
		Description:      subFeature.Description,
		Status:           string(subFeature.Status),
		EntitlementCount: entitlementCount,
		CreatedAt:        subFeature.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        subFeature.UpdatedAt.Format(time.RFC3339),
	}
}
