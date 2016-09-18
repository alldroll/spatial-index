package main

import (
	//"github.com/alldroll/quadtree/quadtree"
	"encoding/json"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"os"
)

type AppConf struct {
	Host         string
	GoogleApiKey string
	Zoom         int
}

var (
	indexTemplate = template.Must(template.ParseFiles("client/index.html"))
	upgrader      = websocket.Upgrader{}
	appConf       = AppConf{}
)

func serveWS(w http.ResponseWriter, r *http.Request) {
	log.Printf("new connection")

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

	data := struct {
		WSPath string
		ApiKey string
		Zoom   int
	}{
		WSPath: "ws://" + appConf.Host + "/ws",
		ApiKey: appConf.GoogleApiKey,
		Zoom:   appConf.Zoom,
	}

	indexTemplate.Execute(w, data)
}

func readConfig() error {
	file, _ := os.Open("config/conf.json")
	decoder := json.NewDecoder(file)
	return decoder.Decode(&appConf)
}

func main() {
	err := readConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWS)
	err = http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
