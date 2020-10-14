package handler

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"programs/pkgs/log"
	"programs/pkgs/util"
	"programs/src/server/config"
	"programs/src/server/service"
	"time"
)

type Handler struct {
	config *config.Config
	upgrade	websocket.Upgrader
	logger	log.Logger
	messageQueue *chan []byte
	service	service.Service

}

func NewHandler(config *config.Config, logger log.Logger, upgrade websocket.Upgrader, messageQueue *chan []byte) *Handler {
	serviceInstance := service.NewSendMessageService(messageQueue)
	return &Handler{
		config: config,
		upgrade: upgrade,
		logger: logger,
		messageQueue: messageQueue,
		service: serviceInstance,
	}
}
/*
impl health check for server i think every apis server need it,
 */
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	resp := util.BuildResponse(http.StatusOK,nil, "Alive")
	json.NewEncoder(w).Encode(resp)
	return
}

func (h *Handler) BroadCastMessageHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		h.logger.Warn("not support the others Method")
		json.NewEncoder(w).Encode(util.BuildResponse(http.StatusNotImplemented, nil, "not support the others Method"))
		return
	}
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("server internal error", zap.String("req_body", string(reqBody)), zap.String("error", err.Error()))
		json.NewEncoder(w).Encode(util.BuildResponse(http.StatusInternalServerError, nil, "server internal error"))
		return
	}
	type Body struct {
		Message string `json:"message"`
	}
	var body Body
	if err := json.Unmarshal(reqBody, &body); err != nil {
		h.logger.Error("unmarshal failed", zap.String("req_body", string(reqBody)), zap.Error(err))
		if err := json.NewEncoder(w).Encode(util.BuildResponse(http.StatusInternalServerError, nil, "server internal error")); err != nil {
			h.logger.Error("send response failed", zap.Error(err))
		}
		return
	}
	h.logger.Info(" receiving message from the publishing client", zap.String("content", body.Message))
	// broadcast message
	timestamp := time.Now().Unix()
	if err := h.service.BroadcastMessage(timestamp, body.Message); err != nil {
		h.logger.Error("broadcast to queue failed", zap.Error(err))
		json.NewEncoder(w).Encode(util.BuildResponse(http.StatusInternalServerError, nil, "server internal error"))
		return
	}
	h.logger.Debug("Broadcast message success", zap.String("message", body.Message), zap.Int64("timestamp", timestamp))
	if err := json.NewEncoder(w).Encode(util.BuildResponse(http.StatusOK, nil, "success")); err != nil {
		h.logger.Error("send response failed", zap.Error(err))
	}
	return
}

func (h *Handler) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// init connection
	conn, err := h.upgrade.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("init connection socket failed", zap.String("error", err.Error()))
		return
	}
	h.logger.Debug("init connection socket success")

	for {
		// receive message from websocket client
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			h.logger.Error("read message failed", zap.String("error", err.Error()), zap.Int("message_type", messageType))
			return
		}
		h.logger.Info("receiving message via websocket", zap.String("content", string(p)))
		if err := conn.WriteMessage(messageType, p); err != nil {
			h.logger.Error("write message failed", zap.Error(err), zap.Int("message_type", messageType))
			return
		}
		// read message from queue
		for  {
			select {
			case <- *h.messageQueue:
				var fromPublisher []byte
				fromPublisher = <- *h.messageQueue
				if err := conn.WriteMessage(messageType, fromPublisher); err != nil {
					h.logger.Error("write message failed", zap.Error(err), zap.Int("message_type", messageType), zap.String("content", string(fromPublisher)))
					return
				}
			default:
				h.logger.Debug("Queue empty")
			}
		}
	}

}
