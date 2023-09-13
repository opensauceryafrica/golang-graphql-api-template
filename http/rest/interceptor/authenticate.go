package interceptor

import (
	"blacheapi/primer/constant"
	"blacheapi/primer/enum"
	"blacheapi/primer/typing"
	"blacheapi/repository"
	"blacheapi/services/redis"
	"context"
	"fmt"
	"net/http"
	"strings"
)

// Authenticate handles request authentication
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// get auth header
		authHeader := r.Header.Get("Authorization")

		// in graphql, we only authenticate the request if the auth header is present
		// if the auth header is not present, we assume that the request is not authenticated
		if authHeader != "" {

			if len(strings.Split(authHeader, "Bearer ")) > 1 {
				// get the token from the auth header
				token := strings.Split(authHeader, "Bearer ")[1]

				if session, err := redis.Ral.GetSession(fmt.Sprintf("%s-%s", constant.UserRedisKey, token)); err == nil {
					// try as user

					u := repository.User{}

					if err := u.FByMap(typing.SQLMaps{
						WMaps: []typing.SQLMap{
							{
								Map: map[string]interface{}{
									"id":     session.ID,
									"org_id": session.OrgID,
								},
								JoinOperator:       enum.And,
								ComparisonOperator: enum.Equal,
							},
						},
					}, true); err == nil {
						u.Role = enum.User
						// set the org id in the request context
						r = r.WithContext(context.WithValue(r.Context(), typing.CtxSessionKey{}, u))
					}
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}
