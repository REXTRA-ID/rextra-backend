package utils

import (
	"fmt"
	"time"
)

const XENDITEXTERNALID = `xendit_%s_%s_%s`

func MakeXenditExternalID(paymentType string, userId string) string {
	xenditId := fmt.Sprintf(XENDITEXTERNALID, paymentType, userId, time.Now().UTC().String())
	return xenditId
}
