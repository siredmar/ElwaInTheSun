package server

import (
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var clients = make(map[*websocket.Conn]bool)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Serve WebSocket connections
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorln("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	clients[conn] = true
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			delete(clients, conn)
			break
		}
	}
}

// Serve frontend
func ServeFrontend(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "html/index.html")
}

// Broadcast updates to all WebSocket clients
func BroadcastUpdate(data string) {
	for client := range clients {
		if err := client.WriteMessage(websocket.TextMessage, []byte(data)); err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}

// // Simulated data update loop
// func StartDataLoop() {
// 	for {
// 		time.Sleep(5 * time.Second) // Simulate new data every 5 seconds
// 		data := fmt.Sprintf("{\"power\": %d, \"timestamp\": \"%s\"}", time.Now().Unix()%5000, time.Now().Format(time.RFC3339))
// 		broadcastUpdate(data)
// 	}
// }

// // Broadcast updates to all WebSocket clients
// func broadcastUpdate(data string) {
// 	for client := range clients {
// 		if err := client.WriteMessage(websocket.TextMessage, []byte(data)); err != nil {
// 			client.Close()
// 			delete(clients, client)
// 		}
// 	}
// }
