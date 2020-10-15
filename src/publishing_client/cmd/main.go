package main

import (
	"go.uber.org/zap"
	"programs/pkgs/log"
	"programs/pkgs/util"
	"programs/src/publishing_client/config"
	"programs/src/publishing_client/processing"
	"time"
)

func main() {

	config, err := config.GetConfig()
	if err != nil {
		panic(err)
	}
	// init logger service
	logger, err := log.NewLoggerService(config.EnvMode, config.PublishingLogPath)
	if err != nil {
		panic(err)
	}
	utilInstance :=  util.NewUtil()
	client := processing.NewPublishingClient("/broadcast", config, utilInstance)
	// send with time T
	ticker := time.NewTicker(time.Duration(config.IntervalTime) * time.Second)
	stopSignal := make(chan bool)
	for {
		select {

		// Case statement
		case <-stopSignal:
			return

		// Case to print current time
		case <-ticker.C:
			//todo
			err := client.PushMessage()
			if err != nil {
				logger.Error("push message failed", zap.Error(err))
				logger.Info("server is not available")
			}
		}
	}

}
