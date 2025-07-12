package webhook

import (
	"crypto/rand"
	"encoding/base64"
)

const (
	// Must be between 24 and 64 bytes according to the reference.
	secretSize = 32
)

func NewSecret() string {
	b := make([]byte, secretSize)
	_, _ = rand.Read(b) // Never fails according to docs
	secret := base64.StdEncoding.EncodeToString(b)
	return "whsec_" + secret
}
