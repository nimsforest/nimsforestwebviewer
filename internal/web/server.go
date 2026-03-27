package web

import (
	"encoding/json"
	"net/http"
	"time"

	nwc "github.com/nimsforest/nimsforestwebcomponents"
)

// Server is the HTTP handler for the dashboard.
type Server struct {
	mux      *http.ServeMux
	state    *StateCache
	renderer *nwc.Renderer
	version  string
}

// DashboardData is passed to the overview template.
type DashboardData struct {
	State                 *ForestState
	LandNodes             []*LandNode
	RunningContainerCount int
	LastUpdated           time.Time
	Version               string
	Connected             bool
}

// NewServer creates a new dashboard HTTP server.
func NewServer(state *StateCache, version string) *Server {
	s := &Server{
		mux:     http.NewServeMux(),
		state:   state,
		version: version,
	}

	s.renderer = nwc.NewRenderer(
		templateFS, "templates",
		[]string{
			"index.html", "lands.html", "containers.html",
			"trees.html", "tree-detail.html",
			"treehouses.html", "treehouse-detail.html",
			"nims.html", "nim-detail.html",
			"songbirds.html", "songbird-detail.html",
			"agents.html",
			"sources.html", "source-detail.html",
			"infrastructure.html", "infrastructure-detail.html",
		},
		nil,
		nwc.AppConfig{
			Name:  "Forest",
			Emoji: "🌲",
			NavItems: []nwc.NavItem{
				{Label: "Overview", Href: "/"},
				{Label: "Sources", Href: "/sources"},
				{Label: "Trees", Href: "/trees"},
				{Label: "Treehouses", Href: "/treehouses"},
				{Label: "Nims", Href: "/nims"},
				{Label: "Songbirds", Href: "/songbirds"},
				{Label: "Agents", Href: "/agents"},
				{Label: "Lands", Href: "/lands"},
				{Label: "Containers", Href: "/containers"},
				{Label: "Infrastructure", Href: "/infrastructure"},
			},
			Footer: "NimsForest Dashboard · Live updates via HTMX",
		},
	)

	s.mux.Handle("GET /static/", http.StripPrefix("/static/", nwc.StaticHandler()))
	s.mux.HandleFunc("GET /", s.handleIndex)
	s.mux.HandleFunc("GET /lands", s.handleLands)
	s.mux.HandleFunc("GET /containers", s.handleContainers)
	s.mux.HandleFunc("GET /trees", s.handleTrees)
	s.mux.HandleFunc("GET /trees/{name}", s.handleTreeDetail)
	s.mux.HandleFunc("GET /treehouses", s.handleTreehouses)
	s.mux.HandleFunc("GET /treehouses/{name}", s.handleTreehouseDetail)
	s.mux.HandleFunc("GET /nims", s.handleNims)
	s.mux.HandleFunc("GET /nims/{name}", s.handleNimDetail)
	s.mux.HandleFunc("GET /songbirds", s.handleSongbirds)
	s.mux.HandleFunc("GET /songbirds/{name}", s.handleSongbirdDetail)
	s.mux.HandleFunc("GET /agents", s.handleAgents)
	s.mux.HandleFunc("GET /sources", s.handleSources)
	s.mux.HandleFunc("GET /sources/{name}", s.handleSourceDetail)
	s.mux.HandleFunc("GET /infrastructure", s.handleInfrastructure)
	s.mux.HandleFunc("GET /infrastructure/{name}", s.handleInfrastructureDetail)
	s.mux.HandleFunc("GET /api/state", s.handleAPIState)

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// render calls RenderFragment for HTMX requests, Render for full page loads.
func (s *Server) render(w http.ResponseWriter, r *http.Request, page, title string, data any) {
	if r.Header.Get("HX-Request") == "true" {
		s.renderer.RenderFragment(w, page, title, data)
		return
	}
	s.renderer.Render(w, page, title, data)
}

func (s *Server) dashboardData() *DashboardData {
	lands := s.state.GetLandNodes()
	var containerCount int
	for _, ln := range lands {
		if ln.Heartbeat == nil {
			continue
		}
		for _, c := range ln.Heartbeat.Containers {
			if c.Status == "running" {
				containerCount++
			}
		}
	}
	return &DashboardData{
		State:                 s.state.GetState(),
		LandNodes:             lands,
		RunningContainerCount: containerCount,
		LastUpdated:           s.state.LastReceived(),
		Version:               s.version,
		Connected:             !s.state.LastReceived().IsZero(),
	}
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	s.render(w, r, "index.html", "NimsForest Dashboard", s.dashboardData())
}

func (s *Server) handleLands(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "lands.html", "Lands - NimsForest", s.dashboardData())
}

func (s *Server) handleContainers(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "containers.html", "Containers - NimsForest", s.dashboardData())
}

func (s *Server) handleTrees(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "trees.html", "Trees - NimsForest", s.dashboardData())
}

func (s *Server) handleTreeDetail(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	data := s.dashboardData()
	if data.State == nil {
		http.NotFound(w, r)
		return
	}
	tree, land := data.State.FindTree(name)
	if tree == nil {
		http.NotFound(w, r)
		return
	}
	s.render(w, r, "tree-detail.html", tree.Name+" - Trees - NimsForest", &DetailData{Dashboard: data, Tree: tree, Land: land})
}

func (s *Server) handleTreehouses(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "treehouses.html", "Treehouses - NimsForest", s.dashboardData())
}

func (s *Server) handleTreehouseDetail(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	data := s.dashboardData()
	if data.State == nil {
		http.NotFound(w, r)
		return
	}
	th, land := data.State.FindTreehouse(name)
	if th == nil {
		http.NotFound(w, r)
		return
	}
	s.render(w, r, "treehouse-detail.html", th.Name+" - Treehouses - NimsForest", &DetailData{Dashboard: data, Treehouse: th, Land: land})
}

func (s *Server) handleNims(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "nims.html", "Nims - NimsForest", s.dashboardData())
}

func (s *Server) handleNimDetail(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	data := s.dashboardData()
	if data.State == nil {
		http.NotFound(w, r)
		return
	}
	nim, land := data.State.FindNim(name)
	if nim == nil {
		http.NotFound(w, r)
		return
	}
	s.render(w, r, "nim-detail.html", nim.Name+" - Nims - NimsForest", &DetailData{Dashboard: data, Nim: nim, Land: land})
}

func (s *Server) handleSongbirds(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "songbirds.html", "Songbirds - NimsForest", s.dashboardData())
}

func (s *Server) handleAgents(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "agents.html", "Agents - NimsForest", s.dashboardData())
}

func (s *Server) handleSongbirdDetail(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	data := s.dashboardData()
	if data.State == nil {
		http.NotFound(w, r)
		return
	}
	sb, land := data.State.FindSongbird(name)
	if sb == nil {
		http.NotFound(w, r)
		return
	}
	s.render(w, r, "songbird-detail.html", sb.Name+" - Songbirds - NimsForest", &DetailData{Dashboard: data, Songbird: sb, Land: land})
}

func (s *Server) handleSources(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "sources.html", "Sources - NimsForest", s.dashboardData())
}

func (s *Server) handleSourceDetail(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	data := s.dashboardData()
	if data.State == nil {
		http.NotFound(w, r)
		return
	}
	src, land := data.State.FindSource(name)
	if src == nil {
		http.NotFound(w, r)
		return
	}
	s.render(w, r, "source-detail.html", src.Name+" - Sources - NimsForest", &DetailData{Dashboard: data, Source: src, Land: land})
}

// DetailData is passed to detail page templates.
type DetailData struct {
	Dashboard *DashboardData
	Tree      *TreeVM
	Treehouse *TreehouseVM
	Nim       *NimVM
	Songbird  *SongbirdVM
	Source    *SourceVM
	Land      *LandVM
}

// InfraDetailData is passed to infrastructure detail templates.
type InfraDetailData struct {
	Dashboard *DashboardData
	Infra     *InfraDetailVM
}

// InfraDetailVM holds display data for a single infrastructure component.
type InfraDetailVM struct {
	Name         string
	Slug         string
	Icon         string
	TypeLabel    string
	BadgeClass   string
	Description  string
	Active       bool
	HasMetrics   bool
	ShowMessages bool
	ShowBytes    bool
	ShowConsumers bool
	ShowKeys     bool
	Messages     uint64
	Bytes        uint64
	Consumers    int
	Keys         uint64
	Subjects     []string
}

func (s *Server) handleInfrastructure(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "infrastructure.html", "Infrastructure - NimsForest", s.dashboardData())
}

func (s *Server) handleInfrastructureDetail(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	data := s.dashboardData()

	var vm *InfraDetailVM

	switch name {
	case "river":
		vm = &InfraDetailVM{
			Name:        "River",
			Slug:        "river",
			Icon:        "\U0001F30A",
			TypeLabel:   "JetStream WorkQueue",
			BadgeClass:  "bg-cyan-400/10 text-cyan-400",
			Description: "JetStream stream for external unstructured data. WorkQueue retention policy ensures each message is processed exactly once. 24-hour max age, subject pattern river.>",
			Active:      true,
			HasMetrics:  true,
		}
		if data.State != nil {
			if st := data.State.FindStream("RIVER"); st != nil {
				vm.ShowMessages = true
				vm.ShowBytes = true
				vm.ShowConsumers = true
				vm.Messages = st.Messages
				vm.Bytes = st.Bytes
				vm.Consumers = st.Consumers
				vm.Subjects = st.Subjects
			}
		}
	case "wind":
		vm = &InfraDetailVM{
			Name:        "Wind",
			Slug:        "wind",
			Icon:        "\U0001F4A8",
			TypeLabel:   "Core Pub/Sub",
			BadgeClass:  "bg-sky-400/10 text-sky-400",
			Description: "NATS core pub/sub transport for carrying leaves between processes. Ephemeral — messages are not persisted, only delivered to active subscribers. Used by Drop() and Catch() in the nim package.",
			Active:      data.State != nil && data.State.Infrastructure.WindActive,
		}
	case "soil":
		vm = &InfraDetailVM{
			Name:        "Soil",
			Slug:        "soil",
			Icon:        "\U0001F4E6",
			TypeLabel:   "KV Bucket",
			BadgeClass:  "bg-amber-400/10 text-amber-400",
			Description: "JetStream KV bucket for shared cluster state. Stores configuration, process definitions, and runtime state as key-value pairs accessible from any node.",
			Active:      true,
			HasMetrics:  true,
		}
		if data.State != nil {
			if kv := data.State.FindKVStore("SOIL"); kv != nil {
				vm.ShowKeys = true
				vm.ShowBytes = true
				vm.Keys = kv.Keys
				vm.Bytes = kv.Bytes
			}
		}
	case "humus":
		vm = &InfraDetailVM{
			Name:        "Humus",
			Slug:        "humus",
			Icon:        "\U0001F4DC",
			TypeLabel:   "JetStream Limits",
			BadgeClass:  "bg-emerald-400/10 text-emerald-400",
			Description: "JetStream stream for audit logging. Limits retention policy with 7-day max age. Records system events, process lifecycle changes, and operational actions on subject pattern humus.>",
			Active:      true,
			HasMetrics:  true,
		}
		if data.State != nil {
			if st := data.State.FindStream("HUMUS"); st != nil {
				vm.ShowMessages = true
				vm.ShowBytes = true
				vm.ShowConsumers = true
				vm.Messages = st.Messages
				vm.Bytes = st.Bytes
				vm.Consumers = st.Consumers
				vm.Subjects = st.Subjects
			}
		}
	default:
		http.NotFound(w, r)
		return
	}

	s.render(w, r, "infrastructure-detail.html", vm.Name+" - Infrastructure - NimsForest", &InfraDetailData{Dashboard: data, Infra: vm})
}

func (s *Server) handleAPIState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := struct {
		ForestState *ForestState `json:"forest_state"`
		LandNodes   []*LandNode `json:"land_nodes"`
		LastUpdated time.Time   `json:"last_updated"`
	}{
		ForestState: s.state.GetState(),
		LandNodes:   s.state.GetLandNodes(),
		LastUpdated: s.state.LastReceived(),
	}
	json.NewEncoder(w).Encode(data)
}
