package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	log2 "programs/pkgs/log"
	"programs/pkgs/util"
	"programs/src/websocket_client/client"
	config2 "programs/src/websocket_client/config"
	"time"
)

func main() {
	config, err := config2.GetConfig()
	if err != nil {
		panic(err)
	}
	//// init logger
	logger, err := log2.NewLoggerService(config.EnvMode, config.WebSocketLogPath)
	if err != nil {
		panic(err)
	}
	//// init websocket
	id := util.NewUtil().GenerateRandomIdentifier()
	socketClient := client.NewWebSocketClient(config, logger, id)
	socketClient.Dial()
	defer socketClient.Close()
	for {
		flag.Parse()
		log.SetFlags(0)

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)
		done := make(chan struct{})

		go func() {
			defer close(done)
			for {
				_, message, err := socketClient.ReadMessage()
				if err == nil {
					logger.Info("received message from server", zap.String("content", string(message)))
				}

				err = socketClient.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("websocket client (%s) received message %s", socketClient.Id, string(message))))
				if err == nil {
					logger.Info("send message to server", zap.String("content", fmt.Sprintf("websocket client (%s) received message %s", socketClient.Id, string(message))))
				}
			}
		}()

		for {
			select {
			case <-done:
				return
			case <-interrupt:
				// Cleanly close the connection by sending a close message and then
				// waiting (with timeout) for the server to close the connection.
				err := socketClient.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					return
				}
				select {
				case <-done:
				case <-time.After(time.Second):
				}
				return
			}
		}
	}

}
