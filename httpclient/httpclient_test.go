package httpclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestGetAndJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test") != "yes" {
			http.Error(w, "missing header", http.StatusBadRequest)
			return
		}
		if got := r.URL.Query()["a"]; len(got) != 2 || got[0] != "1" || got[1] != "2" {
			http.Error(w, "bad query", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"value":7}`))
	}))
	defer server.Close()

	response, err := Get(server.URL, Options{
		Headers: map[string]string{"X-Test": "yes"},
		Query:   url.Values{"a": {"1", "2"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !response.OK() || response.StatusCode != 200 {
		t.Fatalf("response = %#v", response)
	}
	var payload struct {
		OK    bool `json:"ok"`
		Value int  `json:"value"`
	}
	if err := response.JSON(&payload); err != nil || !payload.OK || payload.Value != 7 {
		t.Fatalf("JSON = (%#v,%v)", payload, err)
	}

	var direct struct {
		Value int `json:"value"`
	}
	if err := GetJSON(server.URL, &direct, Options{
		Headers: map[string]string{"X-Test": "yes"},
		Query:   url.Values{"a": {"1", "2"}},
	}); err != nil || direct.Value != 7 {
		t.Fatalf("GetJSON = (%#v,%v)", direct, err)
	}
}

func TestPostJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "bad content type", http.StatusBadRequest)
			return
		}
		var payload map[string]any
		if err := decodeJSON(r, &payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if payload["name"] != "Tom" {
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}
		_, _ = w.Write([]byte("created"))
	}))
	defer server.Close()

	response, err := PostJSON(server.URL, map[string]any{"name": "Tom"})
	if err != nil {
		t.Fatal(err)
	}
	if response.Text() != "created" {
		t.Fatalf("Text = %q", response.Text())
	}

	response, err = Post(server.URL, []byte(`{"name":"Tom"}`), Options{Headers: map[string]string{"Content-Type": "application/json"}})
	if err != nil || response.Text() != "created" {
		t.Fatalf("Post = (%q,%v)", response.Text(), err)
	}
}

func TestTimeoutAndError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		_, _ = w.Write([]byte("slow"))
	}))
	defer server.Close()

	start := time.Now()
	_, err := Get(server.URL, Options{Timeout: 20 * time.Millisecond})
	if err == nil || !strings.Contains(err.Error(), "deadline") {
		t.Fatalf("timeout err = %v", err)
	}
	if time.Since(start) > time.Second {
		t.Fatal("timeout too slow")
	}
	if _, err := Get("://bad-url"); err == nil {
		t.Fatal("invalid URL should error")
	}
}

func decodeJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}
