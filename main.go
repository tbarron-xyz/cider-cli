package main

import (
	"bufio"
	"os"
	"fmt"
	"strings"
	"github.com/gorilla/websocket"
)

var Txtm = websocket.TextMessage
var reader = bufio.NewReader(os.Stdin)

func readOnce (ws *websocket.Conn) (err error) {
	var msg []byte
	_, msg, err = ws.ReadMessage()
	if err != nil {
		fmt.Println("Read error: " + err.Error())
		ws.Close()
	}
	fmt.Println("<", string(msg))
	return
}

func main () {
	fmt.Println("cider cli")
	fmt.Print("Connecting... ")
	ws, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%d/", *_url, *_port), nil)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}
	fmt.Println("Connected.")
	for {
		fmt.Print("> ")
		var input []byte
    	input, _ = reader.ReadBytes('\n')
    	err = ws.WriteMessage(Txtm, []byte(strings.TrimSpace(string(input))))
    	if err != nil {
    		fmt.Println("Write error: " + err.Error())
    		ws.Close()
    		break
    	}
    	readOnce(ws)
	}
}