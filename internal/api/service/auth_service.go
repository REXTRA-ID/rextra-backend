package service

import (
	"context"
	"errors"
	"fmt"

	"rextra-backend/internal/api/repository"
	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	mailer "rextra-backend/internal/pkg/email"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/google/oauth"
	myjwt "rextra-backend/internal/pkg/jwt"
	"rextra-backend/internal/utils"

	"github.com/google/uuid"

	"os"
	"time"

	"gorm.io/gorm"
)

type (
	AuthService interface {
		Register(ctx context.Context, req dto_request.RegisterRequest) (dto_response.RegisterResponse, error)
		Login(ctx context.Context, req dto_request.LoginRequest) (dto_response.LoginResponse, error)
		Verify(ctx context.Context, authtoken string) error
		// ForgotPassword(ctx context.Context, req authDto.ForgotPasswordRequest) error
		// ChangePassword(ctx context.Context, req authDto.ChangePasswordRequest) error
		GetMe(ctx context.Context, userId string) (dto_response.GetMe, error)
		LoginWithGoogle(ctx context.Context, code, state string) (dto_response.LoginWithGoogleResponse, error)
	}

	authService struct {
		userRepository    repository.UserRepository
		sessionRepository repository.SessionRepository
		mailService       mailer.Mailer
		oauthService      oauth.Oauth
		db                *gorm.DB
	}
)

func NewAuth(userRepository repository.UserRepository,
	sessionRepository repository.SessionRepository,
	mailService mailer.Mailer,
	oauthService oauth.Oauth,
	db *gorm.DB) AuthService {
	return &authService{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		mailService:       mailService,
		oauthService:      oauthService,
		db:                db,
	}
}

func (s *authService) Register(ctx context.Context, req dto_request.RegisterRequest) (dto_response.RegisterResponse, error) {
	_, err := s.userRepository.GetByEmail(ctx, nil, req.Email)
	if err == nil {
		return dto_response.RegisterResponse{}, myerror.New("user with this email already exist", myerror.Error_RecordAlreadyExist)
	}

	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return dto_response.RegisterResponse{}, myerror.ProcessingError(err)
	}

	userCreation := entity.User{
		Fullname:    req.Fullname,
		Email:       req.Email,
		Password:    hashPassword,
		PhoneNumber: req.PhoneNumber,
	}

	createResult, err := s.userRepository.Create(ctx, nil, userCreation)
	if err != nil {
		return dto_response.RegisterResponse{}, err
	}

	token, err := myjwt.GenerateToken(map[string]string{
		"user_id": createResult.ID.String(),
		"email":   createResult.Email,
	}, 24*time.Hour)
	if err != nil {
		return dto_response.RegisterResponse{}, err
	}

	token = fmt.Sprintf("%s/auth/verify?token=%s", os.Getenv("APP_URL"), token)
	if err := s.mailService.MakeMail("./internal/pkg/email/template/verification_email.html", map[string]any{
		"Fullname": createResult.Fullname,
		"Verify":   token,
	}).Send(createResult.Email, "Verify Your Account").Error; err != nil {
		return dto_response.RegisterResponse{}, err
	}

	return dto_response.RegisterResponse{
		ID:          createResult.ID.String(),
		Fullname:    createResult.Fullname,
		Email:       createResult.Email,
		PhoneNumber: createResult.PhoneNumber,
		Role:        string(createResult.Role),
	}, nil
}

func (s *authService) Verify(ctx context.Context, token string) error {
	payloadToken, err := myjwt.GetPayloadInsideToken(token)
	if err != nil {
		return err
	}

	user, err := s.userRepository.GetByEmail(ctx, nil, payloadToken["email"])
	if err != nil {
		return err
	}

	user.IsVerified = true

	_, err = s.userRepository.Update(ctx, nil, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *authService) Login(ctx context.Context, req dto_request.LoginRequest) (dto_response.LoginResponse, error) {
	user, err := s.userRepository.GetByEmail(ctx, nil, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.LoginResponse{}, myerror.InvalidCreds()
		}
		return dto_response.LoginResponse{}, err
	}

	if !user.IsVerified {
		return dto_response.LoginResponse{}, myerror.New("user is not verify", myerror.Error_Unauthorized)
	}

	checkPassword, err := utils.CheckPassword(user.Password, []byte(req.Password))
	if !checkPassword || err != nil {
		return dto_response.LoginResponse{}, myerror.InvalidCreds()
	}

	accessToken, err := myjwt.GenerateToken(map[string]string{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    string(user.Role),
	}, 24*time.Hour)
	if err != nil {
		return dto_response.LoginResponse{}, err
	}

	// create refresh token
	token, err := utils.RandomData(32)
	if err != nil {
		return dto_response.LoginResponse{}, myerror.ProcessingError(err)
	}

	refreshToken, err := s.sessionRepository.Create(ctx, nil, entity.SessionToken{
		UserID:       user.ID.String(),
		Token:        token,
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour),
		IsActive:     true,
		AuthProvider: "local",
		DeviceInfo:   nil,
	})
	if err != nil {
		return dto_response.LoginResponse{}, myerror.DatabaseError(err)
	}

	return dto_response.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
		Role:         string(user.Role),
	}, nil
}

func (s *authService) GetMe(ctx context.Context, userId string) (dto_response.GetMe, error) {
	user, err := s.userRepository.GetById(ctx, nil, userId)
	if err != nil {
		return dto_response.GetMe{}, err
	}

	return dto_response.GetMe{
		PersonalInfo: dto_response.PersonalInfo{
			ID:          userId,
			Fullname:    user.Fullname,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			Role:        string(user.Role),
		},
	}, nil
}

func (s *authService) LoginWithGoogle(ctx context.Context, code, state string) (dto_response.LoginWithGoogleResponse, error) {
	tokenOauth, err := s.oauthService.Config.Exchange(ctx, code)
	if err != nil {
		return dto_response.LoginWithGoogleResponse{}, err
	}

	userInfo, err := s.oauthService.GetUserInfo(tokenOauth)
	if err != nil {
		return dto_response.LoginWithGoogleResponse{}, err
	}

	registerToken := ""
	needRegistration := false
	id := uuid.New()
	user, err := s.userRepository.GetByEmail(ctx, nil, userInfo.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto_response.LoginWithGoogleResponse{}, err
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		if !userInfo.VerifiedEmail {
			return dto_response.LoginWithGoogleResponse{}, myerror.New("user with this email not verified", myerror.Error_Unauthorized)
		} else {
			needRegistration = true
			registerToken, err = myjwt.GenerateToken(map[string]string{
				"user_id": id.String(),
				"email":   userInfo.Email,
				"state":   state,
			}, 10*time.Minute)
			if err != nil {
				return dto_response.LoginWithGoogleResponse{}, err
			}
		}
	}

	if !needRegistration {
		id = user.ID
	}

	token, err := myjwt.GenerateToken(map[string]string{
		"user_id": id.String(),
		"email":   user.Email,
	}, 24*time.Hour)
	if err != nil {
		return dto_response.LoginWithGoogleResponse{}, err
	}

	return dto_response.LoginWithGoogleResponse{
		Token:         token,
		Role:          string(user.Role),
		RegisterToken: registerToken,
	}, nil
}
