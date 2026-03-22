# Testbook: Hub Forest Dashboard (forest.hub.nimsforest.com)

Validates the forest webviewer on the hub shows correct infrastructure state
across all pages using headless browser checks.

## Prerequisites

- Hub (land-shared-one) running with nimsforestwebviewer and nimsforestviewmodeltreehouse
- `npx @browserbasehq/browse-cli` available on neoremote
- Dashboard accessible at https://forest.hub.nimsforest.com

## 1. Overview page

```bash
npx @browserbasehq/browse-cli "Go to https://forest.hub.nimsforest.com/ — \
  verify the page title says 'Forest Overview', \
  there is a 'Processes' section, \
  there is a 'Resources' section, \
  and the navigation bar has links for Sources, Trees, Treehouses, Nims, Songbirds, Agents, Lands, Containers, Infrastructure"
```

- [ ] Page loads with "Forest Overview" heading
- [ ] Processes section visible
- [ ] Resources section visible
- [ ] All navigation links present

## 2. Lands page

```bash
npx @browserbasehq/browse-cli "Go to https://forest.hub.nimsforest.com/lands — \
  verify the page title says 'Lands', \
  and check if there is at least one land node listed (look for a card or row with a hostname)"
```

- [ ] Page loads with "Lands" heading
- [ ] At least one land node visible (land-shared-one)

## 3. Containers page

```bash
npx @browserbasehq/browse-cli "Go to https://forest.hub.nimsforest.com/containers — \
  verify the page shows a list of containers, \
  and check that these container names are visible: nimsforest, mycelium, landregistry, \
  hetznertreehouse, pantheon, nimsforestecommerce, nimsforestwebviewer"
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
npx @browserbasehq/browse-cli "Go to https://forest.hub.nimsforest.com/trees — \
  verify the page loads with a 'Trees' heading. \
  Trees parse River data into Leaves. There may be trees listed or the page may show zero trees \
  if the hub forest has no trees configured. Report what you see."
```

- [ ] Page loads

## 5. Treehouses page

```bash
npx @browserbasehq/browse-cli "Go to https://forest.hub.nimsforest.com/treehouses — \
  verify the page loads with a 'Treehouses' heading. Report any treehouses listed."
```

- [ ] Page loads

## 6. Nims page

```bash
npx @browserbasehq/browse-cli "Go to https://forest.hub.nimsforest.com/nims — \
  verify the page loads with a 'Nims' heading. Report any nims listed."
```

- [ ] Page loads

## 7. Songbirds page

```bash
npx @browserbasehq/browse-cli "Go to https://forest.hub.nimsforest.com/songbirds — \
  verify the page loads with a 'Songbirds' heading. Report any songbirds listed."
```

- [ ] Page loads

## 8. Agents page

```bash
npx @browserbasehq/browse-cli "Go to https://forest.hub.nimsforest.com/agents — \
  verify the page loads with an 'Agents' heading. Report any agents listed."
```

- [ ] Page loads

## 9. Sources page

```bash
npx @browserbasehq/browse-cli "Go to https://forest.hub.nimsforest.com/sources — \
  verify the page loads with a 'Sources' heading. \
  The hub has webhook sources configured (stripe, github, releases, etc). \
  Check if any sources are listed."
```

- [ ] Page loads
- [ ] Sources visible (stripe, github, releases, etc.)

## 10. Infrastructure page

```bash
npx @browserbasehq/browse-cli "Go to https://forest.hub.nimsforest.com/infrastructure — \
  verify the page loads with an 'Infrastructure' heading. \
  This page shows the NimsForest infrastructure components (Wind, River, Soil, Humus, Taproot). \
  Check that at least Wind and Soil are listed."
```

- [ ] Page loads
- [ ] Infrastructure components visible (Wind, Soil, etc.)

## Pass Criteria

- [ ] All 10 pages load without errors
- [ ] Overview shows process and resource counts
- [ ] Containers page shows all hub services
- [ ] Navigation works across all pages
- [ ] TLS works (https://forest.hub.nimsforest.com)
