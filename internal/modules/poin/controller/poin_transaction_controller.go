package controller

import (
	"rextra-backend/internal/modules/poin/service"

	"github.com/gin-gonic/gin"
)

type (
	PoinTransactionController interface {
		UsePoin(ctx *gin.Context)
		EarnPoin(ctx *gin.Context)
		RedeemPoin(ctx *gin.Context)
	}

	poinTransaactionController struct {
		poinTransactionService service.PoinTransactionService
	}
)

func NewPoinTransactionController(poinTransactionService service.PoinTransactionService) PoinTransactionController {
	return &poinTransaactionController{
		poinTransactionService: poinTransactionService,
	}
}

func (c *poinTransaactionController) UsePoin(ctx *gin.Context) {

}

func (c *poinTransaactionController) EarnPoin(ctx *gin.Context) {

}

func (c *poinTransaactionController) RedeemPoin(ctx *gin.Context) {

}
