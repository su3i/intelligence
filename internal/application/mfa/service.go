package mfa

import (
	"fmt"
	"time"

	"github.com/darksuei/suei-intelligence/internal/domain/mfa"
)

var ISSUER = "suei-intelligence"

func RetrieveTotpURI(email string, secret string) (string, error) {
	uri := fmt.Sprintf(
		"otpauth://totp/%s:%s?secret=%s&issuer=%s",
		ISSUER,
		email,
		secret,
		ISSUER,
	)
	return uri, nil
}

func VerifyTOTP(secret string, userCode uint32, now time.Time) bool {
	// Allow ±1 time window (30s) for clock drift
	for i := -1; i <= 1; i++ {
		t := now.Add(time.Duration(i*30) * time.Second)
		expected, err := mfa.GenerateTOTP(secret, t)
		if err != nil {
			continue
		}

		if mfa.ConstantTimeCompare(expected, userCode) {
			return true
		}
	}
	return false
}
