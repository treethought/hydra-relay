package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/websocket"
)

var portFlag = flag.String("port", "8088", "listen address")

func newClient(hub *Hub, port int) *Client {

	host := fmt.Sprintf("localhost:%d", port)
	url := url.URL{Scheme: "ws", Host: host, Path: "/ws"}

	fmt.Printf("connecting to %s\n", url.String())
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	return client
}

func runRepl(hub *Hub, port int) {

	fmt.Println("starting repl")
	client := newClient(hub, port)

	repl := Repl{client: client}
	// repl.help()
	repl.start()

}

func main() {
	flag.Parse()

	port, err := strconv.Atoi(*portFlag)
	if err != nil {
		log.Fatal("invalid port")
	}

	hub := newHub()
	go hub.run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	go func() {
		fmt.Printf("starting server on port: %d\n", port)
		err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}

	}()
	// fmt.Println("https://localhost:8000/?code=Y29uc3QlMjBzb2NrZXQlMjAlM0QlMjBuZXclMjBXZWJTb2NrZXQoJ3dzJTNBJTJGJTJGbG9jYWxob3N0JTNBODA4OCUyRndzJyklM0IlMjAlMEElMEElMEElMkYlMkYlMjBDb25uZWN0aW9uJTIwb3BlbmVkJTBBc29ja2V0LmFkZEV2ZW50TGlzdGVuZXIoJ29wZW4nJTJDJTIwZnVuY3Rpb24lMjAoZXZlbnQpJTIwJTdCJTBBJTIwJTIwc29ja2V0LnNlbmQoJ2dyZWV0aW5ncyUyMGZyb20lMjBoeWRyYScpJTNCJTBBJTdEKSUzQiUwQSUwQSUwQWZ1bmN0aW9uJTIwZXZhbFJlbGF5KGJsb2NrKSUyMCU3QiUwQSUyMCUyMHRyeSUyMCU3QiUwQSUwOSUyMCUyMGNvbnNvbGUubG9nKGV2YWwoYmxvY2spKSUwQSUyMCUyMCU3RCUyMGNhdGNoJTIwKGVyciklMjAlN0IlMEElMjAlMjAlMjAlMjBjb25zb2xlLmxvZyglMjJlcnJvciUyMGluJTIwcmVsYXklMjBtZXNzYWdlJTIwZXZhbCUyMiklMEElMjAlMjAlMjAlMjBzb2NrZXQuc2VuZChKU09OLnN0cmluZ2lmeShlcnIubWVzc2FnZSkpJTBBJTIwJTIwJTdEJTBBJTdEJTBBJTBBc29ja2V0LmFkZEV2ZW50TGlzdGVuZXIoJ21lc3NhZ2UnJTJDJTIwZnVuY3Rpb24lMjAoZXZlbnQpJTIwJTdCJTBBJTIwJTIwY29uc29sZS5sb2coZXZlbnQuZGF0YSklM0IlMEElMjAlMjBldmFsUmVsYXkoZXZlbnQuZGF0YSklM0IlMEElN0QpJTNCJTBB")
	runRepl(hub, port)

}
