package web

import (
	"sync"
	"time"
)

// ForestState mirrors the PublishedState from nimsforest2 viewmodel.
type ForestState struct {
	Timestamp  time.Time         `json:"timestamp"`
	Summary    ForestSummary     `json:"summary"`
	Lands      []LandVM          `json:"lands"`
	Trees      []ProcessVM       `json:"trees"`
	Treehouses []TreehouseVM     `json:"treehouses"`
	Nims       []NimVM           `json:"nims"`
}

type ForestSummary struct {
	LandCount         int     `json:"LandCount"`
	ManalandCount     int     `json:"ManalandCount"`
	TotalRAM          uint64  `json:"TotalRAM"`
	TotalCPUCores     int     `json:"TotalCPUCores"`
	TotalManaVram     uint64  `json:"TotalManaVram"`
	TreeCount         int     `json:"TreeCount"`
	TreehouseCount    int     `json:"TreehouseCount"`
	NimCount          int     `json:"NimCount"`
	TotalRAMAllocated uint64  `json:"TotalRAMAllocated"`
	Occupancy         float64 `json:"Occupancy"`
}

type LandVM struct {
	ID         string      `json:"id"`
	Hostname   string      `json:"hostname"`
	RAMTotal   uint64      `json:"ram_total"`
	CPUCores   int         `json:"cpu_cores"`
	GPUVram    uint64      `json:"gpu_vram"`
	Trees      []ProcessVM `json:"trees"`
	Treehouses []TreehouseVM `json:"treehouses"`
	Nims       []NimVM     `json:"nims"`
	JoinedAt   time.Time   `json:"joined_at"`
	LastSeen   time.Time   `json:"last_seen"`
}

func (l *LandVM) HasMana() bool {
	return l.GPUVram > 0
}

func (l *LandVM) ProcessCount() int {
	return len(l.Trees) + len(l.Treehouses) + len(l.Nims)
}

func (l *LandVM) RAMAllocated() uint64 {
	var total uint64
	for _, t := range l.Trees {
		total += t.RAMAllocated
	}
	for _, t := range l.Treehouses {
		total += t.RAMAllocated
	}
	for _, n := range l.Nims {
		total += n.RAMAllocated
	}
	return total
}

func (l *LandVM) Occupancy() float64 {
	if l.RAMTotal == 0 {
		return 0
	}
	return float64(l.RAMAllocated()) / float64(l.RAMTotal) * 100
}

type ProcessVM struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	RAMAllocated uint64    `json:"ram_allocated"`
	LandID       string    `json:"land_id"`
	StartedAt    time.Time `json:"started_at"`
}

type TreehouseVM struct {
	ProcessVM
	ScriptPath string `json:"script_path"`
}

type NimVM struct {
	ProcessVM
	AIEnabled bool   `json:"ai_enabled"`
	Model     string `json:"model"`
}

// LandHeartbeat matches the payload from land's heartbeat.
type LandHeartbeat struct {
	Hostname   string               `json:"hostname"`
	CPUCores   int                  `json:"cpu_cores"`
	RAMTotalMB uint64               `json:"ram_total_mb"`
	RAMUsedMB  uint64               `json:"ram_used_mb"`
	Mana       ManaInfo             `json:"mana"`
	Containers []HeartbeatContainer `json:"containers"`
}

type ManaInfo struct {
	Available bool   `json:"available"`
	Vendor    string `json:"vendor,omitempty"`
	Model     string `json:"model,omitempty"`
	VRAM      uint64 `json:"vram,omitempty"`
}

type HeartbeatContainer struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Image  string `json:"image,omitempty"`
	Uptime string `json:"uptime,omitempty"`
}

// LandNode represents a land node tracked via heartbeats.
type LandNode struct {
	Source     string
	Heartbeat *LandHeartbeat
	LastSeen  time.Time
}

// StateCache holds the latest forest state and land heartbeats.
type StateCache struct {
	mu           sync.RWMutex
	forestState  *ForestState
	landNodes    map[string]*LandNode
	lastReceived time.Time
}

// NewStateCache creates a new empty state cache.
func NewStateCache() *StateCache {
	return &StateCache{
		landNodes: make(map[string]*LandNode),
	}
}

// Update replaces the cached forest state.
func (s *StateCache) Update(state *ForestState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.forestState = state
	s.lastReceived = time.Now()
}

// UpdateLandHeartbeat updates the cached heartbeat for a land node.
func (s *StateCache) UpdateLandHeartbeat(source string, hb *LandHeartbeat) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.landNodes[source] = &LandNode{
		Source:    source,
		Heartbeat: hb,
		LastSeen:  time.Now(),
	}
}

// GetState returns the latest forest state (may be nil).
func (s *StateCache) GetState() *ForestState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.forestState
}

// GetLandNodes returns all tracked land nodes.
func (s *StateCache) GetLandNodes() []*LandNode {
	s.mu.RLock()
	defer s.mu.RUnlock()
	nodes := make([]*LandNode, 0, len(s.landNodes))
	for _, n := range s.landNodes {
		nodes = append(nodes, n)
	}
	return nodes
}

// LastReceived returns the time of the last state update.
func (s *StateCache) LastReceived() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastReceived
}
