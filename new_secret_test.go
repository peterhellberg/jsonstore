package jsonstore

import "testing"

func TestNewSecret(t *testing.T) {
	secret, err := NewSecret()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got, want := len(secret), 64; got != want {
		t.Fatalf("len(secret) = %d, want %d", got, want)
	}
}
