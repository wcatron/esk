package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/markbates/pkger"
	"github.com/wcatron/esk/pkg/config"
	"github.com/wcatron/esk/pkg/datasource"
	"github.com/wcatron/esk/pkg/websocket"
)

func handleConnection(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("main:handleConnection:Connection attempting...")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "main:handleConnection:error:%+v\n", err)
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
		ID:   pool.NextClientId(),
	}
	fmt.Println("main:handleConnection:New client...")

	pool.Register <- client
	client.Read()
}

func setupStaticPlayground() {
	fs := http.FileServer(pkger.Dir("/playground/build"))
	http.Handle("/", fs)
}

func setupWebsocketEndpoint() {
	datasource := datasource.NewDataSource()
	go datasource.Listen()
	config := config.NewConfig(datasource.GenericHandler)
	pool := websocket.NewPool(config)
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleConnection(pool, w, r)
	})
}

func main() {
	port := flag.String("port", "8080", "Port to start the server on.")
	flag.Parse()
	fmt.Println("main:Starting server...")
	setupStaticPlayground()
	setupWebsocketEndpoint()
	fmt.Printf("main:Listening at http://localhost:%s and ws://localhost:%s/ws\n", *port, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
