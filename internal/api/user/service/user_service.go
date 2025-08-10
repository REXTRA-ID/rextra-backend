package userService

import (
	"context"

	userDto_response "rextra-backend/internal/api/user/dto/response"
	userRepository "rextra-backend/internal/api/user/repository"

	"gorm.io/gorm"
)

type (
	UserService interface {
		GetById(ctx context.Context, userId string) (userDto_response.UserResponse, error)
	}

	userService struct {
		userRepository userRepository.UserRepository
		db             *gorm.DB
	}
)

func New(userRepository userRepository.UserRepository,
	db *gorm.DB) UserService {
	return &userService{
		userRepository: userRepository,
		db:             db,
	}
}

func (s *userService) GetById(ctx context.Context, userId string) (userDto_response.UserResponse, error) {
	// user, err := s.userRepository.GetByIdWithFilmList(ctx, nil, userId)
	// if err != nil {
	// 	return userDto_response.UserResponse{}, err
	// }

	// return userDto_response.UserResponse{
	// 	ID:          user.ID.String(),
	// 	Username:    user.Username,
	// 	PhoneNumber: us,
	// }, nil
	return userDto_response.UserResponse{}, nil
}
