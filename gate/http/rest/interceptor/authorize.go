package interceptor

import (
	"context"

	"cendit.io/auth/schema"
	"cendit.io/garage/primer/typing"
	"cendit.io/garage/primitive"
	"cendit.io/gate/http/graphql/exception"
)

// Authorize handles request authorization
// Authorize returns a session if the request is authorized or a gql error if it is not
func Authorize(ctx context.Context, permissions ...primitive.Array) (*typing.Session, interface{}, error) {

	// get session from context
	session := ctx.Value(typing.CtxSessionKey{})

	// if session is nil, return error
	if session == nil {
		return nil, nil, exception.MakeError("Please login to continue", 401)
	}

	auth := typing.Session{}

	// assert session as user

	u, ok := session.(*schema.User)
	if ok {
		// if permissions are provided, check if the user has the permission

		auth = typing.Session{
			Email: u.Email,
			ID:    u.ID,
			Role:  u.Role,
		}

		return &auth, *ctx.Value(typing.CtxSessionKey{}).(*schema.User), nil
	}

	return nil, nil, exception.MakeError("You are not authorized to perform this action", 403)
}
