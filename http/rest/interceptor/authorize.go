package interceptor

import (
	"blacheapi/primer/gql"
	"blacheapi/primer/primitive"
	"blacheapi/primer/typing"
	"blacheapi/repository"
	"blacheapi/services/redis"
	"context"
)

// Authorize handles request authorization
// Authorize returns a session if the request is authorized or a gql error if it is not
func Authorize(ctx context.Context, permissions ...primitive.Array) (*redis.Session, interface{}, error) {

	// get session from context
	session := ctx.Value(typing.CtxSessionKey{})

	// if session is nil, return error
	if session == nil {
		return nil, nil, gql.MakeError("Please login to continue", 401)
	}

	auth := redis.Session{}

	// assert session as user

	u, ok := session.(repository.User)
	if ok {
		// if permissions are provided, check if the user has the permission

		auth = redis.Session{
			Email: u.Email,
			OrgID: u.OrgID,
			ID:    u.ID,
			Role:  u.Role,
		}

		return &auth, ctx.Value(typing.CtxSessionKey{}), nil
	}

	return nil, nil, gql.MakeError("You are not authorized to perform this action", 403)
}
