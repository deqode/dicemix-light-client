package server

import (
	"flag"
	"log"
	"net/url"
	"os"

	"../commons"
	"../dc"
	"../nike"
	"../utils"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

// server configurations
var addr = flag.String("addr", "localhost:8080", "http service address")
var dialer = websocket.Dialer{} // use default options

// Exposed interfaces
var iNike nike.NIKE
var iDcNet dc.DC

type connection struct {
	Server
}

// NewConnection creates a new Server instance
func NewConnection() Server {
	initialize()

	return &connection{}
}

// Register - requests to C_JOIN_REQUEST
func (c *connection) Register(state *utils.State) {
	var connection = connect()
	listener(connection, state)

	defer connection.Close()
}

// performs some basic initializations
func initialize() {
	// initailze exposed interfaes for further use
	iNike = nike.NewNike()
	iDcNet = dc.NewDCNetwork()
}

// connects to server and extablishes a web socket connection
func connect() *websocket.Conn {
	url := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("Connecting to %s", url.String())
	conn, _, err := dialer.Dial(url.String(), nil)

	checkError(err)

	return conn
}

// listens for responses from server side
func listener(c *websocket.Conn, state *utils.State) {
	for {
		_, message, err := c.ReadMessage()
		checkError(err)

		response := &commons.GenericResponse{}
		err = proto.Unmarshal(message, response)
		checkError(err)

		// handles response and take further actions
		// based on response.Code
		handleMessage(c, message, response.Code, state)
	}
}

// checks for any potential errors
// exists program if one found
func checkError(err error) {
	if err != nil {
		log.Fatal("Error Occured:", err)
		os.Exit(1)
	}
}