<!-- togo-header -->
<div align="center">
  <img src=".github/assets/togo-mark.svg" alt="togo" height="64" />
  <h1>togo-framework/dns-caddy</h1>
  <p>Caddy reverse-proxy driver for togo dns.</p>
  <p>
    <a href="https://to-go.dev/marketplace"><img src="https://img.shields.io/badge/marketplace-to--go.dev-1FC7DC" alt="marketplace" /></a>
    <a href="https://pkg.go.dev/github.com/togo-framework/dns-caddy"><img src="https://pkg.go.dev/badge/github.com/togo-framework/dns-caddy.svg" alt="pkg.go.dev" /></a>
    <img src="https://img.shields.io/badge/license-MIT-blue" alt="MIT" />
  </p>
  <p><strong>Part of the <a href="https://to-go.dev">togo</a> framework.</strong></p>
</div>

## Install

```bash
togo install togo-framework/dns-caddy
```
<!-- /togo-header -->

[Caddy](https://caddyserver.com/) reverse-proxy driver for togo's
[`dns`](https://github.com/togo-framework/dns) subsystem, driven through Caddy's
admin API. Both proxy hosts and gateway routes become Caddy HTTP routes
(host/path matcher → `reverse_proxy`). Caddy provisions TLS automatically.

## Config

| Env | Meaning |
|-----|---------|
| `DNS_DRIVER` | set to `caddy` |
| `CADDY_ADMIN` | admin endpoint (default `http://localhost:2019`) |
| `CADDY_SERVER` | server key in the Caddy config (default `srv0`) |

```go
svc, _ := dns.FromKernel(k)
svc.UpsertProxyHost(ctx, dns.ProxyHost{Domain: "app.example.com", Upstream: "http://127.0.0.1:8080"})
svc.UpsertRoute(ctx, dns.Route{Domain: "api.example.com", Path: "/v1", Upstream: "http://gateway:8000"})
```

Routes are tagged with a stable `@id` so an upsert replaces in place. DNS records
return `dns.ErrUnsupported`.

<!-- togo-sponsors -->
---
<div align="center">
  <h3>Premium sponsors</h3>
  <p><a href="https://id8media.com"><strong>ID8 Media</strong></a> &nbsp;·&nbsp; <a href="https://one-studio.co"><strong>One Studio</strong></a></p>
  <p><sub>Support togo — <a href="https://github.com/sponsors/fadymondy">become a sponsor</a>.</sub></p>
</div>
<!-- /togo-sponsors -->
