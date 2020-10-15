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
	config  *config.Config
	upgrade websocket.Upgrader
	logger  log.Logger
	service service.Service
	utilInstance	util.Util
}

func NewHandler(config *config.Config, logger log.Logger, upgrade websocket.Upgrader, utilInstance util.Util) *Handler {
	serviceInstance := service.NewSendMessageService(config)
	return &Handler{
		config:  config,
		upgrade: upgrade,
		logger:  logger,
		service: serviceInstance,
		utilInstance: utilInstance,
	}
}

/*
impl health check for server i think every apis server need it,
*/
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	resp := util.NewUtil().BuildResponse(http.StatusOK, nil, "Alive")
	json.NewEncoder(w).Encode(resp)
	return
}

func (h *Handler) BroadCastMessageHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		h.logger.Warn("not support the others Method")
		json.NewEncoder(w).Encode(h.utilInstance.BuildResponse(http.StatusNotImplemented, nil, "not support the others Method"))
		return
	}
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("server internal error", zap.String("req_body", string(reqBody)), zap.String("error", err.Error()))
		json.NewEncoder(w).Encode(h.utilInstance.BuildResponse(http.StatusInternalServerError, nil, "server internal error"))
		return
	}
	type Body struct {
		Message string `json:"message"`
	}
	var body Body
	if err := json.Unmarshal(reqBody, &body); err != nil {
		h.logger.Error("unmarshal failed", zap.String("req_body", string(reqBody)), zap.Error(err))
		if err := json.NewEncoder(w).Encode(h.utilInstance.BuildResponse(http.StatusInternalServerError, nil, "server internal error")); err != nil {
			h.logger.Error("send response failed", zap.Error(err))
		}
		return
	}
	h.logger.Info(" receiving message from the publishing client", zap.String("content", body.Message))
	// broadcast message
	timestamp := time.Now().Unix()
	if err := h.service.BroadcastMessage(timestamp, body.Message); err != nil {
		h.logger.Error("broadcast to queue failed", zap.Error(err))
		_ = json.NewEncoder(w).Encode(h.utilInstance.BuildResponse(http.StatusInternalServerError, nil, "server internal error"))
		return
	}
	h.logger.Debug("Broadcast message success", zap.String("message", body.Message), zap.Int64("timestamp", timestamp))
	if err := json.NewEncoder(w).Encode(h.utilInstance.BuildResponse(http.StatusOK, nil, "success")); err != nil {
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
		// read message from broadcast
		var err error
		var raw []byte
		var previousData []byte
		for {
			raw, err = ioutil.ReadFile(h.config.TemporaryFile)
			if err != nil {
				h.logger.Error("write message failed", zap.Error(err), zap.Int("message_type", 1))
				break
			}
			if string(raw) == string(previousData) {
				continue
			}
			previousData = raw
			if len(raw) != 0 {
				if err = conn.WriteMessage(websocket.TextMessage, raw); err != nil {
					h.logger.Error("write message failed", zap.Error(err), zap.Int("message_type", websocket.TextMessage), zap.String("content", string(raw)))
					break
				}

				// receive message from websocket client
				var messageType int
				var p []byte
				messageType, p, err = conn.ReadMessage()
				if err != nil {
					h.logger.Debug("read message failed", zap.String("error", err.Error()), zap.Int("message_type", messageType))
					break
				}
				h.logger.Info("receiving message via websocket", zap.String("content", string(p)))
			}

		}

	}
}
