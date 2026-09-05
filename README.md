# rep

Mobile-first cross-platform MMORPG — MVP scaffold

This repository contains a minimal scaffold for a mobile-first MMORPG vertical slice: a Unity client skeleton (C# player controller) and a small Go server stub that validates movement and demonstrates a basic loot + pity token mechanic.

Quickstart

1. Server (Go)
   - cd server/go
   - go run main.go
   - Server runs on http://localhost:8080

2. Client (Unity)
   - Open the client/Unity folder in Unity Editor (recommended Unity 2020+)
   - Attach the PlayerController.cs script to a GameObject and set Server URL to http://localhost:8080

What’s included
- server/go/main.go — simple server: movement validation, loot/pity endpoints
- client/Unity/Assets/Scripts/PlayerController.cs — sample movement client (keyboard + basic touch input hooks)
- backlog/MVP-backlog.md — initial tickets and estimates
- LICENSE (MIT)
- .gitignore

Notes
- This is a starting scaffold. It is minimal by design to be useful for prototyping and hiring. Expand with real networking (UDP, WebSocket), authoritative physics, and persistent DB for production.
