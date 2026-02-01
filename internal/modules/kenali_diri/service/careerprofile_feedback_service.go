package service

import (
	"context"

	dto_req "rextra-backend/internal/dto/request"
	dto_res "rextra-backend/internal/dto/response"
	"rextra-backend/internal/modules/kenali_diri/repository"
)

type CareerProfileFeedbackService interface {
	GetStudentFeedbacks(ctx context.Context, req dto_req.GetStudentFeedbackRequest) (*dto_res.CareerProfileStudentFeedbackListResponse, error)
	GetExpertFeedbacks(ctx context.Context, req dto_req.GetExpertFeedbackRequest) (*dto_res.CareerProfileExpertFeedbackListResponse, error)
	GetExpertFeedbackDetail(ctx context.Context, id int64) (*dto_res.CareerProfileExpertFeedbackDetailResponse, error)
	GetFeedbackMetadata(ctx context.Context) (*dto_res.CareerProfileFeedbackMetadataResponse, error)
}

type careerProfileFeedbackService struct {
	repo repository.CareerProfileFeedbackRepository
}

func NewCareerProfileFeedbackService(repo repository.CareerProfileFeedbackRepository) CareerProfileFeedbackService {
	return &careerProfileFeedbackService{
		repo: repo,
	}
}

func (s *careerProfileFeedbackService) GetStudentFeedbacks(ctx context.Context, req dto_req.GetStudentFeedbackRequest) (*dto_res.CareerProfileStudentFeedbackListResponse, error) {
	items, total, err := s.repo.GetStudentFeedbackList(ctx, req)
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if req.PageSize > 0 {
		totalPages = int(total) / req.PageSize
		if int(total)%req.PageSize > 0 {
			totalPages++
		}
	}

	return &dto_res.CareerProfileStudentFeedbackListResponse{
		Items: items,
		Pagination: dto_res.CareerProfilePaginationMeta{
			Page:       req.Page,
			PageSize:   req.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
			HasNext:    req.Page < totalPages,
			HasPrev:    req.Page > 1,
		},
	}, nil
}

func (s *careerProfileFeedbackService) GetExpertFeedbacks(ctx context.Context, req dto_req.GetExpertFeedbackRequest) (*dto_res.CareerProfileExpertFeedbackListResponse, error) {
	items, total, err := s.repo.GetExpertFeedbackList(ctx, req)
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if req.PageSize > 0 {
		totalPages = int(total) / req.PageSize
		if int(total)%req.PageSize > 0 {
			totalPages++
		}
	}

	return &dto_res.CareerProfileExpertFeedbackListResponse{
		Items: items,
		Pagination: dto_res.CareerProfilePaginationMeta{
			Page:       req.Page,
			PageSize:   req.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
			HasNext:    req.Page < totalPages,
			HasPrev:    req.Page > 1,
		},
	}, nil
}

func (s *careerProfileFeedbackService) GetExpertFeedbackDetail(ctx context.Context, id int64) (*dto_res.CareerProfileExpertFeedbackDetailResponse, error) {
	return s.repo.GetExpertFeedbackDetail(ctx, id)
}
func (s *careerProfileFeedbackService) GetFeedbackMetadata(ctx context.Context) (*dto_res.CareerProfileFeedbackMetadataResponse, error) {
	return s.repo.GetFeedbackMetadata(ctx)
}
