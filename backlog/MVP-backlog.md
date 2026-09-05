MVP Backlog (updated)

Priority: High

- [MVP-1] Scaffold server (Go) - movement validation, health endpoint, loot/pity
  - Estimate: 2d (done)
- [MVP-2] Unity client skeleton - PlayerController script, sample scene
  - Estimate: 2d (done)
- [MVP-3] Integration: client -> server movement validation over HTTP
  - Estimate: 1d (done)
- [MVP-4] Basic loot API + pity tokens implemented and testable
  - Estimate: 1d (done)
- [MVP-5] README, run instructions, license
  - Estimate: 0.5d (done)
- [M-1] Replace HTTP with WebSocket/UDP for realtime movement and prediction (server authoritative tick loop)
  - Estimate: 5d (in progress)

Priority: Medium

- [M-2] Implement authentication & persistent storage (Postgres)
  - Estimate: 3d
- [M-3] Simple 5-player dungeon instance with shared boss and LFG
  - Estimate: 10d

Notes
- WebSocket server implemented in server/go/server_ws.go. Use server/ws_client.js to simulate a client.
