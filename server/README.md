Server scaffold for the MVP.

This simple Go server exposes a few endpoints useful for local testing and early prototyping:

- GET  /health                  -- returns "ok"
- POST /move                    -- validates movement (json MoveRequest)
- POST /loot/kill               -- simulates loot roll with pity token
- GET  /pity?playerId=<id>      -- returns current pity tokens for player

Run:

    cd server/go
    go run main.go

Notes:
- This is intentionally minimal. Production requires persistence (database), player authentication,
  authoritative position tracking across ticks, and a real networking layer (UDP/websocket).
