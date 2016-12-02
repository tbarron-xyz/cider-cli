package main

import (
	//"bufio"
	//"os"
	"fmt"
	//"strings"
	"encoding/json"
	"flag"
	"github.com/gorilla/websocket"
	"testing"
)

var _clients = flag.Int("C", 50, "Number of concurrent clients")
var clients int
var _pipeline = flag.Int("P", 64, "Number of requests to pipeline")
var pipeline int
var _verbose = flag.Bool("v", false, "verbose")

func NewWebsocket() (ws *websocket.Conn) {
	ws, _, _ = websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%d/", *_url, *_port), nil)
	return
}

func verbose(args ...interface{}) {
	if *_verbose {
		fmt.Println(args...)
	}
}

func writeread_n(ws *websocket.Conn, b []byte, n int) <-chan struct{} {
	c := make(chan struct{})
	go func() {
		for i := 0; i < n; i++ {
			ws.WriteMessage(Txtm, b)
			ws.ReadMessage()
		}
		c <- struct{}{}
	}()
	return c
}

func general_benchmark(precmd, cmd string, b *testing.B) {
	var wss = make([]*websocket.Conn, clients)
	for i := 0; i < clients; i++ {
		wss[i] = NewWebsocket()
	}
	wss[0].WriteMessage(Txtm, []byte(precmd))
	wss[0].ReadMessage()
	b.ResetTimer()
	args := make([]string, pipeline)
	for i := 0; i < pipeline; i++ {
		args[i] = cmd
	}
	tosend, _ := json.Marshal(args)
	chans := make([]<-chan struct{}, clients)
	for i, e := range wss {
		chans[i] = writeread_n(e, tosend, b.N/clients/pipeline)
	}
	for _, e := range chans {
		<-e
	}
}

func general_benchmark_multicore(precmd, cmd string, b *testing.B) {
	// cmd should have a "%d" where you want the client # to be substed
	var wss = make([]*websocket.Conn, clients)
	for i := 0; i < clients; i++ {
		wss[i] = NewWebsocket()
	}
	wss[0].WriteMessage(Txtm, []byte(precmd))
	wss[0].ReadMessage()
	var tosend = make([][]byte, clients)
	for i, _ := range wss {
		args := make([]string, pipeline)
		for j := 0; j < pipeline; j++ {
			args[j] = fmt.Sprintf(cmd, i)
		}
		tosend[i], _ = json.Marshal(args)
	}
	b.ResetTimer()
	chans := make([]<-chan struct{}, clients)
	for i, e := range wss {
		chans[i] = writeread_n(e, tosend[i], b.N/clients/pipeline)
	}
	for _, e := range chans {
		<-e
	}
}

func BenchmarkHGET(b *testing.B) {
	general_benchmark("HSET mykey myvalue `Long Ass String`", "HGET mykey myvalue", b)
}

func BenchmarkHGETMulticore(b *testing.B) {
	general_benchmark_multicore("", "HGET mykey%d myvalue", b)
}

func BenchmarkHSET(b *testing.B) {
	general_benchmark("", "HSET mykey myvalue `Long Ass String`", b)
}

func BenchmarkHSETMulticore(b *testing.B) {
	general_benchmark_multicore("", "HSET mykey%d myvalue `Long Ass String`", b)
}

func BenchmarkGET(b *testing.B) {
	general_benchmark("SET mykey `Long Ass String`", "GET mykey", b)
}

func BenchmarkSET(b *testing.B) {
	general_benchmark("", "SET mykey `Long Ass String`", b)
}

func init() {
	flag.Parse()
	clients = *_clients
	pipeline = *_pipeline
	fmt.Printf("clients: %v\npipeline: %v\n", clients, pipeline)
}
