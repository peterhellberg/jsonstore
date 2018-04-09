package jsonstore

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		c := New()

		if got, want := c.baseURL.Scheme, defaultScheme; got != want {
			t.Fatalf("c.baseURL.Scheme = %q, want %q", got, want)
		}

		if got, want := c.baseURL.Host, defaultHost; got != want {
			t.Fatalf("c.baseURL.Host = %q, want %q", got, want)
		}

		if got, want := len(c.secret), 64; got != want {
			t.Fatalf("len(c.secret) = %d, want %d", got, want)
		}
	})

	t.Run("HTTPClient", func(t *testing.T) {
		for _, hc := range []*http.Client{
			&http.Client{Timeout: 1},
			&http.Client{Timeout: 2},
			&http.Client{Timeout: 3},
		} {
			c := New(HTTPClient(hc))

			if got, want := c.httpClient.Timeout, hc.Timeout; got != want {
				t.Fatalf("c.httpClient.Timeout = %q, want %q", got, want)
			}
		}
	})

	t.Run("BaseURL", func(t *testing.T) {
		for _, rawurl := range []string{
			"example.com",
			"example.net",
			"example.org",
		} {
			c := New(BaseURL(rawurl))

			if got, want := c.baseURL.String(), rawurl; got != want {
				t.Fatalf("c.baseURL.String() = %q, want %q", got, want)
			}
		}
	})

	t.Run("Secret", func(t *testing.T) {
		for _, secret := range []string{
			"foo",
			"bar",
			"baz",
		} {
			c := New(Secret(secret))

			if got, want := c.secret, secret; got != want {
				t.Fatalf("c.secret = %q, want %q", got, want)
			}

			if got, want := c.Secret(), secret; got != want {
				t.Fatalf("c.Secret() = %q, want %q", got, want)
			}
		}
	})
}

func TestGet(t *testing.T) {
	ctx := context.Background()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]int{
			"test": 6543,
		})
	}))
	defer ts.Close()

	c := New(BaseURL(ts.URL))

	v := map[string]int{}

	c.Get(ctx, "", &v)

	if got, want := v["test"], 6543; got != want {
		t.Fatalf(`v["test"] = %d, want %d`, got, want)
	}
}
