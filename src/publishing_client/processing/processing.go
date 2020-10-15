package processing

import (
	"encoding/json"
	"fmt"
	"programs/pkgs/util"
	"programs/src/publishing_client/config"
)

type Client interface {
	PushMessage() (string, error)
}
type publishingClient struct {
	url    string
	config *config.Config
	utilInstance	util.Util
}

func NewPublishingClient(url string, config2 *config.Config, utilInstance util.Util) Client {
	return &publishingClient{
		url:    url,
		config: config2,
		utilInstance: utilInstance,
	}
}

func (i *publishingClient) PushMessage() (string, error) {
	uri := fmt.Sprintf("http://%s%s", i.config.Host, i.url)

	body := map[string]interface{}{
		"message": i.utilInstance.RandomString(10),
	}
	rawBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	var rs interface{}
	err = i.utilInstance.ExecuteRequest(rawBody, &rs, uri)
	if err != nil {
		return "", err
	}
	return string(rawBody), err
}
