package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	"rextra-backend/internal/modules/action_category/repository"
	myerror "rextra-backend/internal/pkg/error"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	ActionCategoryService interface {
		Create(ctx context.Context, req dto_request.CreateActionCategoryRequest) (dto_response.GetActionCategoryResponse, error)
		GetAll(ctx context.Context) ([]dto_response.GetActionCategoryResponse, error)
		GetById(ctx context.Context, id string) (dto_response.GetActionCategoryResponse, error)
		Update(ctx context.Context, id string, req dto_request.UpdateActionCategoryRequest) (dto_response.GetActionCategoryResponse, error)
		Delete(ctx context.Context, id string) error
	}

	actionCategoryService struct {
		actionCategoryRepository repository.ActionCategoryRepository
	}
)

func NewActionCategoryService(repo repository.ActionCategoryRepository) ActionCategoryService {
	return &actionCategoryService{
		actionCategoryRepository: repo,
	}
}

func (s *actionCategoryService) Create(ctx context.Context, req dto_request.CreateActionCategoryRequest) (dto_response.GetActionCategoryResponse, error) {
	_, err := s.actionCategoryRepository.GetBySlug(ctx, nil, req.Slug)
	if err == nil {
		return dto_response.GetActionCategoryResponse{}, myerror.RecordAlreadyExist("action category with slug '" + req.Slug + "'")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto_response.GetActionCategoryResponse{}, myerror.DatabaseError(err)
	}

	newCategory := entity.NewActionCategory(req.Name, req.Slug, req.Description, "active")
	result, err := s.actionCategoryRepository.Create(ctx, nil, newCategory)
	if err != nil {
		return dto_response.GetActionCategoryResponse{}, myerror.DatabaseError(err)
	}

	return toActionCategoryResponse(result, 0), nil
}

func (s *actionCategoryService) GetAll(ctx context.Context) ([]dto_response.GetActionCategoryResponse, error) {
	results, countMap, err := s.actionCategoryRepository.GetAllWithEntitlementCount(ctx, nil)
	if err != nil {
		return nil, myerror.DatabaseError(err)
	}

	var res []dto_response.GetActionCategoryResponse
	for _, r := range results {
		res = append(res, toActionCategoryResponse(r, int(countMap[r.ID])))
	}
	return res, nil
}

func (s *actionCategoryService) GetById(ctx context.Context, id string) (dto_response.GetActionCategoryResponse, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return dto_response.GetActionCategoryResponse{}, myerror.InvalidRequest(err)
	}

	result, err := s.actionCategoryRepository.GetById(ctx, nil, parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.GetActionCategoryResponse{}, myerror.RecordNotFound("action category")
		}
		return dto_response.GetActionCategoryResponse{}, myerror.DatabaseError(err)
	}

	count, _ := s.actionCategoryRepository.CountEntitlements(ctx, nil, result.ID)
	return toActionCategoryResponse(result, int(count)), nil
}

func (s *actionCategoryService) Update(ctx context.Context, id string, req dto_request.UpdateActionCategoryRequest) (dto_response.GetActionCategoryResponse, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return dto_response.GetActionCategoryResponse{}, myerror.InvalidRequest(err)
	}

	existing, err := s.actionCategoryRepository.GetById(ctx, nil, parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.GetActionCategoryResponse{}, myerror.RecordNotFound("action category")
		}
		return dto_response.GetActionCategoryResponse{}, myerror.DatabaseError(err)
	}

	if req.Slug != existing.Slug {
		count, err := s.actionCategoryRepository.CountEntitlements(ctx, nil, parsedID)
		if err != nil {
			return dto_response.GetActionCategoryResponse{}, myerror.DatabaseError(err)
		}
		if count > 0 {
			return dto_response.GetActionCategoryResponse{}, myerror.New(
				"slug cannot be changed because it is already used by "+
					formatCount(count)+" entitlement(s)",
				myerror.Error_InvalidRequest,
			)
		}

		_, err = s.actionCategoryRepository.GetBySlug(ctx, nil, req.Slug)
		if err == nil {
			return dto_response.GetActionCategoryResponse{}, myerror.RecordAlreadyExist("action category with slug '" + req.Slug + "'")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.GetActionCategoryResponse{}, myerror.DatabaseError(err)
		}
	}

	existing.Name = req.Name
	existing.Slug = req.Slug
	existing.Description = req.Description
	existing.Status = entity.ActionCategoryStatus(req.Status)
	existing.UpdatedAt = time.Now().UTC()

	result, err := s.actionCategoryRepository.Update(ctx, nil, existing)
	if err != nil {
		return dto_response.GetActionCategoryResponse{}, myerror.DatabaseError(err)
	}

	count, _ := s.actionCategoryRepository.CountEntitlements(ctx, nil, result.ID)
	return toActionCategoryResponse(result, int(count)), nil
}

func (s *actionCategoryService) Delete(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return myerror.InvalidRequest(err)
	}

	_, err = s.actionCategoryRepository.GetById(ctx, nil, parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return myerror.RecordNotFound("action category")
		}
		return myerror.DatabaseError(err)
	}

	count, err := s.actionCategoryRepository.CountEntitlements(ctx, nil, parsedID)
	if err != nil {
		return myerror.DatabaseError(err)
	}
	if count > 0 {
		return myerror.New(
			"action category cannot be deleted because it is used by "+
				formatCount(count)+" entitlement(s)",
			myerror.Error_InvalidRequest,
		)
	}

	if err := s.actionCategoryRepository.Delete(ctx, nil, parsedID); err != nil {
		return myerror.DatabaseError(err)
	}
	return nil
}

func toActionCategoryResponse(a entity.ActionCategory, entitlementCount int) dto_response.GetActionCategoryResponse {
	createdAt := ""
	updatedAt := ""
	if !a.CreatedAt.IsZero() {
		createdAt = a.CreatedAt.Format(time.RFC3339)
	}
	if !a.UpdatedAt.IsZero() {
		updatedAt = a.UpdatedAt.Format(time.RFC3339)
	}
	return dto_response.GetActionCategoryResponse{
		Id:               a.ID.String(),
		Name:             a.Name,
		Slug:             a.Slug,
		Description:      a.Description,
		Status:           string(a.Status),
		EntitlementCount: entitlementCount,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}
}

func formatCount(count int64) string {
	return fmt.Sprintf("%d", count)
}
