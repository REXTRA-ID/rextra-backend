package service

import (
	"context"

	"rextra-backend/internal/modules/user/repository"
	dto_response "rextra-backend/internal/dto/response"

	"gorm.io/gorm"
)

type (
	UserService interface {
		GetById(ctx context.Context, userId string) (dto_response.UserResponse, error)
	}

	userService struct {
		userRepository repository.UserRepository
		db             *gorm.DB
	}
)

func NewUser(userRepository repository.UserRepository,
	db *gorm.DB) UserService {
	return &userService{
		userRepository: userRepository,
		db:             db,
	}
}

func (s *userService) GetById(ctx context.Context, userId string) (dto_response.UserResponse, error) {
	// user, err := s.userRepository.GetByIdWithFilmList(ctx, nil, userId)
	// if err != nil {
	// 	return dto_response.UserResponse{}, err
	// }

	// return dto_response.UserResponse{
	// 	ID:          user.ID.String(),
	// 	Username:    user.Username,
	// 	PhoneNumber: us,
	// }, nil
	return dto_response.UserResponse{}, nil
}
