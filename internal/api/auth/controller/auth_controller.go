package authController

import (
	"net/http"

	authservice "rextra-backend/internal/api/auth/service"
	"rextra-backend/internal/dto"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/google/oauth"
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
		ForgotPassword(ctx *gin.Context)
		ChangePassword(ctx *gin.Context)
		Me(ctx *gin.Context)
		LoginWithGoogle(ctx *gin.Context)
		CallbackGoogle(ctx *gin.Context)
	}

	authController struct {
		authService authservice.AuthService
	}
)

func New(authService authservice.AuthService) AuthController {
	return &authController{
		authService: authService,
	}
}

func (c *authController) Register(ctx *gin.Context) {
	token := ctx.Query("token")
	if token != "" {
		verified, err := myjwt.IsValid(token)
		if err != nil || !verified {
			response.NewFailed("failed get data from body", myerror.New("token invalid", http.StatusBadRequest)).Send(ctx)
			return
		}
	}

	var req dto.RegisterRequest

	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.New(err.Error(), http.StatusBadRequest)).Send(ctx)
		return
	}

	user, err := c.authService.Register(ctx, req, token)
	if err != nil {
		response.NewFailed("failed register account", err).Send(ctx)
		return
	}

	response.NewSuccess("success register account", user).Send(ctx)
}

func (c *authController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
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

func (c *authController) ForgotPassword(ctx *gin.Context) {

}

func (c *authController) ChangePassword(ctx *gin.Context) {

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
	googleOauthConfig := oauth.GetConfig()
	oauthState := oauth.RandomState()

	domain, secure := utils.GetDomain()
	ctx.SetCookie("oauthstate", oauthState, 300, "/", domain, secure, true)
	url := googleOauthConfig.AuthCodeURL(oauthState)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (c *authController) CallbackGoogle(ctx *gin.Context) {
	state := ctx.Query("state")
	stateFromCookie, _ := ctx.Cookie("oauthstate")
	if state != stateFromCookie {
		response.NewFailed("failed get login callback", myerror.New("invalid oauth state", http.StatusBadRequest)).Send(ctx)
		return
	}

	code := ctx.Query("code")
	result, err := c.authService.LoginWithGoogle(ctx, code, state)
	if err != nil {
		response.NewFailed("failed login with google", err).Send(ctx)
		return
	}

	response.NewSuccess("success login with google", result).Send(ctx)
}
