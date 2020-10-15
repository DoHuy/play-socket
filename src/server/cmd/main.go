package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"programs/pkgs/log"
	"programs/pkgs/util"
	"programs/src/server/config"
	"programs/src/server/handler"
)

func main() {

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
	utilInstance := util.NewUtil()
	// init new handler
	handler := handler.NewHandler(config, logger, upgrade, utilInstance)

	// run server
	http.HandleFunc("/broadcast", func(writer http.ResponseWriter, request *http.Request) {
		handler.BroadCastMessageHandler(writer, request)
		return
	})
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		handler.WebSocketHandler(writer, request)
	})

	fmt.Println("server listening on  ", config.ListenAddress)

	http.ListenAndServe(config.ListenAddress, nil)
}
