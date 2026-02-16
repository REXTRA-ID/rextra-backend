package service

import (
	"context"
	"errors"
	"fmt"

	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	authen "rextra-backend/internal/modules/auth/repository"
	membership "rextra-backend/internal/modules/membership/repository"
	user "rextra-backend/internal/modules/user/repository"
	mailer "rextra-backend/internal/pkg/email"
	myerror "rextra-backend/internal/pkg/error"
	myjwt "rextra-backend/internal/pkg/jwt"
	"rextra-backend/internal/utils"

	"os"
	"time"

	"firebase.google.com/go/v4/auth"
	"gorm.io/gorm"
)

type (
	AuthService interface {
		Register(ctx context.Context, req dto_request.RegisterRequest) (dto_response.RegisterResponse, error)
		Login(ctx context.Context, req dto_request.LoginRequest) (dto_response.LoginResponse, error)
		Verify(ctx context.Context, authtoken string) error
		ForgetPassword(ctx context.Context, req dto_request.ForgetPasswordRequest) error
		ChangePassword(ctx context.Context, req dto_request.ChangePasswordRequest) error
		GetMe(ctx context.Context, userId string) (dto_response.GetMe, error)
		SendVerificationEmail(ctx context.Context, email string) error
		// LoginWithGoogle(ctx context.Context, idToken string) (dto_response.LoginResponse, error)
		Logout(ctx context.Context, req dto_request.LogoutRequest) error
	}

	authService struct {
		userRepository               user.UserRepository
		membershipRepository         membership.MembershipRepository
		membershipDurationRepository membership.MembershipDurationRepository
		membershipPlanRepository     membership.MembershipPlanRepository
		sessionRepository            authen.SessionRepository
		mailService                  mailer.Mailer
		// firebaseClient               *auth.Client
		db *gorm.DB
	}
)

func NewAuth(userRepository user.UserRepository,
	membershipRepository membership.MembershipRepository,
	membershipDurationRepository membership.MembershipDurationRepository,
	membershipPlanRepository membership.MembershipPlanRepository,
	sessionRepository authen.SessionRepository,
	mailService mailer.Mailer,
	firebaseClient *auth.Client,
	db *gorm.DB) AuthService {
	return &authService{
		userRepository:               userRepository,
		membershipRepository:         membershipRepository,
		membershipDurationRepository: membershipDurationRepository,
		membershipPlanRepository:     membershipPlanRepository,
		sessionRepository:            sessionRepository,
		mailService:                  mailService,
		db:                           db,
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

	starterPlan, err := s.membershipPlanRepository.GetByPlanName(ctx, nil, string(entity.PLANSTARTER))
	if err != nil {
		return err
	}

	const DURATION_STARTER = 30

	membershipDuration, err := s.membershipDurationRepository.GetByDurationMonth(ctx, nil, DURATION_STARTER)
	if err != nil {
		return err
	}

	userMembership := entity.NewMembership(user.ID, &starterPlan, &membershipDuration)

	_, err = s.membershipRepository.Create(ctx, nil, userMembership)
	if err != nil {
		return err
	}

	user.IsVerified = true
	user.Membership = userMembership

	_, err = s.userRepository.Update(ctx, nil, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *authService) ForgetPassword(ctx context.Context, req dto_request.ForgetPasswordRequest) error {
	user, err := s.userRepository.GetByEmail(ctx, nil, req.Email)
	if err != nil {
		return err
	}

	if !user.IsVerified {
		return errors.New("user not verified")
	}

	token, err := myjwt.GenerateToken(map[string]string{
		"user_id": user.ID.String(),
		"email":   user.Email,
	}, 24*time.Hour)
	if err != nil {
		return err
	}

	// generate token
	token = fmt.Sprintf("%s/auth/change?token=%s", os.Getenv("APP_URL"), token)
	if err := s.mailService.MakeMail("./internal/pkg/email/template/forget_password_email.html", map[string]any{
		"Fullname": user.Fullname,
		"Link":     token,
	}).Send(user.Email, "Forget Password").Error; err != nil {
		return err
	}

	return nil
}

func (s *authService) ChangePassword(ctx context.Context, req dto_request.ChangePasswordRequest) error {
	user, err := s.userRepository.GetByEmail(ctx, nil, req.Email)
	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

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

func (s *authService) SendVerificationEmail(ctx context.Context, email string) error {
	user, err := s.userRepository.GetByEmail(ctx, nil, email)
	if err == gorm.ErrRecordNotFound || err != nil {
		return err
	}

	token, err := myjwt.GenerateToken(map[string]string{
		"user_id": user.ID.String(),
		"email":   user.Email,
	}, 24*time.Hour)
	if err != nil {
		return err
	}

	// generate token
	token = fmt.Sprintf("%s/auth/verify?token=%s", os.Getenv("APP_URL"), token)
	if err := s.mailService.MakeMail("./internal/pkg/email/template/verification_email.html", map[string]any{
		"Fullname": user.Fullname,
		"Verify":   token,
	}).Send(user.Email, "Verify Your Account").Error; err != nil {
		return err
	}

	return nil
}

// func (s *authService) LoginWithGoogle(ctx context.Context, idToken string) (dto_response.LoginResponse, error) {
// 	authToken, err := s.firebaseClient.VerifyIDToken(ctx, idToken)
// 	if err != nil {
// 		return dto_response.LoginResponse{}, myerror.InvalidToken()
// 	}

// 	uid := authToken.UID

// 	userRecord, err := s.firebaseClient.GetUser(ctx, uid)
// 	if err != nil {
// 		return dto_response.LoginResponse{}, myerror.ProcessingError(err)
// 	}

// 	user, err := s.userRepository.GetByEmail(ctx, nil, userRecord.Email)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			phoneNumber, err := myfirebase.GetGooglePhoneNumber(idToken)
// 			if err != nil {
// 				return dto_response.LoginResponse{}, myerror.ProcessingError(err)
// 			}

// 			user, err = s.userRepository.Create(ctx, nil, entity.User{
// 				Fullname:    userRecord.DisplayName,
// 				Email:       userRecord.Email,
// 				PhoneNumber: phoneNumber,
// 				Password:    uuid.New().String(),
// 				IsVerified:  true,
// 				Role:        entity.RoleUser,
// 			})
// 			if err != nil {
// 				return dto_response.LoginResponse{}, myerror.DatabaseError(err)
// 			}
// 		} else {
// 			return dto_response.LoginResponse{}, myerror.DatabaseError(err)
// 		}
// 	}

// 	accessToken, err := myjwt.GenerateToken(map[string]string{
// 		"user_id": user.ID.String(),
// 		"email":   user.Email,
// 		"role":    string(user.Role),
// 	}, 24*time.Hour)
// 	if err != nil {
// 		return dto_response.LoginResponse{}, err
// 	}

// 	// create refresh token
// 	token, err := utils.RandomData(32)
// 	if err != nil {
// 		return dto_response.LoginResponse{}, myerror.ProcessingError(err)
// 	}

// 	refreshToken, err := s.sessionRepository.Create(ctx, nil, entity.SessionToken{
// 		UserID:       user.ID.String(),
// 		Token:        token,
// 		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour),
// 		IsActive:     true,
// 		AuthProvider: authToken.Firebase.SignInProvider,
// 		DeviceInfo:   nil,
// 	})
// 	if err != nil {
// 		return dto_response.LoginResponse{}, myerror.DatabaseError(err)
// 	}

// 	return dto_response.LoginResponse{
// 		AccessToken:  accessToken,
// 		RefreshToken: refreshToken.Token,
// 		Role:         string(user.Role),
// 	}, nil
// }

func (s *authService) Logout(ctx context.Context, req dto_request.LogoutRequest) error {
	session, err := s.sessionRepository.GetByToken(ctx, nil, req.RefreshToken)
	if err != nil {
		return myerror.DatabaseError(err)
	}

	if !session.IsActive {
		return myerror.New("session already inactive", myerror.Error_InvalidRequest)
	}

	session.IsActive = false
	_, err = s.sessionRepository.Update(ctx, nil, session)
	if err != nil {
		return myerror.DatabaseError(err)
	}

	if err := s.sessionRepository.Delete(ctx, nil, session); err != nil {
		return myerror.DatabaseError(err)
	}

	return nil
}
