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
	State       *ForestState
	LandNodes   []*LandNode
	LastUpdated time.Time
	Version     string
	Connected   bool
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
			"sources.html", "source-detail.html",
		},
		nil,
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
	s.mux.HandleFunc("GET /sources", s.handleSources)
	s.mux.HandleFunc("GET /sources/{name}", s.handleSourceDetail)
	s.mux.HandleFunc("GET /api/state", s.handleAPIState)

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) dashboardData() *DashboardData {
	return &DashboardData{
		State:       s.state.GetState(),
		LandNodes:   s.state.GetLandNodes(),
		LastUpdated: s.state.LastReceived(),
		Version:     s.version,
		Connected:   !s.state.LastReceived().IsZero(),
	}
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	s.renderer.Render(w, "index.html", "NimsForest Dashboard", s.dashboardData())
}

func (s *Server) handleLands(w http.ResponseWriter, r *http.Request) {
	s.renderer.Render(w, "lands.html", "Lands - NimsForest", s.dashboardData())
}

func (s *Server) handleContainers(w http.ResponseWriter, r *http.Request) {
	s.renderer.Render(w, "containers.html", "Containers - NimsForest", s.dashboardData())
}

func (s *Server) handleTrees(w http.ResponseWriter, r *http.Request) {
	s.renderer.Render(w, "trees.html", "Trees - NimsForest", s.dashboardData())
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
	s.renderer.Render(w, "tree-detail.html", tree.Name+" - Trees - NimsForest", &DetailData{Dashboard: data, Process: tree, Land: land})
}

func (s *Server) handleTreehouses(w http.ResponseWriter, r *http.Request) {
	s.renderer.Render(w, "treehouses.html", "Treehouses - NimsForest", s.dashboardData())
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
	s.renderer.Render(w, "treehouse-detail.html", th.Name+" - Treehouses - NimsForest", &DetailData{Dashboard: data, Treehouse: th, Land: land})
}

func (s *Server) handleNims(w http.ResponseWriter, r *http.Request) {
	s.renderer.Render(w, "nims.html", "Nims - NimsForest", s.dashboardData())
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
	s.renderer.Render(w, "nim-detail.html", nim.Name+" - Nims - NimsForest", &DetailData{Dashboard: data, Nim: nim, Land: land})
}

func (s *Server) handleSongbirds(w http.ResponseWriter, r *http.Request) {
	s.renderer.Render(w, "songbirds.html", "Songbirds - NimsForest", s.dashboardData())
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
	s.renderer.Render(w, "songbird-detail.html", sb.Name+" - Songbirds - NimsForest", &DetailData{Dashboard: data, Songbird: sb, Land: land})
}

func (s *Server) handleSources(w http.ResponseWriter, r *http.Request) {
	s.renderer.Render(w, "sources.html", "Sources - NimsForest", s.dashboardData())
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
	s.renderer.Render(w, "source-detail.html", src.Name+" - Sources - NimsForest", &DetailData{Dashboard: data, Source: src, Land: land})
}

// DetailData is passed to detail page templates.
type DetailData struct {
	Dashboard *DashboardData
	Process   *ProcessVM
	Treehouse *TreehouseVM
	Nim       *NimVM
	Songbird  *SongbirdVM
	Source    *SourceVM
	Land      *LandVM
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
