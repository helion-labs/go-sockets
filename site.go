package main

import (
	_"encoding/json"
	"log"
	"html/template"
	"net/http"
	"github.com/gorilla/websocket"
	"strings"
	"fmt"
)

var upgrader = websocket.Upgrader{}

func index(w http.ResponseWriter, req *http.Request) {
	log.Print("Executing template")
	t, _ := template.ParseFiles("main.html")
	t.Execute(w, "null")
}

func upgrade(w http.ResponseWriter, r *http.Request) {
	log.Print("Upgrading websocket...")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer func() {
		log.Print("Closing websocket...")
		c.Close()
	}()

	for {
		distance_json := <- Distance_chan
		fmt.Println("Writting -> ", string(distance_json), " to site!")
	  c.WriteMessage(websocket.TextMessage , distance_json)
	}
}

// sanitized the file-server not to show the "root"
// at ./static/
func neuter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main_site() {
	fmt.Println("Starting site!")

	mux := http.NewServeMux()
	mux.HandleFunc("/index", index)
	mux.HandleFunc("/upgrade", upgrade)
	mux.HandleFunc("/", index)
	
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", neuter(fileServer)))
	log.Fatal(http.ListenAndServe(":8090", mux))
}
