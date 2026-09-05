Unity client skeleton notes

- Open client/Unity folder in Unity Editor (2020.3 LTS or later recommended)
- Create a simple scene with a ground plane and a capsule for the player.
- Attach client/Unity/Assets/Scripts/PlayerController.cs to the player object.
- Ensure the main camera follows the player or is positioned to see the scene.
- This script demonstrates keyboard movement and a tap-to-move placeholder for mobile.
- For production, replace HTTP move validation with websocket/UDP reliable movement and prediction.
