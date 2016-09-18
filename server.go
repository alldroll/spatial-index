package main

import (
	//"github.com/alldroll/quadtree/quadtree"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
)

var (
	indexTemplate = template.Must(template.ParseFiles("client/index.html"))
	upgrader      = websocket.Upgrader{}
)

func serveWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("something with upgrader", err)
		return
	}

	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("something with read message", err)
			break
		}

		log.Printf("recieve: %s", message)
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Printf("something with write message", err)
			break
		}
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if !(r.Method == "GET" && r.URL.Path == "/") {
		http.Error(w, "Not found", 404)
		return
	}

	indexTemplate.Execute(w, "ws://"+r.Host+"/ws")
}

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWS)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
