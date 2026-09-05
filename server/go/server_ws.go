package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocket-based authoritative movement server with server tick loop.
// - Clients connect via WS to /ws
// - Client sends JSON messages of type "move" with position and timestamp
// - Server validates speed against last known pos and broadcasts authoritative positions each tick

const (
	tickRate      = 20                    // server ticks per second
	maxSpeed      = 10.0                  // units per second
	maxSpeedLeeway = 1.5                  // allow some leeway
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// messages
type WSMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type MoveMsg struct {
	PlayerID string  `json:"playerId"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Z        float64 `json:"z"`
	TS       int64   `json:"ts"`
}

type PlayerState struct {
	PlayerID string  `json:"playerId"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Z        float64 `json:"z"`
	TS       int64   `json:"ts"`
}

// connection wrapper
type Client struct {
	conn *websocket.Conn
	mu   sync.Mutex // protects write
	id   string
	state PlayerState
}

var (
	clientsMu sync.Mutex
	clients   = map[string]*Client{}
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	defer c.Close()

	// read initial registration message
	var reg WSMessage
	if err := c.ReadJSON(&reg); err != nil {
		log.Println("read initial registration failed:", err)
		return
	}
	if reg.Type != "register" {
		log.Println("expected register message")
		return
	}
	var regData struct{ PlayerID string `json:"playerId"` }
	if err := json.Unmarshal(reg.Data, &regData); err != nil {
		log.Println("unmarshal register:", err)
		return
	}

	client := &Client{conn: c, id: regData.PlayerID, state: PlayerState{PlayerID: regData.PlayerID}}
	clientsMu.Lock()
	clients[client.id] = client
	clientsMu.Unlock()

	log.Printf("client connected: %s\n", client.id)

	// read loop
	for {
		var msg WSMessage
		if err := c.ReadJSON(&msg); err != nil {
			log.Printf("read json error for %s: %v\n", client.id, err)
			break
		}
		switch msg.Type {
		case "move":
			var m MoveMsg
			if err := json.Unmarshal(msg.Data, &m); err != nil {
				log.Println("bad move msg:", err)
				continue
			}
			// update the client's tentative state (server authoritative will validate on tick)
			clientsMu.Lock()
			if _, ok := clients[client.id]; ok {
				clients[client.id].state.X = m.X
				clients[client.id].state.Y = m.Y
				clients[client.id].state.Z = m.Z
				clients[client.id].state.TS = m.TS
			}
			clientsMu.Unlock()
		default:
			// ignore
		}
	}

	// cleanup
	clientsMu.Lock()
	delete(clients, client.id)
	clientsMu.Unlock()
	log.Printf("client disconnected: %s\n", client.id)
}

func broadcastStates(states []PlayerState) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	if len(clients) == 0 {
		return
	}
	payload := map[string]interface{}{"type": "states", "data": states}
	for _, cl := range clients {
		cl.mu.Lock()
		if err := cl.conn.WriteJSON(payload); err != nil {
			log.Printf("write to %s error: %v\n", cl.id, err)
		}
		cl.mu.Unlock()
	}
}

func tickLoop() {
	ticker := time.NewTicker(time.Second / tickRate)
	defer ticker.Stop()
	// simple map to hold last authoritative positions
	lastPos := map[string]PlayerState{}

	for range ticker.C {
		// build authoritative states
		states := []PlayerState{}
		clientsMu.Lock()
		for id, cl := range clients {
			// validate movement vs lastPos
			prev, ok := lastPos[id]
			if ok {
				dt := float64(cl.state.TS - prev.TS) / 1000.0
				if dt <= 0 {
					dt = 1.0 / tickRate
				}
				dx := cl.state.X - prev.X
				dy := cl.state.Y - prev.Y
				dz := cl.state.Z - prev.Z
				dist := math.Sqrt(dx*dx + dy*dy + dz*dz)
				speed := dist / dt
				if speed > maxSpeed*maxSpeedLeeway {
					// clamp to allowed position
					// compute allowed movement vector
					allowedDist := maxSpeed * dt * maxSpeedLeeway
					ratio := allowedDist / dist
					if ratio < 1 {
						cl.state.X = prev.X + dx*ratio
						cl.state.Y = prev.Y + dy*ratio
						cl.state.Z = prev.Z + dz*ratio
					}
				}
			}
			// update lastPos
			lastPos[id] = cl.state
			states = append(states, cl.state)
		}
		clientsMu.Unlock()

		// broadcast to clients
		broadcastStates(states)
	}
}

func main() {
	go func() {
		log.Println("Starting tick loop")
		tickLoop()
	}()

	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	addr := ":8080"
	log.Printf("WebSocket server starting on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
