// +build integration

package jsonstore

import (
	"context"
	"testing"
)

func TestIntegration(t *testing.T) {
	ctx := context.Background()
	key := "test"

	// Create a new client with a random secret
	c := New()

	// Post some data to jsonstore.io
	if err := c.Post(ctx, key, testData{Foo: 123, Bar: false, Baz: "abc"}); err != nil {
		t.Fatalf("c.Post for %s returned error: %v", c.URL(key), err)
	}

	// Modify the Baz value
	if err := c.Put(ctx, key+"/Baz", "cde"); err != nil {
		t.Fatalf("c.Put for %s returned error: %v", c.URL(key+"/Baz"), err)
	}

	var td testData

	// Retrieve the test data from the key
	if err := c.Get(ctx, key, &td); err != nil {
		t.Fatalf("c.Get for %s returned error: %v", c.URL(key), err)
	}

	if got, want := td.Foo, 123; got != want {
		t.Fatalf("td.Foo = %d, want %d", got, want)
	}

	if got, want := td.Bar, false; got != want {
		t.Fatalf("td.Bar = %v, want %v", got, want)
	}

	if got, want := td.Baz, "cde"; got != want {
		t.Fatalf("td.Baz = %q, want %q", got, want)
	}

	// Delete all of the data
	if err := c.Delete(ctx, key); err != nil {
		t.Fatalf("c.Delete for %s returned error: %v", c.URL(key), err)
	}
}

type testData struct {
	Foo int
	Bar bool
	Baz string
}
