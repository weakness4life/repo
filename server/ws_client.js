// Node.js test WebSocket client
// Usage:
// 1) npm install ws
// 2) node ws_client.js

const WebSocket = require('ws');
const ws = new WebSocket('ws://localhost:8080/ws');

ws.on('open', function open() {
  console.log('connected');
  // register
  ws.send(JSON.stringify({ type: 'register', data: { playerId: 'player1' } }));

  // send a move every 200ms
  let x = 0;
  setInterval(() => {
    x += 0.5;
    const move = { playerId: 'player1', x: x, y: 0, z: 0, ts: Date.now() };
    ws.send(JSON.stringify({ type: 'move', data: move }));
  }, 200);
});

ws.on('message', function incoming(data) {
  console.log('received: %s', data);
});
