package authService

import (
	"context"
	"errors"
	"fmt"
	userRepository "rextra-backend/internal/api/user/repository"
	"rextra-backend/internal/dto"
	"rextra-backend/internal/entity"
	mailer "rextra-backend/internal/pkg/email"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/google/oauth"
	myjwt "rextra-backend/internal/pkg/jwt"
	"rextra-backend/internal/utils"

	"github.com/google/uuid"

	"net/http"
	"os"
	"time"

	"gorm.io/gorm"
)

type (
	AuthService interface {
		Register(ctx context.Context, req dto.RegisterRequest, token string) (dto.RegisterResponse, error)
		Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
		Verify(ctx context.Context, authtoken string) error
		// ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error
		// ChangePassword(ctx context.Context, req dto.ChangePasswordRequest) error
		GetMe(ctx context.Context, userId string) (dto.GetMe, error)
		LoginWithGoogle(ctx context.Context, code, state string) (dto.LoginWithGoogleResponse, error)
	}

	authService struct {
		userRepository userRepository.UserRepository
		mailService    mailer.Mailer
		oauthService   oauth.Oauth
		db             *gorm.DB
	}
)

func New(userRepository userRepository.UserRepository,
	mailService mailer.Mailer,
	oauthService oauth.Oauth,
	db *gorm.DB) AuthService {
	return &authService{
		userRepository: userRepository,
		mailService:    mailService,
		oauthService:   oauthService,
		db:             db,
	}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest, authtoken string) (dto.RegisterResponse, error) {
	_, err := s.userRepository.GetByEmail(ctx, nil, req.Email)
	if err == nil {
		return dto.RegisterResponse{}, myerror.New("user with this email already exist", http.StatusConflict)
	}

	userCreation := entity.User{
		Username:    req.Username,
		Email:       req.Email,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
	}

	if authtoken != "" {
		payload, err := myjwt.GetPayloadInsideToken(authtoken)
		if err != nil {
			return dto.RegisterResponse{}, myerror.New("failed get payload token", http.StatusBadRequest)
		} else if payload["email"] == "" || payload["email"] != req.Email {
			return dto.RegisterResponse{}, myerror.New("email not match with token payload", http.StatusBadRequest)
		}

		userCreation.ID = uuid.MustParse(payload["user_id"])
		userCreation.IsVerified = true
	}

	createResult, err := s.userRepository.Create(ctx, nil, userCreation)
	if err != nil {
		return dto.RegisterResponse{}, err
	}

	if authtoken == "" {
		token, err := myjwt.GenerateToken(map[string]string{
			"user_id": createResult.ID.String(),
			"email":   createResult.Email,
		}, 24*time.Hour)
		if err != nil {
			return dto.RegisterResponse{}, err
		}

		token = fmt.Sprintf("%s/auth/verify?token=%s", os.Getenv("APP_URL"), token)
		if err := s.mailService.MakeMail("./internal/pkg/email/template/verification_email.html", map[string]any{
			"Username": createResult.Username,
			"Verify":   token,
		}).Send(createResult.Email, "Verify Your Account").Error; err != nil {
			return dto.RegisterResponse{}, err
		}
	}

	return dto.RegisterResponse{
		ID:          createResult.ID.String(),
		Username:    createResult.Username,
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

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	user, err := s.userRepository.GetByEmail(ctx, nil, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.LoginResponse{}, myerror.New("email or password invalid", http.StatusBadRequest)
		}
		return dto.LoginResponse{}, err
	}

	if !user.IsVerified {
		return dto.LoginResponse{}, myerror.New("user is not verify", http.StatusUnauthorized)
	}

	checkPassword, err := utils.CheckPassword(user.Password, []byte(req.Password))
	if !checkPassword || err != nil {
		return dto.LoginResponse{}, myerror.New("email or password invalid", http.StatusBadRequest)
	}

	token, err := myjwt.GenerateToken(map[string]string{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    string(user.Role),
	}, 24*time.Hour)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{
		Token: token,
		Role:  string(user.Role),
	}, nil
}

func (s *authService) GetMe(ctx context.Context, userId string) (dto.GetMe, error) {
	user, err := s.userRepository.GetById(ctx, nil, userId)
	if err != nil {
		return dto.GetMe{}, err
	}

	return dto.GetMe{
		PersonalInfo: dto.PersonalInfo{
			ID:          userId,
			Username:    user.Username,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			Role:        string(user.Role),
		},
	}, nil
}

func (s *authService) LoginWithGoogle(ctx context.Context, code, state string) (dto.LoginWithGoogleResponse, error) {
	tokenOauth, err := s.oauthService.Config.Exchange(ctx, code)
	if err != nil {
		return dto.LoginWithGoogleResponse{}, err
	}

	userInfo, err := s.oauthService.GetUserInfo(tokenOauth)
	if err != nil {
		return dto.LoginWithGoogleResponse{}, err
	}

	registerToken := ""
	needRegistration := false
	id := uuid.New()
	user, err := s.userRepository.GetByEmail(ctx, nil, userInfo.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.LoginWithGoogleResponse{}, err
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		if !userInfo.VerifiedEmail {
			return dto.LoginWithGoogleResponse{}, myerror.New("user with this email not verified", http.StatusBadRequest)
		} else {
			needRegistration = true
			registerToken, err = myjwt.GenerateToken(map[string]string{
				"user_id": id.String(),
				"email":   userInfo.Email,
				"state":   state,
			}, 10*time.Minute)
			if err != nil {
				return dto.LoginWithGoogleResponse{}, err
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
		return dto.LoginWithGoogleResponse{}, err
	}

	return dto.LoginWithGoogleResponse{
		NeedRegistration: needRegistration,
		Token:            token,
		Role:             string(user.Role),
		RegisterToken:    registerToken,
	}, nil
}
