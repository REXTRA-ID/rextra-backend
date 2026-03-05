package service

import (
	"context"

	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/modules/jelajah_profesi/repository"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/utils"

	"github.com/google/uuid"
)

type (
	ProfessionService interface {
		ListProfessions(ctx context.Context, req dto_request.ListProfessionsRequest, userID uuid.UUID) (dto_response.ProfessionListResponse, error)
		GetProfessionDetail(ctx context.Context, slug string, userID uuid.UUID) (dto_response.ProfessionDetailResponse, error)
		ListMainCategories(ctx context.Context) (dto_response.MainCategoryListResponse, error)
		ListSubCategories(ctx context.Context, mainCategoryID int64) (dto_response.SubCategoryListResponse, error)
	}

	professionService struct {
		professionRepo repository.ProfessionRepository
		categoryRepo   repository.ProfessionCategoryRepository
		detailRepo     repository.ProfessionDetailRepository
		favoriteRepo   repository.FavoriteProfessionRepository
	}
)

func NewProfession(
	professionRepo repository.ProfessionRepository,
	categoryRepo repository.ProfessionCategoryRepository,
	detailRepo repository.ProfessionDetailRepository,
	favoriteRepo repository.FavoriteProfessionRepository,
) ProfessionService {
	return &professionService{
		professionRepo: professionRepo,
		categoryRepo:   categoryRepo,
		detailRepo:     detailRepo,
		favoriteRepo:   favoriteRepo,
	}
}

func (s *professionService) ListProfessions(ctx context.Context, req dto_request.ListProfessionsRequest, userID uuid.UUID) (dto_response.ProfessionListResponse, error) {
	pag := utils.NewPaginationParams(req.Page, req.Limit)

	filters := repository.ProfessionFilters{
		Search:         req.Search,
		MainCategoryID: req.MainCategoryID,
		SubCategoryID:  req.SubCategoryID,
	}

	professions, total, err := s.professionRepo.List(ctx, filters, pag.Limit, pag.GetOffset())
	if err != nil {
		return dto_response.ProfessionListResponse{}, err
	}

	// Batch-check favorit status
	profIDs := make([]int64, len(professions))
	for i, p := range professions {
		profIDs[i] = p.ID
	}

	favMap := make(map[int64]bool)
	if userID != uuid.Nil {
		favMap, _ = s.favoriteRepo.GetFavoritedMap(ctx, userID, profIDs)
	}

	items := make([]dto_response.ProfessionListItem, 0, len(professions))
	for _, p := range professions {
		items = append(items, dto_response.ProfessionListItem{
			ID:              p.ID,
			Slug:            p.Slug,
			Name:            p.Name,
			ImageURL:        p.ImageURL,
			SubCategoryName: p.SubCategory.Name,
			IsFavorited:     favMap[p.ID],
		})
	}

	meta := utils.CalculatePaginationMeta(pag.Page, pag.Limit, total)
	return dto_response.ProfessionListResponse{
		Data: items,
		Pagination: dto_response.PaginationMeta{
			CurrentPage:  meta.CurrentPage,
			TotalPages:   meta.TotalPages,
			TotalRecords: meta.TotalRecords,
			PerPage:      meta.PerPage,
		},
	}, nil
}

func (s *professionService) GetProfessionDetail(ctx context.Context, slug string, userID uuid.UUID) (dto_response.ProfessionDetailResponse, error) {
	profession, err := s.professionRepo.GetBySlug(ctx, slug)
	if err != nil {
		return dto_response.ProfessionDetailResponse{}, myerror.RecordNotFound("profession")
	}

	// Fetch all relations concurrently-safe (sequential for simplicity)
	activities, _ := s.detailRepo.GetActivities(ctx, profession.ID)
	skills, _ := s.detailRepo.GetSkills(ctx, profession.ID)
	tools, _ := s.detailRepo.GetTools(ctx, profession.ID)
	careerPaths, _ := s.detailRepo.GetCareerPaths(ctx, profession.ID)
	marketInsights, _ := s.detailRepo.GetMarketInsights(ctx, profession.ID)
	studyPrograms, _ := s.detailRepo.GetStudyPrograms(ctx, profession.ID)
	aliases, _ := s.detailRepo.GetAliases(ctx, profession.ID)

	// Check favorit
	isFav := false
	if userID != uuid.Nil {
		isFav, _ = s.favoriteRepo.IsFavorited(ctx, userID, profession.ID)
	}

	resp := dto_response.ProfessionDetailResponse{
		ID:                profession.ID,
		Slug:              profession.Slug,
		Name:              profession.Name,
		ImageURL:          profession.ImageURL,
		AboutDescription:  profession.AboutDescription,
		RiasecDescription: profession.RiasecDescription,
		IsFavorited:       isFav,
		MainCategory: dto_response.ProfessionCategoryItem{
			ID:   profession.MainCategory.ID,
			Code: profession.MainCategory.Code,
			Name: profession.MainCategory.Name,
		},
		SubCategory: dto_response.ProfessionCategoryItem{
			ID:   profession.SubCategory.ID,
			Code: profession.SubCategory.Code,
			Name: profession.SubCategory.Name,
		},
	}

	// Riasec code
	if profession.RiasecCode != nil {
		resp.RiasecCode = &dto_response.ProfessionRiasecCodeItem{
			ID:          profession.RiasecCode.ID,
			RiasecCode:  profession.RiasecCode.RiasecCode,
			RiasecTitle: profession.RiasecCode.RiasecTitle,
		}
	}

	// Aliases
	resp.Aliases = make([]dto_response.ProfessionAliasItem, 0, len(aliases))
	for _, a := range aliases {
		resp.Aliases = append(resp.Aliases, dto_response.ProfessionAliasItem{
			ID:        a.ID,
			AliasName: a.AliasName,
		})
	}

	// Activities
	resp.Activities = make([]dto_response.ProfessionActivityItem, 0, len(activities))
	for _, a := range activities {
		resp.Activities = append(resp.Activities, dto_response.ProfessionActivityItem{
			ID:          a.ID,
			Description: a.Description,
			SortOrder:   a.SortOrder,
		})
	}

	// Skills
	resp.Skills = make([]dto_response.ProfessionSkillItem, 0, len(skills))
	for _, s := range skills {
		resp.Skills = append(resp.Skills, dto_response.ProfessionSkillItem{
			SkillID:   s.SkillID,
			SkillName: s.Skill.Name,
			SkillType: s.SkillType,
			Priority:  s.Priority,
		})
	}

	// Tools
	resp.Tools = make([]dto_response.ProfessionToolItem, 0, len(tools))
	for _, t := range tools {
		resp.Tools = append(resp.Tools, dto_response.ProfessionToolItem{
			ToolID:    t.ToolID,
			ToolName:  t.Tool.Name,
			UsageType: t.UsageType,
		})
	}

	// Career paths
	resp.CareerPaths = make([]dto_response.ProfessionCareerPathItem, 0, len(careerPaths))
	for _, c := range careerPaths {
		resp.CareerPaths = append(resp.CareerPaths, dto_response.ProfessionCareerPathItem{
			ID:              c.ID,
			Title:           c.Title,
			ExperienceRange: c.ExperienceRange,
			SalaryMin:       c.SalaryMin,
			SalaryMax:       c.SalaryMax,
			SortOrder:       c.SortOrder,
		})
	}

	// Market insights
	resp.MarketInsights = make([]dto_response.ProfessionMarketInsightItem, 0, len(marketInsights))
	for _, m := range marketInsights {
		resp.MarketInsights = append(resp.MarketInsights, dto_response.ProfessionMarketInsightItem{
			ID:          m.ID,
			Description: m.Description,
			SortOrder:   m.SortOrder,
		})
	}

	// Study programs
	resp.StudyPrograms = make([]dto_response.ProfessionStudyProgramItem, 0, len(studyPrograms))
	for _, sp := range studyPrograms {
		resp.StudyPrograms = append(resp.StudyPrograms, dto_response.ProfessionStudyProgramItem{
			StudyProgramID:   sp.StudyProgramID,
			StudyProgramName: sp.StudyProgram.Name,
		})
	}

	return resp, nil
}

func (s *professionService) ListMainCategories(ctx context.Context) (dto_response.MainCategoryListResponse, error) {
	categories, err := s.categoryRepo.ListMainCategories(ctx)
	if err != nil {
		return dto_response.MainCategoryListResponse{}, err
	}

	items := make([]dto_response.MainCategoryItem, 0, len(categories))
	for _, c := range categories {
		items = append(items, dto_response.MainCategoryItem{
			ID:          c.ID,
			Code:        c.Code,
			Name:        c.Name,
			Description: c.Description,
		})
	}

	return dto_response.MainCategoryListResponse{Data: items}, nil
}

func (s *professionService) ListSubCategories(ctx context.Context, mainCategoryID int64) (dto_response.SubCategoryListResponse, error) {
	categories, err := s.categoryRepo.ListSubCategories(ctx, mainCategoryID)
	if err != nil {
		return dto_response.SubCategoryListResponse{}, err
	}

	items := make([]dto_response.SubCategoryItem, 0, len(categories))
	for _, c := range categories {
		items = append(items, dto_response.SubCategoryItem{
			ID:          c.ID,
			Code:        c.Code,
			Name:        c.Name,
			Description: c.Description,
		})
	}

	return dto_response.SubCategoryListResponse{Data: items}, nil
}
