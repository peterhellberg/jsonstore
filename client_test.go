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
			{Timeout: 1},
			{Timeout: 2},
			{Timeout: 3},
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

func TestSecret(t *testing.T) {
	s, err := NewSecret()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	c := New(Secret(s))

	if got, want := c.Secret(), s; got != want {
		t.Fatalf("c.Secret() = %q, want %q", got, want)
	}
}

func TestURL(t *testing.T) {
	s, err := NewSecret()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	c := New(Secret(s), BaseURL("https://example.org"))

	want := "https://example.org/" + s + "/foo/bar"

	if got := c.URL("foo", "bar").String(); got != want {
		t.Fatalf(`c.URL("foo", "bar").String() = %q, want %q`, got, want)
	}
}

func TestGet(t *testing.T) {
	c, ts := testClientAndServer(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]int{
			"result": 6543,
		})
	})
	defer ts.Close()

	var v int

	c.Get(context.Background(), "", &v)

	if got, want := v, 6543; got != want {
		t.Fatalf(`v = %d, want %d`, got, want)
	}
}

func TestPost(t *testing.T) {
	c, ts := testClientAndServer(func(w http.ResponseWriter, r *http.Request) {
	})
	defer ts.Close()

	if err := c.Post(context.Background(), "test", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPut(t *testing.T) {
	c, ts := testClientAndServer(func(w http.ResponseWriter, r *http.Request) {
	})
	defer ts.Close()

	if err := c.Put(context.Background(), "test", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDelete(t *testing.T) {
	c, ts := testClientAndServer(func(w http.ResponseWriter, r *http.Request) {
	})
	defer ts.Close()

	if err := c.Delete(context.Background(), "test"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRequest(t *testing.T) {
	t.Run("with no secret", func(t *testing.T) {
		c := &Client{}

		if _, err := c.request(context.Background(), "", "", nil); err != ErrNoSecret {
			t.Fatalf("expected ErrNoSecret, got %v", err)
		}
	})
}

func testClientAndServer(h http.HandlerFunc, options ...Option) (*Client, *httptest.Server) {
	ts := httptest.NewServer(h)
	return New(append(options, BaseURL(ts.URL))...), ts
}
