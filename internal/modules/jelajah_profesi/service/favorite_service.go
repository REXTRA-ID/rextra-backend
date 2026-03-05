package service

import (
	"context"
	"errors"
	"time"

	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/modules/jelajah_profesi/repository"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const timeFormatFavorite = "02 Jan 2006 15:04"

type (
	FavoriteService interface {
		ListFavorites(ctx context.Context, userID uuid.UUID, page, limit int) (dto_response.FavoriteProfessionListResponse, error)
		AddFavorite(ctx context.Context, userID uuid.UUID, professionID int64) error
		RemoveFavorite(ctx context.Context, userID uuid.UUID, professionID int64) error
	}

	favoriteService struct {
		favoriteRepo   repository.FavoriteProfessionRepository
		professionRepo repository.ProfessionRepository
	}
)

func NewFavorite(
	favoriteRepo repository.FavoriteProfessionRepository,
	professionRepo repository.ProfessionRepository,
) FavoriteService {
	return &favoriteService{
		favoriteRepo:   favoriteRepo,
		professionRepo: professionRepo,
	}
}

func (s *favoriteService) ListFavorites(ctx context.Context, userID uuid.UUID, page, limit int) (dto_response.FavoriteProfessionListResponse, error) {
	pag := utils.NewPaginationParams(page, limit)

	favorites, total, err := s.favoriteRepo.ListByUserID(ctx, userID, pag.Limit, pag.GetOffset())
	if err != nil {
		return dto_response.FavoriteProfessionListResponse{}, err
	}

	items := make([]dto_response.FavoriteProfessionItem, 0, len(favorites))
	for _, f := range favorites {
		items = append(items, dto_response.FavoriteProfessionItem{
			ID:               f.ID,
			ProfessionID:     f.ProfessionID,
			ProfessionSlug:   f.Profession.Slug,
			ProfessionName:   f.Profession.Name,
			ProfessionImage:  f.Profession.ImageURL,
			MainCategoryName: f.Profession.MainCategory.Name,
			FavoritedAt:      f.CreatedAt.Format(time.RFC3339),
		})
	}

	meta := utils.CalculatePaginationMeta(pag.Page, pag.Limit, total)
	return dto_response.FavoriteProfessionListResponse{
		Data: items,
		Pagination: dto_response.PaginationMeta{
			CurrentPage:  meta.CurrentPage,
			TotalPages:   meta.TotalPages,
			TotalRecords: meta.TotalRecords,
			PerPage:      meta.PerPage,
		},
	}, nil
}

func (s *favoriteService) AddFavorite(ctx context.Context, userID uuid.UUID, professionID int64) error {
	// Check if already favorited
	already, err := s.favoriteRepo.IsFavorited(ctx, userID, professionID)
	if err != nil {
		return err
	}
	if already {
		return myerror.RecordAlreadyExist("profesi favorit")
	}

	return s.favoriteRepo.Add(ctx, userID, professionID)
}

func (s *favoriteService) RemoveFavorite(ctx context.Context, userID uuid.UUID, professionID int64) error {
	exists, err := s.favoriteRepo.IsFavorited(ctx, userID, professionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return myerror.RecordNotFound("profesi favorit")
		}
		return err
	}
	if !exists {
		return myerror.RecordNotFound("profesi favorit")
	}

	return s.favoriteRepo.Remove(ctx, userID, professionID)
}
