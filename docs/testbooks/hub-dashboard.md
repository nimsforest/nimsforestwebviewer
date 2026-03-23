# Testbook: Hub Forest Dashboard (forest.hub.nimsforest.com)

Validates the forest webviewer on the hub shows correct infrastructure state
across all pages using headless browser checks via PinchTab.

## Prerequisites

- Hub (land-shared-one) running with nimsforestwebviewer and nimsforestviewmodeltreehouse
- PinchTab running on neoremote (`pinchtab health` returns ok)
- Dashboard accessible at https://forest.hub.nimsforest.com

## Setup

```bash
pinchtab health
pinchtab nav https://forest.hub.nimsforest.com/
```

## 1. Overview page

```bash
pinchtab nav https://forest.hub.nimsforest.com/
pinchtab snap -c
```

- [ ] Page loads with "Forest Overview" heading
- [ ] Processes section visible
- [ ] Resources section visible
- [ ] All navigation links present

## 2. Lands page

```bash
pinchtab nav https://forest.hub.nimsforest.com/lands
pinchtab snap -c
```

- [ ] Page loads with "Lands" heading
- [ ] At least one land node visible (land-shared-one)

## 3. Containers page

```bash
pinchtab nav https://forest.hub.nimsforest.com/containers
pinchtab snap -c
```

- [ ] Page loads with container list
- [ ] nimsforest visible
- [ ] mycelium visible
- [ ] landregistry visible
- [ ] hetznertreehouse visible
- [ ] pantheon visible
- [ ] nimsforestecommerce visible
- [ ] nimsforestwebviewer visible

## 4. Trees page

```bash
pinchtab nav https://forest.hub.nimsforest.com/trees
pinchtab snap -c
```

- [ ] Page loads

## 5. Treehouses page

```bash
pinchtab nav https://forest.hub.nimsforest.com/treehouses
pinchtab snap -c
```

- [ ] Page loads

## 6. Nims page

```bash
pinchtab nav https://forest.hub.nimsforest.com/nims
pinchtab snap -c
```

- [ ] Page loads

## 7. Songbirds page

```bash
pinchtab nav https://forest.hub.nimsforest.com/songbirds
pinchtab snap -c
```

- [ ] Page loads

## 8. Agents page

```bash
pinchtab nav https://forest.hub.nimsforest.com/agents
pinchtab snap -c
```

- [ ] Page loads

## 9. Sources page

```bash
pinchtab nav https://forest.hub.nimsforest.com/sources
pinchtab snap -c
```

- [ ] Page loads
- [ ] Sources visible (stripe, github, releases, etc.)

## 10. Infrastructure page

```bash
pinchtab nav https://forest.hub.nimsforest.com/infrastructure
pinchtab snap -c
```

- [ ] Page loads
- [ ] Infrastructure components visible (Wind, Soil, etc.)

## Pass Criteria

- [ ] All 10 pages load without errors
- [ ] Overview shows process and resource counts
- [ ] Containers page shows all hub services
- [ ] Navigation works across all pages
- [ ] TLS works (https://forest.hub.nimsforest.com)
