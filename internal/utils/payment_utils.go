package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

const PAYMENTEXTERNALID = `pay_%s_%s_%s`

func PaymentExternalID(paymentType string, userId string) string {
	paymentExternaId := fmt.Sprintf(PAYMENTEXTERNALID, paymentType, userId, time.Now().UTC().String())
	return paymentExternaId
}

func CalculateMembershipPrice(basePrice float64, duration int) float64 {
	return basePrice * float64(duration)
}

func ApplyDiscountPercentage(price float64, discountPercent int) float64 {
	discount := (price * float64(discountPercent)) / 100
	if price < discount {
		return price
	}
	return price - discount
}

func ApplyDiscountFixed(price float64, discountAmount int) float64 {
	if price < float64(discountAmount) {
		return price
	}
	return price - float64(discountAmount)
}

func GenerateMerchantReff() string {
	timestamp := time.Now().Format("20060102150405")

	// Generate 4 random bytes (8 hex chars)
	randomBytes := make([]byte, 4)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic("failed to generate random merchant ref")
	}

	randomStr := hex.EncodeToString(randomBytes)

	return fmt.Sprintf("INV-%s-%s", timestamp, randomStr)
}
