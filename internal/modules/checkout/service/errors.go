package service

import "fmt"

type ErrPendingTransactionExists struct {
	TransactionID string
}

func (e *ErrPendingTransactionExists) Error() string {
	return fmt.Sprintf("anda memiliki transaksi pending yang belum dibayar: %s", e.TransactionID)
}
