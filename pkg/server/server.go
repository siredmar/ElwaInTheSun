package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/siredmar/ElwaInTheSun/pkg/args"
	"github.com/siredmar/ElwaInTheSun/pkg/mypv"
	"github.com/siredmar/ElwaInTheSun/pkg/sonnen"
	log "github.com/sirupsen/logrus"
)

var clients = make(map[*websocket.Conn]bool)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Server struct {
	sonnenClient *sonnen.Client
	mypvClient   *mypv.Client
}

func NewServer(sonnenClient *sonnen.Client, mypvClient *mypv.Client) *Server {
	s := &Server{
		sonnenClient: sonnenClient,
		mypvClient:   mypvClient,
	}

	http.HandleFunc("/", s.ServeFrontend)
	http.HandleFunc("/ws", s.HandleWebSocket)
	http.HandleFunc("/settings", s.GetConfigHandler)
	http.HandleFunc("/powerdata", s.PowerDataHandler)
	return s
}

func (s *Server) Run() error {
	// go server.StartDataLoop()
	log.Infof("Server running on :%d\n", args.Port)
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", args.Port), nil)
}

// Serve WebSocket connections
func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
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
func (s *Server) ServeFrontend(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "html/index.html")
}

// Broadcast updates to all WebSocket clients
func (s *Server) BroadcastUpdate(data string) {
	for client := range clients {
		if err := client.WriteMessage(websocket.TextMessage, []byte(data)); err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}

type PowerData struct {
	GridWatt          int `json:"gridWatt"`
	ConsumptionWatt   int `json:"consumptionWatt"`
	HeaterWatt        int `json:"heaterWatt"`
	BatteryPercentage int `json:"batteryPercentage"`
	BatteryWatt       int `json:"batteryWatt"`
	ProductionWatt    int `json:"productionWatt"`
	Temp1             int `json:"temp1"`
	Temp2             int `json:"temp2"`
}

func (s *Server) PowerDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		liveData, err := s.mypvClient.LiveData()
		if err != nil {
			http.Error(w, "Failed to load power data", http.StatusInternalServerError)
			return
		}
		status, err := s.sonnenClient.Status()
		if err != nil {
			http.Error(w, "Failed to load sonnen status", http.StatusInternalServerError)
			return
		}

		data := PowerData{
			GridWatt:          int(status.GridFeedInW),
			ConsumptionWatt:   int(status.ConsumptionW),
			HeaterWatt:        liveData.PowerElwa2,
			BatteryWatt:       int(status.PacTotalW),
			BatteryPercentage: int(status.Rsoc),
			ProductionWatt:    int(status.ProductionW),
			Temp1:             liveData.Temp1,
			Temp2:             liveData.Temp2,
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			http.Error(w, "Failed to load power data", http.StatusInternalServerError)
			return
		}
		_, err = w.Write(jsonData)
		if err != nil {
			log.Println("Error writing response:", err)
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
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
