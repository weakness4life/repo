package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// Simple in-memory server demonstrating authoritative movement validation
// and a basic loot + pity token system for MVP prototyping.

type MoveRequest struct {
	PlayerID string  `json:"playerId"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Z        float64 `json:"z"`
	Speed    float64 `json:"speed"`
	TS       int64   `json:"ts"`
}

type MoveResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message,omitempty"`
}

var (
	maxSpeed = 10.0 // units per second (tune for your game)
	mu       sync.Mutex
	pity     = map[string]int{} // playerId -> pity tokens
)

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/move", handleMove)
	http.HandleFunc("/loot/kill", handleLootKill)
	http.HandleFunc("/pity", handlePity)

	addr := ":8080"
	log.Printf("Server starting on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handleMove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req MoveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Authoritative validation: simple speed check
	if req.Speed > maxSpeed*1.5 { // small leniency for latency/prediction
		resp := MoveResponse{Valid: false, Message: "speed too high"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	// In a real server you'd check last known position, timestamps, and do server-side physics.
	resp := MoveResponse{Valid: true}
	json.NewEncoder(w).Encode(resp)
}

// Loot API: POST /loot/kill  { "playerId":"p1", "bossId":"boss_alpha" }
// Simulates an RNG drop with pity: 10% legendary drop chance, but every failed kill grants a pity token.
func handleLootKill(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		PlayerID string `json:"playerId"`
		BossID   string `json:"bossId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Simple RNG
	roll := rand.Float64()
	const legendaryChance = 0.10

	mu.Lock()
	defer mu.Unlock()
	if roll < legendaryChance {
		// drop legendary, reset pity
		pity[body.PlayerID] = 0
		json.NewEncoder(w).Encode(map[string]interface{}{"drop":"legendary", "pity": 0})
		return
	}
	// no drop, give a pity token
	pity[body.PlayerID]++
	// if pity reaches threshold, auto-grant
	const pityThreshold = 10
	if pity[body.PlayerID] >= pityThreshold {
		pity[body.PlayerID] = 0
		json.NewEncoder(w).Encode(map[string]interface{}{"drop":"legendary (pity)", "pity": 0})
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"drop":"none", "pity": pity[body.PlayerID]})
}

// GET /pity?playerId=p1
func handlePity(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Query().Get("playerId")
	if player == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	json.NewEncoder(w).Encode(map[string]interface{}{"playerId": player, "pity": pity[player]})
}
