package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/nimsforest/nimsforest2/pkg/nim"
	"github.com/nimsforest/nimsforestwebviewer/internal/web"
)

var version = "dev"

func main() {
	natsURL := flag.String("nats", "nats://nimsforest.nimsforest.com:4222", "NATS server URL")
	listen := flag.String("listen", ":8082", "HTTP listen address")
	flag.Parse()

	// Allow env override
	if env := os.Getenv("NATS_URL"); env != "" {
		*natsURL = env
	}

	log.Printf("nimsforestwebviewer %s starting", version)

	// Connect to Wind
	wind, err := nim.Connect(*natsURL, "nimsforestwebviewer")
	if err != nil {
		log.Fatalf("failed to connect to NATS: %v", err)
	}
	defer wind.Close()

	// Create state cache
	state := web.NewStateCache()

	// Subscribe to viewmodel state
	_, err = wind.Catch("forest.viewmodel.state", func(leaf nim.Leaf) {
		var published web.ForestState
		if err := json.Unmarshal(leaf.Data, &published); err != nil {
			log.Printf("failed to parse viewmodel state: %v", err)
			return
		}
		state.Update(&published)
	})
	if err != nil {
		log.Fatalf("failed to subscribe to forest.viewmodel.state: %v", err)
	}

	// Also subscribe to land heartbeats to track land nodes directly
	_, err = wind.Catch("forest.land.heartbeat", func(leaf nim.Leaf) {
		var hb web.LandHeartbeat
		if err := json.Unmarshal(leaf.Data, &hb); err != nil {
			log.Printf("failed to parse heartbeat: %v", err)
			return
		}
		state.UpdateLandHeartbeat(leaf.Source, &hb)
	})
	if err != nil {
		log.Fatalf("failed to subscribe to heartbeats: %v", err)
	}

	// Serve HTTP
	srv := web.NewServer(state, version)
	log.Printf("serving on %s", *listen)
	if err := http.ListenAndServe(*listen, srv); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
