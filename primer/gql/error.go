package gql

import (
	"blacheapi/http/graphql/model"
	"blacheapi/primer/enum"
	"blacheapi/primer/function"
	"context"

	"blacheapi/logger"
	"blacheapi/repository"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func MakeSubgraphError(message string, status int, errorlog ...string) *model.Error {
	// log activity
	makeErrorLog(message, errorlog...)

	return &model.Error{
		Message: message,
		Code:    CodeForStatus[status],
		Status:  status,
	}
}

type Error struct {
	Message string
	Status  int
}

func (e Error) Error() string {
	return e.Message
}

func MakeError(message string, status int) Error {
	return Error{
		Message: message,
		Status:  status,
	}
}

func makeErrorLog(message string, errorlog ...string) {
	// log activity
	if len(errorlog) > 1 {
		a := repository.Activity{}
		if err := function.Load(errorlog[1], &a); err != nil {
			logger.GetLogger().Debug(`ACTIVITY LOG :: LOAD ERROR :: ` + err.Error())
		}
		a.Status = enum.Failure
		a.Error = message
		if err := a.Create(); err != nil {
			logger.GetLogger().Debug(`ACTIVITY LOG :: SAVE ERROR :: ` + err.Error())
		}
	}

	if len(errorlog) > 0 {
		logger.GetLogger().Debug(`HALTED :: ` + strings.ReplaceAll(errorlog[0], "{message}", message))
	}
}

func MakeGraphQLError(ctx context.Context, message string, status int, errorlog ...string) error {
	// log activity
	makeErrorLog(message, errorlog...)

	return &gqlerror.Error{
		Path:    graphql.GetPath(ctx),
		Message: message,
		Extensions: map[string]interface{}{
			"code":   CodeForStatus[status],
			"status": status,
		},
	}
}
