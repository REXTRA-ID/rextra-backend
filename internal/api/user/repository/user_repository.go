package userRepository

import (
	"context"
	"rextra-backend/internal/entity"

	"gorm.io/gorm"
)

type (
	UserRepository interface {
		Create(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error)
		GetById(ctx context.Context, tx *gorm.DB, userId string) (entity.User, error)
		GetByEmail(ctx context.Context, tx *gorm.DB, email string) (entity.User, error)
		Update(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error)
	}

	userRepository struct {
		db *gorm.DB
	}
)

func New(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (r *userRepository) GetById(ctx context.Context, tx *gorm.DB, userId string) (entity.User, error) {
	if tx == nil {
		tx = r.db
	}

	var user entity.User
	if err := tx.WithContext(ctx).Take(&user, "id = ?", userId).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, tx *gorm.DB, email string) (entity.User, error) {
	if tx == nil {
		tx = r.db
	}

	var user entity.User
	if err := tx.WithContext(ctx).Take(&user, "email = ?", email).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Save(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}
