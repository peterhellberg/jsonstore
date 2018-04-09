package jsonstore

import "testing"

func TestNewSecret(t *testing.T) {
	s, err := NewSecret()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got, want := len(s), 64; got != want {
		t.Fatalf("len(s) = %d, want %d", got, want)
	}
}
