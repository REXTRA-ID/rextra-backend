package repository

import (
	"context"

	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	FavoriteProfessionRepository interface {
		ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.UserFavoriteProfession, int64, error)
		Add(ctx context.Context, userID uuid.UUID, professionID int64) error
		Remove(ctx context.Context, userID uuid.UUID, professionID int64) error
		IsFavorited(ctx context.Context, userID uuid.UUID, professionID int64) (bool, error)
		GetFavoritedMap(ctx context.Context, userID uuid.UUID, professionIDs []int64) (map[int64]bool, error)
	}

	favoriteProfessionRepository struct {
		db *gorm.DB
	}
)

func NewFavoriteProfession(db *gorm.DB) FavoriteProfessionRepository {
	return &favoriteProfessionRepository{db: db}
}

func (r *favoriteProfessionRepository) ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.UserFavoriteProfession, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).
		Model(&entity.UserFavoriteProfession{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var favorites []entity.UserFavoriteProfession
	if err := r.db.WithContext(ctx).
		Preload("Profession").
		Preload("Profession.MainCategory").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&favorites).Error; err != nil {
		return nil, 0, err
	}

	return favorites, total, nil
}

func (r *favoriteProfessionRepository) Add(ctx context.Context, userID uuid.UUID, professionID int64) error {
	fav := entity.UserFavoriteProfession{
		UserID:       userID,
		ProfessionID: professionID,
	}
	return r.db.WithContext(ctx).Create(&fav).Error
}

func (r *favoriteProfessionRepository) Remove(ctx context.Context, userID uuid.UUID, professionID int64) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND profession_id = ?", userID, professionID).
		Delete(&entity.UserFavoriteProfession{}).Error
}

func (r *favoriteProfessionRepository) IsFavorited(ctx context.Context, userID uuid.UUID, professionID int64) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.UserFavoriteProfession{}).
		Where("user_id = ? AND profession_id = ?", userID, professionID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *favoriteProfessionRepository) GetFavoritedMap(ctx context.Context, userID uuid.UUID, professionIDs []int64) (map[int64]bool, error) {
	result := make(map[int64]bool)
	if len(professionIDs) == 0 {
		return result, nil
	}

	var favorites []entity.UserFavoriteProfession
	if err := r.db.WithContext(ctx).
		Select("profession_id").
		Where("user_id = ? AND profession_id IN ?", userID, professionIDs).
		Find(&favorites).Error; err != nil {
		return nil, err
	}

	for _, f := range favorites {
		result[f.ProfessionID] = true
	}
	return result, nil
}
