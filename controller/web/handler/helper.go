package handler

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	StatusCode int    `json:"error_code"`
	Message    string `json:"error_message"`
}

func WriteSuccessResponse(w http.ResponseWriter, statusCode int, data interface{}, headMap map[string]string) {
	w.Header().Add("Content-Type", "application/json")
	if len(headMap) > 0 {
		for key, val := range headMap {
			w.Header().Add(key, val)
		}
	}
	w.WriteHeader(statusCode)
	jsonData, _ := json.Marshal(data)
	w.Write(jsonData)
}

func WriteFailResponse(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errorMsg := Error{
		Message:    err.Error(),
		StatusCode: statusCode,
	}
	jsonData, _ := json.Marshal(errorMsg)

	w.Write(jsonData)
}
