package handler

import (
	"encoding/json"
	"net/http"

	"gopkg.in/macaron.v1"
)

// IndexMetaV1Handler now only helps to know if the server is alive.
func IndexMetaV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": "Update Server Backend REST API Service"})
	return http.StatusOK, result
}
