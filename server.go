package main

import (
	"encoding/json"
	"github.com/alldroll/spatial-index/clustering"
	"github.com/alldroll/spatial-index/quadtree"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
)

type AppConf struct {
	Host         string
	GoogleApiKey string
	Zoom         int
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

	coords := struct {
		Lat1 float64
		Lon1 float64
		Lat2 float64
		Lon2 float64
	}{}

	for {
		err := conn.ReadJSON(&coords)
		if err != nil {
			log.Printf("something with read message", err)
			break
		}

		log.Printf("recieve: %#v\n", coords)

		points, _ := db.GetPoints(coords.Lon1, coords.Lat1, coords.Lon2, coords.Lat2)

		spI, _ := quadtree.NewQuadTree(coords.Lon1, coords.Lat1, coords.Lon2, coords.Lat2, 20)

		eps := (coords.Lat2 - coords.Lat1) / 10

		clustering := cluster.NewClustering(
			spI,
			points,
			eps,
			3,
		)

		clusters := clustering.DBScan(points, eps, 3)

		log.Printf("recieve: %#v\n", clusters)

		res := make([]TypeRes, len(clusters))
		for i, p := range clusters {
			res[i] = TypeRes{
				Lat: p.GetCenter().GetX(),
				Lon: p.GetCenter().GetY(),
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
		Zoom   int
	}{
		WSPath: "ws://" + appConf.Host + "/ws",
		ApiKey: appConf.GoogleApiKey,
		Zoom:   appConf.Zoom,
	}

	indexTemplate.Execute(w, data)
}

const (
	Kdb = 131.78635917382815
	Kdf = 131.9923528261719
	Pdb = 43.224515498757405
	Pdf = 43.024050275744735
)

func generatePoints() {
	for i := 0; i < 400; i++ {
		a := Pdf + rand.Float64()*(Pdb-Pdf)
		b := Kdb + rand.Float64()*(Kdf-Kdb)
		db.Insert(a, b)
	}
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

	db, err = quadtree.NewQuadTree(Pdf, Pdb, Kdb, Kdf, 20)
	if err != nil {
		log.Fatal(err)
		return
	}

	generatePoints()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWS)
	err = http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
