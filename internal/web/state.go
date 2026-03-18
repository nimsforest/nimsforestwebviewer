package web

import (
	"sync"
	"time"

	vm "github.com/nimsforest/nimsforest2/pkg/viewmodel"
)

// Type aliases — map local names to the canonical wire types from pkg/viewmodel.
type ForestState = vm.PublishedState
type ForestSummary = vm.Summary
type LandVM = vm.LandViewModel
type ProcessVM = vm.Process
type TreeVM = vm.TreeViewModel
type TreehouseVM = vm.TreehouseViewModel
type NimVM = vm.NimViewModel
type SongbirdVM = vm.SongbirdViewModel
type SourceVM = vm.SourceViewModel
type AgentVM = vm.AgentViewModel
type InfrastructureState = vm.Infrastructure
type StreamStatusVM = vm.StreamStatus
type KVStatusVM = vm.KVStatus

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
	Source    string
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
