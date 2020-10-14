package processing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"programs/pkgs/util"
	"programs/src/publishing_client/config"
)

type Client interface {
	PushMessage() error
}
type publishingClient struct {
	url    string
	config *config.Config
}

func NewPublishingClient(url string, config2 *config.Config) Client {
	return &publishingClient{
		url:    url,
		config: config2,
	}
}

func (i *publishingClient) PushMessage() error {
	uri := fmt.Sprintf("http://%s%s", i.config.Host, i.url)

	body := map[string]interface{}{
		"message": util.RandomString(10),
	}
	rawBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	request, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(rawBody))
	if err != nil {
		return err
	}
	var rs interface{}
	err = util.ExecuteRequest(request, &rs)
	if err != nil {
		return err
	}
	return err
}
