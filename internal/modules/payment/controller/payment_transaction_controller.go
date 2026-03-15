package controller

import (
	"encoding/json"
	"io"
	"strconv"

	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/internal/modules/payment/service"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"
	"rextra-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type (
	PaymentTransactionController interface {
		CreateTransaction(ctx *gin.Context)
		GetMyTransactions(ctx *gin.Context)
		GetMyTransactionDetail(ctx *gin.Context)
		GetPaymentChannels(ctx *gin.Context)
		CalculatePrice(ctx *gin.Context)
		GetAllTransactions(ctx *gin.Context)
		GetTransactionDetail(ctx *gin.Context)
		CancelTransaction(ctx *gin.Context)
		HandleTripayCallback(ctx *gin.Context)
		CallbackStatusCheck(ctx *gin.Context)
		SimulateTripayPayment(ctx *gin.Context)
	}

	paymentTransactionController struct {
		service service.PaymentTransactionService
	}
)

func NewPaymentTransactionController(svc service.PaymentTransactionService) PaymentTransactionController {
	return &paymentTransactionController{service: svc}
}

func (c *paymentTransactionController) CreateTransaction(ctx *gin.Context) {
	userID, _ := utils.GetUserIdFromCtx(ctx)
	userNameVal, _ := ctx.Get("user_name")
	userEmailVal, _ := ctx.Get("user_email")
	userName, _ := userNameVal.(string)
	userEmail, _ := userEmailVal.(string)

	var req dto_request.MakeNewTransactionMembershipRequest
	if err := ctx.ShouldBindJSON(&req); err != nil { response.NewFailed("invalid body", myerror.InvalidRequest(err)).Send(ctx); return }

	result, err := c.service.CreateTransaction(ctx, req, userID, userName, userEmail)
	if err != nil { response.NewFailed("failed create", err).Send(ctx); return }
	response.NewSuccess("created", result).Send(ctx)
}

func (c *paymentTransactionController) GetMyTransactions(ctx *gin.Context) {
	userID, _ := utils.GetUserIdFromCtx(ctx)
	var req dto_request.MyTransactionFilterRequest
	ctx.ShouldBindQuery(&req)
	result, err := c.service.GetMyTransactions(ctx, userID, req)
	if err != nil { response.NewFailed("failed get", err).Send(ctx); return }
	response.NewSuccess("retrieved", result).Send(ctx)
}

func (c *paymentTransactionController) GetMyTransactionDetail(ctx *gin.Context) {
	userID, _ := utils.GetUserIdFromCtx(ctx)
	result, err := c.service.GetTransactionDetail(ctx, ctx.Param("transactionId"), userID)
	if err != nil { response.NewFailed("failed get", err).Send(ctx); return }
	response.NewSuccess("retrieved", result).Send(ctx)
}

func (c *paymentTransactionController) GetPaymentChannels(ctx *gin.Context) {
	amountStr := ctx.DefaultQuery("amount", "0")
	amount, _ := strconv.ParseInt(amountStr, 10, 64)

	result, err := c.service.GetPaymentChannels(ctx, amount)
	if err != nil { response.NewFailed("failed get channels", err).Send(ctx); return }
	response.NewSuccess("retrieved", result).Send(ctx)
}

func (c *paymentTransactionController) CalculatePrice(ctx *gin.Context) {
	userID, _ := utils.GetUserIdFromCtx(ctx)
	var req dto_request.CalculatePriceRequest
	ctx.ShouldBindJSON(&req)
	result, err := c.service.CalculatePrice(ctx, req, userID)
	if err != nil { response.NewFailed("failed calc", err).Send(ctx); return }
	response.NewSuccess("calculated", result).Send(ctx)
}

func (c *paymentTransactionController) GetAllTransactions(ctx *gin.Context) {
	var req dto_request.PaymentTransactionFilterRequest
	ctx.ShouldBindQuery(&req)
	result, err := c.service.GetAllTransactions(ctx, req)
	if err != nil { response.NewFailed("failed get all", err).Send(ctx); return }
	response.NewSuccess("retrieved", result).Send(ctx)
}

func (c *paymentTransactionController) GetTransactionDetail(ctx *gin.Context) {
	result, err := c.service.GetTransactionDetail(ctx, ctx.Param("transactionId"), "")
	if err != nil { response.NewFailed("failed get", err).Send(ctx); return }
	response.NewSuccess("retrieved", result).Send(ctx)
}

func (c *paymentTransactionController) CancelTransaction(ctx *gin.Context) {
	var req dto_request.CancelTransactionRequest
	ctx.ShouldBindJSON(&req)
	if err := c.service.CancelTransaction(ctx, ctx.Param("transactionId"), req); err != nil { response.NewFailed("failed cancel", err).Send(ctx); return }
	response.NewSuccess("cancelled", nil).Send(ctx)
}

func (c *paymentTransactionController) HandleTripayCallback(ctx *gin.Context) {
	rawBody, _ := io.ReadAll(ctx.Request.Body)
	sig := ctx.GetHeader("X-Callback-Signature")
	var payload dto_request.TripayCallbackRequest
	json.Unmarshal(rawBody, &payload)
	c.service.HandleCallback(ctx, sig, rawBody, payload)
	response.NewSuccess("callback processed", nil).Send(ctx)
}

func (c *paymentTransactionController) CallbackStatusCheck(ctx *gin.Context) {
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>REXTRA Payment Gateway</title>
		<style>
			body { font-family: sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; background: #f4f7f6; }
			.card { background: white; padding: 40px; border-radius: 12px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); text-align: center; max-width: 400px; }
			.status-badge { display: inline-block; padding: 6px 12px; background: #2ecc71; color: white; border-radius: 20px; font-weight: bold; margin-bottom: 20px; }
			h1 { color: #2c3e50; margin-bottom: 10px; }
			p { color: #7f8c8d; line-height: 1.5; }
			.footer { margin-top: 30px; font-size: 12px; color: #bdc3c7; }
		</style>
	</head>
	<body>
		<div class="card">
			<div class="status-badge">ACTIVE</div>
			<h1>REXTRA Callback Handler</h1>
			<p>This endpoint is active and listening for TriPay payment notifications. Please use POST method for callback data.</p>
			<div class="footer">REXTRA Backend &copy; 2026</div>
		</div>
	</body>
	</html>
	`
	ctx.Header("Content-Type", "text/html")
	ctx.String(200, html)
}

func (c *paymentTransactionController) SimulateTripayPayment(ctx *gin.Context) {
	transactionID := ctx.Param("transactionId")
	if err := c.service.SimulateTripayPayment(ctx, transactionID); err != nil {
		response.NewFailed("failed to simulate payment", err).Send(ctx)
		return
	}
	response.NewSuccess("simulated success", nil).Send(ctx)
}
