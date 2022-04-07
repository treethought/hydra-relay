package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Repl struct {
	client *Client
}

func (r Repl) help() {
	fmt.Println("hyrda-relay")
	fmt.Println("Open hydra and evaluate the following:")
	fmt.Println()
	fmt.Println(`
// replace with your relay address
const socket = new WebSocket('ws://localhost:8088/ws'); 


// Connection opened
socket.addEventListener('open', function (event) {
  socket.send('greetings from hydra');
});


function evalRelay(block) {
  try {
	  console.log(eval(block))
  } catch (err) {
    console.log("error in relay message eval")
    socket.send(JSON.stringify(err.message))
  }
}

socket.addEventListener('message', function (event) {
  console.log(event.data);
  evalRelay(event.data);
});
    `)

}

func (r Repl) start() {
	block := []string{}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		// reads user input until \n by default
		scanner.Scan()
		// Holds the string that was scanned
		text := scanner.Text()
		if len(text) != 0 {
			fmt.Println(text)
			block = append(block, text)
		} else {
			// exit if user entered an empty string

			combined := strings.Join(block, "\n")
			fmt.Println("sending:")
			fmt.Println(combined)
			fmt.Println()
			r.client.hub.broadcast <- []byte(combined)
			combined = ""
			block = []string{}
		}

	}

}
