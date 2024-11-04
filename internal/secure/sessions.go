package secure

import (
	crand "crypto/rand"
	"encoding/base64"
)

func GenerateSessionID() (string, error) {
	b := make([]byte, 32) // 256-bit session ID
	_, err := crand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
