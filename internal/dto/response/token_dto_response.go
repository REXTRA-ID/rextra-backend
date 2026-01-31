package dto_response

import (
	"time"

	"github.com/google/uuid"
)

type (
	TokenWalletDTOResponse struct {
		UserId  string `json:"user_id"`
		Balance int64  `json:"balance"`
	}

	TokenLedgerDTOResponse struct {
		ID         string    `json:"id"`
		OccurredAt time.Time `json:"occurred_at"`
		Amount     int64     `json:"amount"`
		Direction  string    `json:"direction"`
		SourceType string    `json:"source_type"`
	}

	TokenLedgerDetailDTOResponse struct {
		ID            string     `json:"id"`
		OccurredAt    time.Time  `json:"occurred_at"`
		Direction     string     `json:"direction"`
		Amount        int64      `json:"amount"`
		BalanceBefore int64      `json:"balance_before"`
		BalanceAfter  int64      `json:"balance_after"`
		SourceType    string     `json:"source_type"`
		SourceID      *uint64    `json:"source_id,omitempty"`
		ReferenceID   *string    `json:"reference_id,omitempty"`
		Description   string     `json:"description"`
		Metadata      any        `json:"metadata"`
		OperatorID    *uuid.UUID `json:"operator_id,omitempty"`
		CreatedAt     time.Time  `json:"created_at"`
	}

	TokenBundleDTOResponse struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		TokenAmount  int64  `json:"token_amount"`
		PriceRp      int64  `json:"price_rp"`
		Label        string `json:"label"`
		DisplayOrder int    `json:"display_order"`
		IsActive     bool   `json:"is_available"`
	}
)
