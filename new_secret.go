package jsonstore

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

// NewSecret generates a new secret based on 64 random bytes
func NewSecret() (string, error) {
	b, err := randomBytes(64)
	if err != nil {
		return "", err
	}

	return sha256sum(b), nil
}

func randomBytes(n int) ([]byte, error) {
	b := make([]byte, n)

	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}

func sha256sum(b []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(b))
}
