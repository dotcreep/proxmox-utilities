package utils

import (
	"encoding/json"
	"net/http"
)

// Json is a standard response wrapper.
//
//	@description	This struct represents the JSON response structure.
//	@Success		field	determines	the	success	status	of	the	API	call.
//	@Result			is the data payload.
//	@Message		provides additional information.
//	@Status			is the HTTP status code.
//	@Error			contains error details, if any.
type Response struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result"`
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Error   interface{} `json:"error"`
}

func (j *Response) newJson(success bool, result interface{}, message string, status int, error interface{}) *Response {
	return &Response{success, result, message, status, error}
}

func (j *Response) Json(success bool, w http.ResponseWriter, result interface{}, message string, status int, error interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(j.newJson(success, result, message, status, error))
}
