package main

import (
	"encoding/json"
	"github.com/alldroll/spatial-index/geometry"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type AppConf struct {
	Host         string
	GoogleApiKey string
	Zoom         int
}

var (
	indexTemplate = template.Must(template.ParseFiles("public/index.html"))
	upgrader      = websocket.Upgrader{}
	appConf       = AppConf{}
	service       *TileService
)

func serveWS(w http.ResponseWriter, r *http.Request) {
	log.Printf("new connection")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("something with upgrader", err)
		return
	}

	defer conn.Close()

	type response struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
		Cnt int     `json:"count"`
	}

	msg := struct {
		Lat1     float64
		Lng1     float64
		Lat2     float64
		Lng2     float64
		QuadKeys []string
		Zoom     int
	}{}

	for {
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("something with read msg", err)
			break
		}

		start := time.Now()

		bounds := shape.NewBoundaryBox(
			shape.NewPoint(msg.Lat1, msg.Lng1),
			shape.NewPoint(msg.Lat2, msg.Lng2),
		)

		clusters := service.RangeQueryQuadKeys(bounds, msg.QuadKeys)

		res := make([]response, len(clusters))
		for i, p := range clusters {
			res[i] = response{
				Lat: p.GetX(),
				Lng: p.GetY(),
				Cnt: p.GetCount(),
			}
		}

		elapsed := time.Since(start)

		log.Printf("Took %s\n", elapsed)

		err = conn.WriteJSON(res)
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
	file, _ := os.Open("config/config.json")
	decoder := json.NewDecoder(file)
	return decoder.Decode(&appConf)
}

func initService() {
	service = NewTileService(
		NewTileRepo("config/vl_points.json"),
	)
}

func main() {

	err := readConfig()

	if err != nil {
		log.Fatal(err)
		return
	}

	initService()

	r := mux.NewRouter()

	r.HandleFunc("/", serveHome)
	r.HandleFunc("/ws", serveWS)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

	log.Fatal(http.ListenAndServe(":8080", r))
}
