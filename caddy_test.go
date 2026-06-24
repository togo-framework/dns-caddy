package caddy

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/togo-framework/dns"
)

func TestUpsertProxyHostPostsRoute(t *testing.T) {
	var posted map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			json.NewDecoder(r.Body).Decode(&posted)
		}
		w.WriteHeader(200)
	}))
	defer srv.Close()

	p := &provider{admin: srv.URL, server: "srv0", hc: srv.Client()}
	id, err := p.UpsertProxyHost(context.Background(), dns.ProxyHost{Domain: "app.example.com", Upstream: "http://127.0.0.1:8080"})
	if err != nil {
		t.Fatalf("upsert: %v", err)
	}
	if id == "" {
		t.Fatal("empty id")
	}
	h := posted["handle"].([]any)[0].(map[string]any)
	up := h["upstreams"].([]any)[0].(map[string]any)
	if up["dial"] != "127.0.0.1:8080" {
		t.Fatalf("dial=%v", up["dial"])
	}
}

func TestDialFor(t *testing.T) {
	if dialFor("https://svc.internal") != "svc.internal:443" {
		t.Fatal(dialFor("https://svc.internal"))
	}
	if dialFor("http://h:9000") != "h:9000" {
		t.Fatal(dialFor("http://h:9000"))
	}
}
