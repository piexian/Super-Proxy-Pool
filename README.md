# Super-Proxy-Pool

Go + SQLite + Mihomo proxy pool management panel.

Current repository status:

- Phase 1 is implemented end-to-end: project structure, SQLite init, login/session auth, four management pages, settings CRUD, manual node CRUD, subscription CRUD, proxy pool CRUD, SSE event channel, Docker build skeleton.
- Phase 2+ entry points already exist: subscription sync parser, node test queue API, pool member management, Mihomo manager placeholder.
- The panel can start now and the main data path is usable. Mihomo hot publish, active probe execution, and richer runtime stats will be filled on top of this base.

## Default behavior

- Panel listen: `0.0.0.0:7890`
- Default password: `admin`
- SQLite path: `/data/app.db` in Docker, `./data/app.db` on Windows local development
- Default subscription sync interval: `3600`
- Default latency URL: `https://www.gstatic.com/generate_204`
- Speed test is off by default

## Local run

```bash
go run ./cmd/app
```

Then open:

- [http://127.0.0.1:7890/login](http://127.0.0.1:7890/login)

## What is already available

### Login

- Password-only login page
- Session cookie auth
- Password stored as bcrypt hash in SQLite
- Password change in system settings page

### Subscription management

- Create, update, delete subscriptions
- Manual sync trigger API
- Subscription detail page
- YAML / Base64 / URI-list parsing logic is in place

### Manual node management

- Import raw nodes from:
  - single URI
  - multiple URI lines
  - Mihomo YAML `proxies` fragments
- Edit / delete / toggle nodes
- Latency and speed test trigger API endpoints are present

### Proxy pools

- Create multiple pools
- HTTP / SOCKS protocol selection
- Port conflict validation against panel port and other pools
- Member selection from manual nodes and subscription nodes
- Publish endpoint and state fields are in place

### System settings

- Panel host and port
- Probe URLs, timeout, concurrency, speed limits
- Controller secret
- Default subscription interval
- Log level
- Restart button

## Docker deployment

### Build

```bash
docker build -t super-proxy-pool .
```

### Run

```bash
docker run -d \
  --name super-proxy-pool \
  --restart unless-stopped \
  -p 7890:7890 \
  -p 18080-18120:18080-18120 \
  -v $PWD/data:/data \
  super-proxy-pool
```

## docker-compose

Bridge mode with pre-opened pool port range:

```bash
docker compose up -d super-proxy-pool
```

Host networking profile:

```bash
docker compose --profile host up -d super-proxy-pool-host
```

## Host network vs bridge

### Option A: host network

Recommended for Linux servers.

Why:

- New proxy pool ports do not need container port re-publishing.
- Operationally simpler when pool ports change frequently.

### Option B: bridge + pre-opened port range

Recommended only when host networking is unavailable.

Why:

- Container stays isolated.
- You must pre-map a fixed port range, for example `18080-18120`.
- Pool listen ports must stay inside the mapped range.

## Common operations

### Add a subscription

1. Open `订阅管理`.
2. Fill in subscription name and URL.
3. Save.
4. Click `立即同步`.

### Add manual nodes

1. Open `节点管理`.
2. Paste one or more URIs or a Mihomo YAML fragment.
3. Save.

### Create a proxy pool

1. Open `代理池设置`.
2. Set pool name, protocol, host and port.
3. Choose members from manual nodes and subscription nodes.
4. Save.

### Change password

1. Open `系统设置`.
2. Enter old password and new password.
3. Submit the password form.

### Restart system

1. Open `系统设置`.
2. Click `重启系统`.
3. In Docker, make sure restart policy is enabled.

## Tests

```bash
go test ./...
```

Covered now:

- subscription parser
- manual node parser
- port conflict validation
- password hash/verify

## Notes on Mihomo

- The Docker image downloads Mihomo from the official MetaCubeX GitHub release assets.
- The Dockerfile exposes `MIHOMO_VERSION` and `MIHOMO_ASSET` as build args, so you can adjust the asset name if upstream naming changes.
- Runtime manager and config file paths are already wired under `/data/runtime`.
