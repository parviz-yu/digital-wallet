package security

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
)

func generateSignature(secretToken string, payloadBody []byte) string {
	mac := hmac.New(sha1.New, []byte(secretToken))
	mac.Write(payloadBody)
	expectedMAC := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(expectedMAC)
}

func VerifyBody(secretToken string, payloadBody []byte, toCompareWith string) bool {
	signature := generateSignature(secretToken, payloadBody)
	return signature == toCompareWith
}
