# dns-caddy — docs

**Caddy.** Reverse-proxy routes with automatic HTTPS via the Caddy admin API.

## Install

```bash
togo install togo-framework/dns-caddy
```

Registers on the [`dns`](https://github.com/togo-framework/dns) base; select it with **dns.provider in togo.yaml (or DNS_DRIVER)**, then use **`togo proxy`**.

## Interface

`Provider` — `UpsertRecord`/`DeleteRecord`/`ListRecords`, `UpsertProxyHost`/`DeleteProxyHost`, `UpsertRoute`/`DeleteRoute`.

## Configuration

| Env var | Description |
|---|---|
| `CADDY_ADMIN` | Caddy admin API URL (default `http://localhost:2019`). |
| `CADDY_SERVER` | Caddy server name to attach routes to (default `srv0`). |

## Usage & notes

Adds/removes `reverse_proxy` routes (stable `@id`, replace-on-upsert) on the running Caddy server. Auto-TLS handled by Caddy.

## Example

```bash
togo proxy:host:add app.example.com http://localhost:3000 --provider caddy --dry-run
```

## Links

- [Caddy admin API](https://caddyserver.com/docs/api)
- [Marketplace](https://to-go.dev/marketplace)
- [Source](https://github.com/togo-framework/dns-caddy)
