package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"gorm.io/gorm"
)

type (
	PersonaRepository interface {
		Create(ctx context.Context, tx *gorm.DB, session entity.Persona) (entity.Persona, error)
		GetByID(ctx context.Context, tx *gorm.DB, id string) (entity.Persona, error)
		GetByUserID(ctx context.Context, tx *gorm.DB, userID string) (entity.Persona, error)
		Update(ctx context.Context, tx *gorm.DB, session entity.Persona) (entity.Persona, error)
	}

	personaRepository struct {
		db *gorm.DB
	}
)

func NewPersona(db *gorm.DB) PersonaRepository {
	return &personaRepository{db}
}

func (r *personaRepository) Create(ctx context.Context, tx *gorm.DB, persona entity.Persona) (entity.Persona, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&persona).Error; err != nil {
		return persona, err
	}

	return persona, nil
}

func (r *personaRepository) GetByID(ctx context.Context, tx *gorm.DB, id string) (entity.Persona, error) {
	if tx == nil {
		tx = r.db
	}

	var persona entity.Persona
	if err := tx.WithContext(ctx).Take(&persona, "id = ?", id).Error; err != nil {
		return entity.Persona{}, err
	}

	return persona, nil
}

func (r *personaRepository) GetByUserID(ctx context.Context, tx *gorm.DB, userID string) (entity.Persona, error) {
	if tx == nil {
		tx = r.db
	}

	var persona entity.Persona
	if err := tx.WithContext(ctx).Take(&persona, "user_id = ?", userID).Error; err != nil {
		return entity.Persona{}, err
	}
	return persona, nil
}

func (r *personaRepository) Update(ctx context.Context, tx *gorm.DB, persona entity.Persona) (entity.Persona, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Save(&persona).Error; err != nil {
		return persona, err
	}

	return persona, nil
}
