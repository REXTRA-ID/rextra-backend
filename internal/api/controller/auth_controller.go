package controller

import (
	"rextra-backend/internal/api/service"
	dto_request "rextra-backend/internal/dto/request"
	myerror "rextra-backend/internal/pkg/error"
	myjwt "rextra-backend/internal/pkg/jwt"
	"rextra-backend/internal/pkg/response"
	"rextra-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type (
	AuthController interface {
		Register(ctx *gin.Context)
		Login(ctx *gin.Context)
		Verify(ctx *gin.Context)
		SendVerificationEmail(ctx *gin.Context)
		ForgetPassword(ctx *gin.Context)
		ChangePassword(ctx *gin.Context)
		Me(ctx *gin.Context)
		LoginWithGoogle(ctx *gin.Context)
		Logout(ctx *gin.Context)
	}

	authController struct {
		authService service.AuthService
	}
)

func NewAuth(authService service.AuthService) AuthController {
	return &authController{
		authService: authService,
	}
}

func (c *authController) Register(ctx *gin.Context) {
	var req dto_request.RegisterRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	user, err := c.authService.Register(ctx, req)
	if err != nil {
		response.NewFailed("failed register account", err).Send(ctx)
		return
	}

	response.NewSuccess("success register account", user).Send(ctx)
}

func (c *authController) Login(ctx *gin.Context) {
	var req dto_request.LoginRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	result, err := c.authService.Login(ctx.Request.Context(), req)
	if err != nil {
		response.NewFailed("failed login", err).Send(ctx)
		return
	}

	response.NewSuccess("success login", result).Send(ctx)
}

func (c *authController) Verify(ctx *gin.Context) {
	token := ctx.Query("token")
	if err := c.authService.Verify(ctx, token); err != nil {
		response.NewFailed("failed verify account", err).Send(ctx)
		return
	}

	response.NewSuccess("success verify account", nil).Send(ctx)
}

func (h *authController) SendVerificationEmail(ctx *gin.Context) {
	var req dto_request.SendVerificationRequest

	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	if err := h.authService.SendVerificationEmail(ctx, req.Email); err != nil {
		response.NewFailed("failed send verification email", err).Send(ctx)
		return
	}

	response.NewSuccess("success send verification email", nil).Send(ctx)
}

func (c *authController) ForgetPassword(ctx *gin.Context) {
	var req dto_request.ForgetPasswordRequest

	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	if err := c.authService.ForgetPassword(ctx, req); err != nil {
		response.NewFailed("failed forget password", err).Send(ctx)
		return
	}

	response.NewSuccess("success forget password", nil).Send(ctx)
}

func (c *authController) ChangePassword(ctx *gin.Context) {
	var req dto_request.ChangePasswordRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", err).Send(ctx)
		return
	}

	token := ctx.Query("token")
	if token == "" {
		response.NewFailed("failed change password", myerror.ErrBodyRequest).Send(ctx)
		return
	}

	claims, err := myjwt.GetPayloadInsideToken(token)
	if err != nil {
		response.NewFailed("failed change password", err).Send(ctx)
		return
	}

	req.Email = claims["email"]
	if err := c.authService.ChangePassword(ctx, req); err != nil {
		response.NewFailed("failed change password", err).Send(ctx)
		return
	}

	response.NewSuccess("success change password", nil).Send(ctx)
}

func (c *authController) Me(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get user id", err).Send(ctx)
		return
	}

	res, err := c.authService.GetMe(ctx.Request.Context(), userId)
	if err != nil {
		response.NewFailed("failed get me", err).Send(ctx)
		return
	}

	response.NewSuccess("success get me", res).Send(ctx)
}

func (c *authController) LoginWithGoogle(ctx *gin.Context) {
	var req dto_request.LoginWithGoogleRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	result, err := c.authService.LoginWithGoogle(ctx.Request.Context(), req.IdToken)
	if err != nil {
		response.NewFailed("failed login with google", err).Send(ctx)
		return
	}

	response.NewSuccess("success login with google", result).Send(ctx)
}

func (c *authController) Logout(ctx *gin.Context) {
	var req dto_request.LogoutRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	if err := c.authService.Logout(ctx.Request.Context(), req); err != nil {
		response.NewFailed("failed logout", err).Send(ctx)
		return
	}

	response.NewSuccess("success logout", nil).Send(ctx)
}
