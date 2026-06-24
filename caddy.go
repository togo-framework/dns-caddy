// Package caddy is a togo dns driver for the Caddy reverse proxy, driven through
// Caddy's admin API. Proxy hosts and gateway routes are both modeled as Caddy
// HTTP routes (host matcher → reverse_proxy upstream). DNS records return
// ErrUnsupported.
//
// Install: `togo install togo-framework/dns-caddy`, set DNS_DRIVER=caddy.
package caddy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/togo-framework/dns"
	"github.com/togo-framework/togo"
)

func init() {
	dns.RegisterDriver("caddy", func(k *togo.Kernel) (dns.Provider, error) {
		admin := os.Getenv("CADDY_ADMIN")
		if admin == "" {
			admin = "http://localhost:2019"
		}
		server := os.Getenv("CADDY_SERVER")
		if server == "" {
			server = "srv0"
		}
		return &provider{admin: strings.TrimRight(admin, "/"), server: server, hc: &http.Client{Timeout: 15 * time.Second}}, nil
	})
}

type provider struct {
	admin, server string
	hc            *http.Client
}

func (p *provider) do(ctx context.Context, method, path string, body any) ([]byte, error) {
	var rdr io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		rdr = bytes.NewReader(b)
	}
	req, _ := http.NewRequestWithContext(ctx, method, p.admin+path, rdr)
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("dns-caddy: %s %s -> %d: %s", method, path, resp.StatusCode, string(raw))
	}
	return raw, nil
}

// dialFor strips the scheme so Caddy gets host:port to dial.
func dialFor(upstream string) string {
	if u, err := url.Parse(upstream); err == nil && u.Host != "" {
		if u.Port() == "" && u.Scheme == "https" {
			return u.Hostname() + ":443"
		}
		if u.Port() == "" {
			return u.Hostname() + ":80"
		}
		return u.Host
	}
	return upstream
}

func (p *provider) addRoute(ctx context.Context, id, domain, pathPrefix, upstream string) (string, error) {
	match := map[string]any{"host": []string{domain}}
	if pathPrefix != "" && pathPrefix != "/" {
		match["path"] = []string{strings.TrimRight(pathPrefix, "*") + "*"}
	}
	route := map[string]any{
		"@id":   id,
		"match": []any{match},
		"handle": []any{map[string]any{
			"handler":   "reverse_proxy",
			"upstreams": []any{map[string]any{"dial": dialFor(upstream)}},
		}},
	}
	path := fmt.Sprintf("/config/apps/http/servers/%s/routes", url.PathEscape(p.server))
	if _, err := p.do(ctx, http.MethodPost, path, route); err != nil {
		return "", err
	}
	return id, nil
}

func (p *provider) UpsertProxyHost(ctx context.Context, h dns.ProxyHost) (string, error) {
	id := "togo_" + strings.ReplaceAll(h.Domain, ".", "_")
	_ = p.do2DeleteByID(ctx, id) // replace-on-upsert
	return p.addRoute(ctx, id, h.Domain, "/", h.Upstream)
}

func (p *provider) UpsertRoute(ctx context.Context, rt dns.Route) (string, error) {
	id := "togo_" + strings.ReplaceAll(rt.Domain+rt.Path, "/", "_")
	id = strings.ReplaceAll(id, ".", "_")
	_ = p.do2DeleteByID(ctx, id)
	return p.addRoute(ctx, id, rt.Domain, rt.Path, rt.Upstream)
}

func (p *provider) do2DeleteByID(ctx context.Context, id string) error {
	_, err := p.do(ctx, http.MethodDelete, "/id/"+url.PathEscape(id), nil)
	return err
}

func (p *provider) DeleteProxyHost(ctx context.Context, id string) error {
	return p.do2DeleteByID(ctx, id)
}
func (p *provider) DeleteRoute(ctx context.Context, id string) error {
	return p.do2DeleteByID(ctx, id)
}

func (p *provider) UpsertRecord(context.Context, string, dns.Record) (string, error) {
	return "", dns.ErrUnsupported
}
func (p *provider) DeleteRecord(context.Context, string, string) error { return dns.ErrUnsupported }
func (p *provider) ListRecords(context.Context, string) ([]dns.Record, error) {
	return nil, dns.ErrUnsupported
}
