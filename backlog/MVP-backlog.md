MVP Backlog

Priority: High

- [MVP-1] Scaffold server (Go) - movement validation, health endpoint, loot/pity
  - Estimate: 2d
- [MVP-2] Unity client skeleton - PlayerController script, sample scene
  - Estimate: 2d
- [MVP-3] Integration: client -> server movement validation over HTTP
  - Estimate: 1d
- [MVP-4] Basic loot API + pity tokens implemented and testable
  - Estimate: 1d
- [MVP-5] README, run instructions, license
  - Estimate: 0.5d

Priority: Medium

- [M-1] Replace HTTP with WebSocket/UDP for realtime movement and prediction (server authoritative tick loop)
  - Estimate: 5d
- [M-2] Implement authentication & persistent storage (Postgres)
  - Estimate: 3d
- [M-3] Simple 5-player dungeon instance with shared boss and LFG
  - Estimate: 10d

Priority: Low

- [L-1] Add basic crafting + auction house microservice
  - Estimate: 7d
- [L-2] Build mobile-optimized HUD and touch controls
  - Estimate: 5d

Notes
- Estimates are team-days with a single engineer per ticket. Parallelize tasks where possible.
