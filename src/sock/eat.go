package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	s "github.com/student020341/go-server-core/TWTServer"
	"golang.org/x/net/websocket"
)

func HandleWeb(w http.ResponseWriter, r *http.Request, path []string) {

	router.Handle(w, r, path)
}

func GetName() string {
	return "sock"
}

func handleHome(w http.ResponseWriter, r *http.Request, args map[string]interface{}) {
	fmt.Fprintln(w, "Use the /connect route to start a socket connection")
}

func handleConnect(w http.ResponseWriter, r *http.Request, args map[string]interface{}) {
	wserver := websocket.Server{Handler: websocket.Handler(socketStuff)}
	wserver.ServeHTTP(w, r)
}

type simpleClient struct {
	Id     int
	Socket *websocket.Conn
	// buffer position in client struct
	Position Vector2
	Name     string
}

type Vector2 struct {
	X float64
	Y float64
}

// list of connected clients
var clients []*simpleClient

// id=0 is the server
var lastId int = 1

func socketStuff(ws *websocket.Conn) {
	fmt.Println("incoming connection:", ws.Request().RemoteAddr)

	clientId := lastId
	lastId += 1

	// tell client its id
	websocket.Message.Send(ws, fmt.Sprintf(`{"id":%d}`, clientId))

	// add this client to the client list
	wsClient := simpleClient{Id: clientId, Socket: ws}
	clients = append(clients, &wsClient)

	var in []byte
	// serve client until it leaves
	for {
		if err := websocket.Message.Receive(ws, &in); err != nil {
			break
		}

		var obj map[string]interface{}
		err := json.Unmarshal(in, &obj)
		if err != nil {
			// weird that user sent nothing, but can be ignored
			if err != io.EOF {
				fmt.Println(err)
			}
			// ignore errors
			continue
		}

		fmt.Printf("client %v says...\n%v\n", wsClient.Id, obj)

		// message, haveMessage := obj["message"]
		// // todo: move to function or use broadcast other at least
		// if (haveMessage) {
		// 	broadcastToOther(clientId,map[string]interface{}{
		// 		"message": message,
		// 		"client": clientId,
		// 	})
		// }
	}

	// after loop

	// drop client from server array
	for index, client := range clients {
		if client.Socket == ws {
			// splice client out of collection
			clients = append(clients[:index], clients[index+1:]...)
		}
	}

	fmt.Printf("client %v disconnected\n", clientId)
}

var router s.SubRouter

func init() {
	// setup router
	router.Register("/", "GET", handleHome)
	router.Register("/connect", "GET", handleConnect)
}

func main() {
	fmt.Println("this is")
}
