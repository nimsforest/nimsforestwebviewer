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
		[]string{"index.html", "lands.html", "containers.html"},
		nil,
	)

	s.mux.Handle("GET /static/", http.StripPrefix("/static/", nwc.StaticHandler()))
	s.mux.HandleFunc("GET /", s.handleIndex)
	s.mux.HandleFunc("GET /lands", s.handleLands)
	s.mux.HandleFunc("GET /containers", s.handleContainers)
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
