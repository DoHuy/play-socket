package service

import (
	"encoding/json"
	"io/ioutil"
	"programs/src/server/config"
)

type Service interface {
	BroadcastMessage(timestamp int64, message string) error
}

type implSendMessageService struct {
	config *config.Config
}

func NewSendMessageService(config *config.Config) Service {
	return &implSendMessageService{config: config}
}

func (i *implSendMessageService) BroadcastMessage(timestamp int64, message string) error {
	mes := map[string]interface{}{
		"timestamp": timestamp,
		"message":   message,
	}
	raw, err := json.Marshal(mes)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(i.config.TemporaryFile, raw, 0777); err != nil {
		return err
	}
	return nil
}
