package util

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"math/rand"
	"net/http"
)

const (
	chars = "abcdefghijklmnopqrstuvwxyz0123456789"
)

// BaseResponse represents struct of base response
type BaseResponse struct {
	Data interface{} `json:"data,omitempty"`
	Meta Meta `json:"meta,omitempty"`
}

type Meta struct{
	Message string `json:"message,omitempty"`
	Code	int    `json:"code,omitempty"`
}

func RandomString(l uint) string {
	s := make([]byte, l)
	for i := 0; i < int(l); i++ {
		s[i] = chars[rand.Intn(len(chars))]
	}
	return string(s)
}

func GenerateRandomIdentifier() string {
	return uuid.NewV4().String()
}

func BuildResponse(code int, body interface{}, message string) BaseResponse {
	return BaseResponse{Data: body, Meta: Meta{Message: message, Code: code}}
}


func ExecuteRequest(request *http.Request, result interface{}) error {
	httpClient := http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}

	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}

	respBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(respBytes, &result)
	if err != nil {
		return err
	}
	return nil
}

