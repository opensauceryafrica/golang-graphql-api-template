package api

import (
	"blacheapi/logger"
	"blacheapi/monitor"
	"blacheapi/primer/constant"
	"context"
	"encoding/json"
	"net/http"

	"github.com/getsentry/sentry-go"
)

// ServerResponse contains information required
// to generate the server response to a request
type ServerResponse struct {
	Err         error           `json:"-"`
	Message     string          `json:"message"`
	Status      string          `json:"status"`
	StatusCode  int             `json:"status_code"`
	Context     context.Context `json:"-"`
	ContentType string          `json:"-"`
	Payload     interface{}     `json:"payload"`
}

type ErrorResponse struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

// writeJSONResponse writes the server handler response to the response writer
func writeJSONResponse(w http.ResponseWriter, statusCode int, content []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write(content); err != nil {
		logger.GetLogger().Sugar().Errorf("Unable to write json resposne: %v", err)
	}
}

func writeErrorResponse(w http.ResponseWriter, status string, errString string) {
	r := RespondWithError(nil, errString, status)
	errorResponse, _ := json.Marshal(r)
	writeJSONResponse(w, r.StatusCode, errorResponse)
}

// RespondWithError parses an error to the ServerResponse type, and captures a sentry exception
// if status is error
func RespondWithError(err error, message, status string) *ServerResponse {

	if status == constant.Error {
		monitor.SendScopeLocalizedError(err, nil, "", 0, sentry.LevelError)
	}

	return &ServerResponse{
		Err:        err,
		Status:     status,
		Message:    message,
		StatusCode: parseStatusCode(status),
	}
}
