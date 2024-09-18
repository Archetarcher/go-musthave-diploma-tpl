package handlers

import (
	"encoding/json"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

func sendResponse(enc *json.Encoder, data interface{}, code int, w http.ResponseWriter) {
	w.WriteHeader(code)

	if err := enc.Encode(data); err != nil {
		logger.Log.Info("error", zap.Any("err", err))

	}

}
