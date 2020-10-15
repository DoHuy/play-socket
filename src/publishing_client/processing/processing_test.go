package processing

import (
	"errors"
	mock2 "github.com/stretchr/testify/mock"
	"programs/pkgs/util"
	"programs/src/publishing_client/config"
	"programs/src/publishing_client/processing/mocks"
	"testing"
)

func Test_publishingClient_PushMessage(t *testing.T) {
	mock:=new(mocks.Util)
	mockFailed:=new(mocks.Util)
	var rs interface{}
	//httpMock, _ := http.NewRequest(http.MethodPost, "/broadcast", bytes.NewBufferString("test"))
	mock.On("ExecuteRequest", mock2.Anything, &rs, "http://0.0.0.0:8080/broadcast").Return(nil)
	mock.On("RandomString", uint(10)).Return("12342edsdf")
	mockFailed.On("ExecuteRequest",mock2.Anything, &rs, "http://0.0.0.0:8080/broadcast").Return(errors.New("server is not available"))
	mockFailed.On("RandomString", uint(10)).Return("12342edsdf")
	type fields struct {
		url    string
		config *config.Config
	}
	tests := []struct {
		name    string
		fields  fields
		utilInstance util.Util
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Request broadcast failed, server is not available",
			fields: fields{
				config: &config.Config{
					Host: "0.0.0.0:8080",
				},
				url: "/broadcast",
			},
			utilInstance: mockFailed,
			wantErr: true,
		},
		{
			name: "Request broadcast success",
			fields: fields{
				config: &config.Config{
					Host: "0.0.0.0:8080",
				},
				url: "/broadcast",

			},
			utilInstance: mock,
			wantErr: false,
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &publishingClient{
				url:    tt.fields.url,
				config: tt.fields.config,
				utilInstance: tt.utilInstance,
			}
			if err := i.PushMessage(); (err != nil) != tt.wantErr {
				t.Errorf("PushMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
