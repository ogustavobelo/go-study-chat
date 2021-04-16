package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
)

var wsChan = make(chan WebSocketPayload)
var clients = make(map[WebSocketConnection]string)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WebSocketPayload struct {
	Action   string              `json:"action"`
	Username string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"-"`
}

type WebSocketConnection struct {
	*websocket.Conn
}

type WebSocketResponse struct {
	Action        string   `json:"action"`
	Message       string   `json:"message"`
	MessageType   string   `json:"message_type"`
	ConnectedUser []string `json:"connected_users"`
}

func WebSocketEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected to endpoint")

	webSocketResponse := WebSocketResponse{Message: `<em><small>Connected to server</small></em>`}

	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""

	err = ws.WriteJSON(webSocketResponse)
	if err != nil {
		log.Println(err)
	}
	go ListenForWs(&conn)
}

func ListenForWs(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error on Listen")
			log.Println("Error:", fmt.Sprintf("%v", r))
		}
	}()

	var payload WebSocketPayload

	for {
		err := conn.ReadJSON(&payload)
		if err == nil {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

func ListenToWebSocketChannel() {
	var response WebSocketResponse

	for {
		event := <-wsChan
		switch event.Action {
		case "username":
			clients[event.Conn] = event.Username
			users := getUsersList()
			response.Action = "list_users"
			response.ConnectedUser = users
			broadcastToAll(response)
		case "broadcast":
			response.Action = event.Action
			response.Message = fmt.Sprintf("<strong>%s</strong>: %s", event.Username, event.Message)
			broadcastToAll(response)
		case "left":
			response.Action = "list_users"
			delete(clients, event.Conn)
			users := getUsersList()
			response.ConnectedUser = users
			broadcastToAll(response)
		}

	}
}

func getUsersList() []string {
	var userList []string
	for _, client := range clients {
		if client != "" {
			userList = append(userList, client)
		}
	}

	sort.Strings(userList)
	return userList
}

func broadcastToAll(response WebSocketResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("WebSocket Error")
			_ = client.Close()
			delete(clients, client)
		}
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		log.Println(err)
	}

}

func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/static/favicon.ico")
}

func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println(err)
		return err
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
