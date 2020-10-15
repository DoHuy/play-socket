package client

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/jpillora/backoff"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"net/url"
	"programs/pkgs/log"
	"programs/src/websocket_client/config"
	"sync"
	"time"
)

type WebSocketClient struct {
	Id         string
	Connection *websocket.Conn
	Config     *config.Config
	Logger     log.Logger
	// default to 2 seconds
	RecIntvlMin time.Duration
	// RecIntvlMax specifies the maximum reconnecting interval,
	// default to 30 seconds
	RecIntvlMax time.Duration
	// RecIntvlFactor specifies the rate of increase of the reconnection
	// interval, default to 1.5
	RecIntvlFactor float64
	// HandshakeTimeout specifies the duration for the handshake to complete,
	// default to 2 seconds
	HandshakeTimeout time.Duration
	// NonVerbose suppress connecting/reconnecting messages.
	NonVerbose bool
	mu          sync.Mutex
	reqHeader   http.Header
	httpResp    *http.Response
	dialErr     error
	isConnected bool
	dialer      *websocket.Dialer
}
// ErrNotConnected is returned when the application read/writes
// a message and the connection is closed
var ErrNotConnected = errors.New("websocket: not connected")

func NewWebSocketClient(config *config.Config, logger log.Logger, id string) *WebSocketClient {
	return &WebSocketClient{
		Id:         id,
		Logger:     logger,
		Config:     config,
	}
}

// Close closes the underlying network connection without
// sending or waiting for a close frame.
func (rc *WebSocketClient) Close() {
	rc.mu.Lock()
	if rc.Connection != nil {
		rc.Connection.Close()
	}
	rc.isConnected = false
	rc.mu.Unlock()
}

// Close And Recconect will try to reconnect.
func (rc *WebSocketClient) closeAndReConnect() {
	rc.Logger.Info("Try reconnect to server")
	rc.Close()
	go func() {
		rc.Connect()
	}()

}

// ReadMessage is a helper method for getting a reader
// using NextReader and reading from that reader to a buffer.
//
// If the connection is closed ErrNotConnected is returned
func (rc *WebSocketClient) ReadMessage() (messageType int, message []byte, err error) {
	err = ErrNotConnected
	if rc.IsConnected() {
		messageType, message, err = rc.Connection.ReadMessage()
		if err != nil {
			rc.Logger.Error("read message failed", zap.Error(err))
			rc.Logger.Warn("Server is not available")
			rc.closeAndReConnect()
		}
	}

	return
}

func (rc *WebSocketClient) WriteMessage(messageType int, data []byte) error {
	err := ErrNotConnected
	if rc.IsConnected() {
		err = rc.Connection.WriteMessage(messageType, data)
		if err != nil {
			rc.Logger.Error("write message failed", zap.Error(err))
			rc.Logger.Warn("Server is not available")
			rc.closeAndReConnect()
		}
	}

	return err
}

func (rc *WebSocketClient) Dial() {
	if rc.RecIntvlMin == 0 {
		rc.RecIntvlMin = 30 * time.Second
	}

	if rc.RecIntvlMax == 0 {
		rc.RecIntvlMax = 30 * time.Second
	}

	if rc.RecIntvlFactor == 0 {
		rc.RecIntvlFactor = 1.5
	}

	if rc.HandshakeTimeout == 0 {
		rc.HandshakeTimeout = 2 * time.Second
	}

	rc.dialer = websocket.DefaultDialer
	rc.dialer.HandshakeTimeout = rc.HandshakeTimeout

	go func() {
		rc.Connect()
	}()

	// wait on first attempt
	time.Sleep(rc.HandshakeTimeout)
}

func (rc *WebSocketClient) Connect() {
	b := &backoff.Backoff{
		Min:    rc.RecIntvlMin,
		Max:    rc.RecIntvlMax,
		Factor: rc.RecIntvlFactor,
		Jitter: true,
	}

	rand.Seed(time.Now().UTC().UnixNano())

	for {
		nextItvl := b.Duration()
		u := url.URL{Scheme: "ws", Host: rc.Config.Host, Path: "/ws"}
		wsConn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

		rc.mu.Lock()
		rc.Connection = wsConn
		rc.dialErr = err
		rc.isConnected = err == nil
		rc.mu.Unlock()

		if err == nil {
			if !rc.NonVerbose {
				rc.Logger.Info("Dial: connection was successfully established", zap.String("host", rc.Config.Host))
			}
			break
		} else {
			if !rc.NonVerbose {
				rc.Logger.Warn("Server is not available")
			}
		}

		time.Sleep(nextItvl)
	}
}

// GetDialError returns the last dialer error.
// nil on successful connection.
func (rc *WebSocketClient) GetDialError() error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	return rc.dialErr
}

// IsConnected returns the WebSocket connection state
func (rc *WebSocketClient) IsConnected() bool {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	return rc.isConnected
}