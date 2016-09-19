package main

import (
	"encoding/json"
	"github.com/alldroll/quadtree/quadtree"
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
	quadTree      *quadtree.QuadTree = nil
)

type TypeRes struct {
	Lat, Lon float64
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

		points, _ := quadTree.GetPoints(coords.Lon1, coords.Lat1, coords.Lon2, coords.Lat2)

		res := make([]TypeRes, len(points))
		for i, p := range points {
			res[i] = TypeRes{
				Lat: p.GetX(),
				Lon: p.GetY(),
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
	for i := 1; i < 500; i++ {
		a := Pdf + rand.Float64()*(Pdb-Pdf)
		b := Kdb + rand.Float64()*(Kdf-Kdb)
		quadTree.Insert(a, b)
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

	quadTree, err = quadtree.NewQuadTree(Pdf, Pdb, Kdb, Kdf, 10)
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
