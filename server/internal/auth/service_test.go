package auth

import (
	"encoding/base64"
	"testing"

	"golang.org/x/crypto/argon2"
)

func TestPasswordHashVerification(t *testing.T) {
	salt := []byte("0123456789abcdef")
	derived := argon2.IDKey([]byte("correct horse battery staple"), salt, 3, 64*1024, 2, 32)
	encoded := base64.RawStdEncoding.EncodeToString(append(salt, derived...))
	if !verifyPassword(encoded, "correct horse battery staple") { t.Fatal("expected password to verify") }
	if verifyPassword(encoded, "wrong password") { t.Fatal("unexpected password verification") }
}
