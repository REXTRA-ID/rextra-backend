package utils

import (
	"fmt"
	"time"
)

const PAYMENTEXTERNALID = `pay_%s_%s_%s`

func PaymentExternalID(paymentType string, userId string) string {
	xenditId := fmt.Sprintf(PAYMENTEXTERNALID, paymentType, userId, time.Now().UTC().String())
	return xenditId
}

func CalculateMembershipPrice(basePrice float64, duration int) float64 {
	return basePrice * float64(duration)
}
