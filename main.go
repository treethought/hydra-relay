package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var port = flag.String("port", "8088", "listen address")

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	fmt.Printf("starting server on port: %s", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", *port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
