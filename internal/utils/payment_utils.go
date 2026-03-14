package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"
)

const PAYMENTEXTERNALID = `pay_%s_%s_%s`

func PaymentExternalID(paymentType string, userId string) string {
	paymentExternaId := fmt.Sprintf(PAYMENTEXTERNALID, paymentType, userId, time.Now().UTC().String())
	return paymentExternaId
}

func GenerateMerchantReff() string {
	timestamp := time.Now().Format("20060102150405")
	randomBytes := make([]byte, 4)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic("failed to generate random merchant ref")
	}
	randomStr := hex.EncodeToString(randomBytes)
	return fmt.Sprintf("INV-%s-%s", timestamp, randomStr)
}

func GenerateTransactionID() string {
	date := time.Now().Format("20060102")
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 4)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return fmt.Sprintf("TRX-%s-%s", date, string(b))
}
