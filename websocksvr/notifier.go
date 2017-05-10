package main

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

//Notifier represents a web sockets notifier
type Notifier struct {
	eventq  chan interface{}
	clients map[*websocket.Conn]bool
	//TODO: add other fields you might need
	//such as another channel or a mutex
	//(either would work)
	//remember that go maps ARE NOT safe for
	//concurrent access, so you must do something
	//to protect the `clients` map
	mu sync.RWMutex
}

//NewNotifier constructs a new Notifer.
func NewNotifier() *Notifier {
	//TODO: create, initialize and return
	//a Notifier struct
	no := &Notifier{
		eventq:  make(chan interface{}),
		clients: make(map[*websocket.Conn]bool),
		mu:      sync.RWMutex{},
	}
	return no
}

//Start begins a loop that checks for new events
//and broadcasts them to all web socket clients.
//This function should be called on a new goroutine
//e.g., `go mynotifer.Start()`
func (n *Notifier) Start() {
	//TODO: implement this function
	//this should check for new events written
	//to the `eventq` channel, and broadcast
	//them to all of the web socket clients
	for {
		event := <-n.eventq
		n.broadcast(event)
	}
}

//AddClient adds a new web socket client to the Notifer
func (n *Notifier) AddClient(client *websocket.Conn) {
	//TODO: implement this
	//But remember that this will be called from
	//an HTTP handler, and each HTTP request is
	//processed on its own goroutine, so your
	//implementation here MUST be safe for concurrent use

	fmt.Println("Add Client")

	n.mu.Lock()
	defer n.mu.Unlock()

	n.clients[client] = true
	go n.readPump(client)
}

//Notify will add a new event to the event queue
func (n *Notifier) Notify(event interface{}) {
	//TODO: add the `event` to the `eventq`
	fmt.Println("event add started")
	n.eventq <- event
	fmt.Println("event add finished")
}

//readPump will read all messages (including control messages)
//send by the client and ignore them. This is necessary in order
//process the control messages. If you don't do this, the
//websocket will get stuck and start producing errors.
//see https://godoc.org/github.com/gorilla/websocket#hdr-Control_Messages
func (n *Notifier) readPump(client *websocket.Conn) {
	//TODO: implement this according to the notes in the
	//Control Message section of the Gorilla Web Socket docs:
	//https://godoc.org/github.com/gorilla/websocket#hdr-Control_Messages
	for {
		if _, _, err := client.NextReader(); err != nil {
			client.Close()
			fmt.Printf("Error reading: %v\n", err.Error())
			break
		}
	}
}

//broadcast sends the event to all client as a JSON-encoded object
func (n *Notifier) broadcast(event interface{}) {
	//TODO: Loop over all of the web socket clients in
	//n.clients and write the `event` parameter to the client
	//as a JSON-encoded object.
	//HINT: https://godoc.org/github.com/gorilla/websocket#Conn.WriteJSON
	//and for even better performance, try using a PreparedMessage:
	//https://godoc.org/github.com/gorilla/websocket#PreparedMessage
	//https://godoc.org/github.com/gorilla/websocket#Conn.WritePreparedMessage
	for client, _ := range n.clients {
		n.mu.Lock()
		err := client.WriteJSON(event)
		if err != nil {
			client.Close()
			delete(n.clients, client)
			fmt.Printf("Error broadcasting: %v\n", err.Error())
			n.mu.Unlock()
			return
		}
		n.mu.Unlock()
	}
	//If you get an error while writing to a client,
	//the client has wandered off, so you should call
	//the `.Close()` method on the client, and delete
	//it from the n.clients map
}
