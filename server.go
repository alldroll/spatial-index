package main

import (
	"encoding/json"
	"github.com/alldroll/spatial-index/clustering"
	"github.com/alldroll/spatial-index/quadtree"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	Kdb = 131.78635917382815
	Kdf = 131.9923528261719
	Pdb = 43.224515498757405
	Pdf = 43.024050275744735
)

type AppConf struct {
	Host         string
	GoogleApiKey string
	Zoom         int
}

type Markets struct {
	Markers []struct {
		Lat string
		Lon string
	}
}

var (
	indexTemplate                    = template.Must(template.ParseFiles("client/index.html"))
	upgrader                         = websocket.Upgrader{}
	appConf                          = AppConf{}
	db            *quadtree.QuadTree = nil
)

type TypeRes struct {
	Lat, Lon float64
	Cnt      int
}

func serveWS(w http.ResponseWriter, r *http.Request) {
	log.Printf("new connection")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("something with upgrader", err)
		return
	}

	defer conn.Close()

	msg := struct {
		Lat1 float64
		Lon1 float64
		Lat2 float64
		Lon2 float64
		Zoom int
	}{}

	for {
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("something with read msg", err)
			break
		}

		points, _ := db.GetPoints(msg.Lon1, msg.Lat1, msg.Lon2, msg.Lat2)

		zoom := (msg.Zoom - appConf.Zoom) + 3
		if zoom <= 0 {
			zoom = 0
		}

		grid := cluster.NewGrid(Pdf, Kdb, Pdb, Kdf, zoom)

		grid.AddChunk(points)

		clusters := grid.GetClusters()

		res := make([]TypeRes, len(clusters))
		for i, p := range clusters {
			center := p.GetCenter()
			res[i] = TypeRes{
				Lat: center.GetX(),
				Lon: center.GetY(),
				Cnt: p.GetLength(),
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

func loadMarkets() {
	points := Markets{}
	file, _ := os.Open("markets_points.json")
	decoder := json.NewDecoder(file)
	e := decoder.Decode(&points)

	if e != nil {
		log.Fatal(e)
		return
	}

	for i, point := range points.Markers {
		x, _ := strconv.ParseFloat(point.Lat, 64)
		y, _ := strconv.ParseFloat(point.Lon, 64)
		db.Insert(x, y)
		if false && i == 50 {
			break
		}
	}
}

func main() {
	err := readConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	db, err = quadtree.NewQuadTree(Pdf, Kdb, Pdb, Kdf, 20)
	if err != nil {
		log.Fatal(err)
		return
	}

	loadMarkets()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWS)
	err = http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
