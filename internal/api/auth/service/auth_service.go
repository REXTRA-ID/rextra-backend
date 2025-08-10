package authService

import (
	"context"
	"errors"
	"fmt"

	authDto "rextra-backend/internal/api/auth/dto/request"
	authDto_request "rextra-backend/internal/api/auth/dto/request"
	authDto_response "rextra-backend/internal/api/auth/dto/response"
	userRepository "rextra-backend/internal/api/user/repository"
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
		Register(ctx context.Context, req authDto_request.RegisterRequest, token string) (authDto_response.RegisterResponse, error)
		Login(ctx context.Context, req authDto.LoginRequest) (authDto_response.LoginResponse, error)
		Verify(ctx context.Context, authtoken string) error
		// ForgotPassword(ctx context.Context, req authDto.ForgotPasswordRequest) error
		// ChangePassword(ctx context.Context, req authDto.ChangePasswordRequest) error
		GetMe(ctx context.Context, userId string) (authDto_response.GetMe, error)
		LoginWithGoogle(ctx context.Context, code, state string) (authDto_response.LoginWithGoogleResponse, error)
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

func (s *authService) Register(ctx context.Context, req authDto_request.RegisterRequest, authtoken string) (authDto_response.RegisterResponse, error) {
	_, err := s.userRepository.GetByEmail(ctx, nil, req.Email)
	if err == nil {
		return authDto_response.RegisterResponse{}, myerror.New("user with this email already exist", myerror.Error_RecordAlreadyExist)
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
			return authDto_response.RegisterResponse{}, myerror.InvalidToken()
		} else if payload["email"] == "" || payload["email"] != req.Email {
			return authDto_response.RegisterResponse{}, myerror.InvalidRequest(
				fmt.Errorf("email not match with token payload"),
			)
		}

		userCreation.ID = uuid.MustParse(payload["user_id"])
		userCreation.IsVerified = true
	}

	createResult, err := s.userRepository.Create(ctx, nil, userCreation)
	if err != nil {
		return authDto_response.RegisterResponse{}, err
	}

	if authtoken == "" {
		token, err := myjwt.GenerateToken(map[string]string{
			"user_id": createResult.ID.String(),
			"email":   createResult.Email,
		}, 24*time.Hour)
		if err != nil {
			return authDto_response.RegisterResponse{}, err
		}

		token = fmt.Sprintf("%s/auth/verify?token=%s", os.Getenv("APP_URL"), token)
		if err := s.mailService.MakeMail("./internal/pkg/email/template/verification_email.html", map[string]any{
			"Username": createResult.Username,
			"Verify":   token,
		}).Send(createResult.Email, "Verify Your Account").Error; err != nil {
			return authDto_response.RegisterResponse{}, err
		}
	}

	return authDto_response.RegisterResponse{
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

func (s *authService) Login(ctx context.Context, req authDto.LoginRequest) (authDto_response.LoginResponse, error) {
	user, err := s.userRepository.GetByEmail(ctx, nil, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return authDto_response.LoginResponse{}, myerror.InvalidCreds()
		}
		return authDto_response.LoginResponse{}, err
	}

	if !user.IsVerified {
		return authDto_response.LoginResponse{}, myerror.New("user is not verify", myerror.Error_Unauthorized)
	}

	checkPassword, err := utils.CheckPassword(user.Password, []byte(req.Password))
	if !checkPassword || err != nil {
		return authDto_response.LoginResponse{}, myerror.InvalidCreds()
	}

	token, err := myjwt.GenerateToken(map[string]string{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    string(user.Role),
	}, 24*time.Hour)
	if err != nil {
		return authDto_response.LoginResponse{}, err
	}

	return authDto_response.LoginResponse{
		Token: token,
		Role:  string(user.Role),
	}, nil
}

func (s *authService) GetMe(ctx context.Context, userId string) (authDto_response.GetMe, error) {
	user, err := s.userRepository.GetById(ctx, nil, userId)
	if err != nil {
		return authDto_response.GetMe{}, err
	}

	return authDto_response.GetMe{
		PersonalInfo: authDto_response.PersonalInfo{
			ID:          userId,
			Username:    user.Username,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			Role:        string(user.Role),
		},
	}, nil
}

func (s *authService) LoginWithGoogle(ctx context.Context, code, state string) (authDto_response.LoginWithGoogleResponse, error) {
	tokenOauth, err := s.oauthService.Config.Exchange(ctx, code)
	if err != nil {
		return authDto_response.LoginWithGoogleResponse{}, err
	}

	userInfo, err := s.oauthService.GetUserInfo(tokenOauth)
	if err != nil {
		return authDto_response.LoginWithGoogleResponse{}, err
	}

	registerToken := ""
	needRegistration := false
	id := uuid.New()
	user, err := s.userRepository.GetByEmail(ctx, nil, userInfo.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return authDto_response.LoginWithGoogleResponse{}, err
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		if !userInfo.VerifiedEmail {
			return authDto_response.LoginWithGoogleResponse{}, myerror.New("user with this email not verified", myerror.Error_Unauthorized)
		} else {
			needRegistration = true
			registerToken, err = myjwt.GenerateToken(map[string]string{
				"user_id": id.String(),
				"email":   userInfo.Email,
				"state":   state,
			}, 10*time.Minute)
			if err != nil {
				return authDto_response.LoginWithGoogleResponse{}, err
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
		return authDto_response.LoginWithGoogleResponse{}, err
	}

	return authDto_response.LoginWithGoogleResponse{
		NeedRegistration: needRegistration,
		Token:            token,
		Role:             string(user.Role),
		RegisterToken:    registerToken,
	}, nil
}
