// Establish connection to server via WebSocket
// @returns conn or null
function connect() {
  if (window["WebSocket"]) {
    return new WebSocket("ws://localhost:9000/ws");
  } else {
    return null;
  }
}
