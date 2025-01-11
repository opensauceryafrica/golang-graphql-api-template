package api

import (
	"net/http"

	"cendit.io/garage/primer/constant"
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
