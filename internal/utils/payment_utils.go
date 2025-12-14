package utils

import (
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
	return price - discount
}

func ApplyDiscountFixed(price float64, discountPercent int) float64 {
	return price - float64(discountPercent)
}
