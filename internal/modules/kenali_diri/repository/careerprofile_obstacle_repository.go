package repository

import (
	"rextra-backend/internal/entity"

	"gorm.io/gorm"
)

type CareerProfileObstacleRepository interface {
	FindAllStudentObstacles() ([]entity.CareerProfileObstacleOption, error)
}

type careerProfileObstacleRepository struct {
	db *gorm.DB
}

func NewCareerProfileObstacleRepository(db *gorm.DB) CareerProfileObstacleRepository {
	return &careerProfileObstacleRepository{
		db: db,
	}
}

func (r *careerProfileObstacleRepository) FindAllStudentObstacles() ([]entity.CareerProfileObstacleOption, error) {
	var obstacles []entity.CareerProfileObstacleOption
	err := r.db.Order("sort_order ASC").Find(&obstacles).Error
	return obstacles, err
}
