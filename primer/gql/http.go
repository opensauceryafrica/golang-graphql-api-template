package gql

import "blacheapi/primer/constant"

// HTTPStatus
const (
	BadRequest          = "BAD_REQUEST"
	InternalServerError = "INTERNAL_SERVER_ERROR"
	NotFound            = "NOT_FOUND"
	Unauthorized        = "UNAUTHORIZED"
	Forbidden           = "FORBIDDEN"
	Unathenticated      = "UNAUTHENTICATED"
	OK                  = "OK"
)

// CodeForStatus
var CodeForStatus = map[int]string{
	constant.CodeBadRequest:   BadRequest,
	constant.CodeUnauthorized: Unauthorized,
	constant.CodeForbidden:    Forbidden,
	constant.CodeNotFound:     NotFound,
	constant.CodeISE:          InternalServerError,
	constant.CodeOK:           OK,
}
