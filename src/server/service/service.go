package service

import "encoding/json"

type Service interface {
	BroadcastMessage(timestamp int64, message string) error
}

type implSendMessageService struct {
	MessageQueue	*chan []byte
}

func NewSendMessageService(messageQueue *chan []byte) Service{
	return &implSendMessageService{MessageQueue: messageQueue}
}

func (i *implSendMessageService) BroadcastMessage(timestamp int64, message string) error  {
	mes := map[string]interface{}{
		"timestamp": timestamp,
		"message": message,
	}
	raw, err := json.Marshal(mes)
	if err != nil {
		return err
	}
	*i.MessageQueue <- raw
	return nil
}

