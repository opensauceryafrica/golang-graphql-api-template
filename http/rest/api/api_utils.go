package api

import (
	"blacheapi/primer/constant"
	"net/http"
)

func parseStatusCode(status string) int {
	switch status {
	case constant.Error:
		return http.StatusInternalServerError
	case constant.Success:
		return http.StatusOK
	default:
		return http.StatusInternalServerError
	}
}
