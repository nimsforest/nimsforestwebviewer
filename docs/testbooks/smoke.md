# Smoke Test — NimsForest Webviewer

Run after each deploy to verify forest.nimsforest.com is working correctly.

## Prerequisites

- Webviewer deployed and container running on land
- Forest orchestrator running (publishes viewmodel state)
- At least one land node sending heartbeats

## Checks

### 1. Container running

```bash
ssh root@46.225.164.179 "docker ps --format '{{.Names}} {{.Image}} {{.Status}}' | grep webviewer"
```

Expected: `nimsforestwebviewer registry.nimsforest.com/nimsforestwebviewer:<version> Up ...`

### 2. NATS connected

```bash
ssh root@46.225.164.179 "docker logs nimsforestwebviewer 2>&1 | grep Wind"
```

Expected: `[Wind] connected to nats://nimsforest.nimsforest.com:4222`

### 3. Page loads

```bash
curl -s https://forest.nimsforest.com/ | grep '<title>'
```

Expected: `<title>NimsForest Dashboard</title>`

### 4. API state endpoint

```bash
curl -s https://forest.nimsforest.com/api/state | python3 -m json.tool | head -5
```

Expected: valid JSON with `forest_state`, `land_nodes`, `last_updated`

### 5. Zero-count process cards hidden

```bash
curl -s https://forest.nimsforest.com/ | grep -oP 'tracking-wider mt-1">[^<]+'
```

Should only show process types with count > 0 (e.g., "Nims", "Containers"). Should NOT show "Sources", "Trees", "Treehouses", "Songbirds", "Agents" if their counts are 0.

### 6. Containers card visible

When land heartbeats report running containers, a "Containers" card should appear in the Processes section linking to `/containers`.

### 7. Detail pages accessible

Spot-check a few detail pages:

```bash
curl -s https://forest.nimsforest.com/nims | grep '<title>'
curl -s https://forest.nimsforest.com/containers | grep '<title>'
curl -s https://forest.nimsforest.com/infrastructure | grep '<title>'
```

All should return valid HTML.

### 8. Land heartbeats section

The "Land Heartbeats" section should show at least one land node with hostname, CPU, RAM, and container list.
