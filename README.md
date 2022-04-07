# hydra-relay

A simple websocket server for relaying messages to hydra from non-hydra peers.

This was primarily created to send code to be evaluated within hydra from
neovim

Currently, this is just a simple websocket server that forwards messages to all
connected clients. Taken almost directly from gorilla websocket library's
[example](https://github.com/gorilla/websocket/tree/master/examples/chat).

## Usage

Start the relay server, passing the port to listen on (default 8088).

```
hydra-relay --port 8088
```

Open an instance of hydra, and paste the following:

```
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
```

Then send blocks of javascript to your relay to have them evaluated in hydra via
any websocket client.







