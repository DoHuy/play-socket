package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"programs/pkgs/log"
	"programs/src/server/config"
	"programs/src/server/handler"
)

func main()  {

	config, err := config.GetConfig()
	if err != nil {
		panic(err)
	}
	// init logger service
	logger, err := log.NewLoggerService(config.EnvMode, config.ServerLogPath)
	if err != nil {
		panic(err)
	}
	//init websocket instance
	upgrade := websocket.Upgrader{}
	// init messageQueue
	messageQueue := make(chan []byte, config.MaximumQueueSize)
	// init new handler
	handler := handler.NewHandler(config, logger, upgrade, &messageQueue)

	// run server
	http.HandleFunc("/broadcast", func(writer http.ResponseWriter, request *http.Request) {
		handler.BroadCastMessageHandler(writer, request)
		return
	})
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		handler.WebSocketHandler(writer, request)
	})

	http.ListenAndServe(config.ListenAddress, nil)
	fmt.Println("server listening on  ", config.ListenAddress)
}