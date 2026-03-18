package service

import "fmt"

type ErrPendingTransactionExists struct {
	TransactionID string
}

func (e *ErrPendingTransactionExists) Error() string {
	return fmt.Sprintf("kamu masih punya transaksi yang belum dibayar: %s", e.TransactionID)
}
