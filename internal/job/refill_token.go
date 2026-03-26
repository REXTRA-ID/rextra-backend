package job

import "rextra-backend/internal/modules/token_usage/service"

type RefillTokenJob struct {
	TokenTransactionService service.TokenTransactionService
}

func (j *RefillTokenJob) Run() {
	j.TokenTransactionService.RefillToken()
}
