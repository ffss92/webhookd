package webhook

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestNewSecret(t *testing.T) {
	secret := NewSecret()

	prefix := "whsec_"
	encoded, found := strings.CutPrefix(secret, prefix)
	if !found {
		t.Fatal("expected secret to have \"whsec_\" prefix")
	}

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("failed to base64 decode the encoded part: %v", err)
	}
	if len(decoded) < 24 || len(decoded) > 64 {
		t.Fatalf("expected the decoded part to have between 24 and 64 bytes but got: %d", len(decoded))
	}
}
