package main

import (
	"github.com/alldroll/spatial-index/geometry"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"os"
)

type AppConf struct {
	Host         string
	GoogleApiKey string
	Zoom         string
}

var (
	indexTemplate = template.Must(template.ParseFiles("cmd/example/public/index.html"))
	upgrader      = websocket.Upgrader{}
	appConf       = AppConf{}
	port          = ""
	service       *Service
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

		bounds := shape.NewBoundaryBox(
			shape.NewPoint(msg.Lat1, msg.Lng1),
			shape.NewPoint(msg.Lat2, msg.Lng2),
		)

		clusters := service.RangeQuery(bounds, msg.QuadKeys)

		res := make([]response, len(clusters))
		for i, p := range clusters {
			res[i] = response{
				Lat: p.GetX(),
				Lng: p.GetY(),
				Cnt: p.GetCount(),
			}
		}

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
		Zoom   string
	}{
		WSPath: "ws://" + appConf.Host + "/ws",
		ApiKey: appConf.GoogleApiKey,
		Zoom:   appConf.Zoom,
	}

	indexTemplate.Execute(w, data)
}

func main() {
	port = os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	appConf.GoogleApiKey = os.Getenv("GOOGLE_API_KEY")
	appConf.Zoom = os.Getenv("ZOOM")
	appConf.Host = os.Getenv("HOST")
	if appConf.Host == "" {
		appConf.Host = "localhost:" + port
	}

	service = NewService()

	r := mux.NewRouter()
	r.HandleFunc("/", serveHome)
	r.HandleFunc("/ws", serveWS).Name("wsRoute")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./cmd/example/public/")))

	log.Fatal(http.ListenAndServe(":"+port, r))
}
