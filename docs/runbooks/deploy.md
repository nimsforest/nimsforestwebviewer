# Deploy NimsForest Webviewer

Webviewer is the operator dashboard at `forest.nimsforest.com`, running on land (46.225.164.179), port 8082.

## Prerequisites

- SSH access to land: `ssh root@46.225.164.179`
- Go 1.25+ for cross-compilation
- Repository cloned at `/home/claude-user/nimsforestwebviewer`

## Deploy Steps

### 1. Build

```bash
cd /home/claude-user/nimsforestwebviewer
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags="-s -w -X main.version=<VERSION>" \
  -o bin/nimsforestwebviewer-linux-amd64 ./cmd/nimsforestwebviewer
```

### 2. Upload binary

```bash
scp bin/nimsforestwebviewer-linux-amd64 root@46.225.164.179:/opt/nimsforestwebviewer/nimsforestwebviewer
```

### 3. Build image and restart

```bash
ssh root@46.225.164.179 "
  cd /opt/nimsforestwebviewer &&
  docker build -t registry.nimsforest.com/nimsforestwebviewer:<VERSION> \
               -t registry.nimsforest.com/nimsforestwebviewer:latest . &&
  docker stop nimsforestwebviewer &&
  docker rm nimsforestwebviewer &&
  docker run -d --name nimsforestwebviewer --network host --restart unless-stopped \
    registry.nimsforest.com/nimsforestwebviewer:<VERSION>
"
```

### 4. Verify

```bash
ssh root@46.225.164.179 "docker logs nimsforestwebviewer 2>&1 | tail -5"
# Should show: connected to NATS, serving on :8082
curl -s https://forest.nimsforest.com/ | grep '<title>'
# Should show: NimsForest Dashboard
```

## Rollback

```bash
ssh root@46.225.164.179 "
  docker stop nimsforestwebviewer && docker rm nimsforestwebviewer &&
  docker run -d --name nimsforestwebviewer --network host --restart unless-stopped \
    registry.nimsforest.com/nimsforestwebviewer:<PREVIOUS_VERSION>
"
```

## Configuration

No config file — webviewer uses hardcoded defaults:
- Listen: `:8082`
- NATS: `nats://nimsforest.nimsforest.com:4222`

## Architecture

- Subscribes to `forest.viewmodel.state` and `forest.land.heartbeat` via NATS
- Multi-page dashboard with detail views for each process type
- HTMX-powered with fragment rendering for navigation
- No database, no persistent state — all data from NATS subscriptions
