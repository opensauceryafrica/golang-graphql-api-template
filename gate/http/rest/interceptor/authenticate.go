package interceptor

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"cendit.io/auth/repository"
	"cendit.io/garage/function"
	"cendit.io/garage/logger"
	"cendit.io/garage/primer/constant"
	"cendit.io/garage/primer/typing"
	"cendit.io/garage/redis"
	"cendit.io/garage/xiao"
)

// Authenticate handles request authentication
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// get auth header
		authHeader := r.Header.Get("Authorization")

		logger.GetLogger().Debug(fmt.Sprintf(`interceptor.Authenticate :: [%v] :: received call to authenticate token: %s`, r.Context().Value(typing.CtxTraceKey{}), authHeader))

		// in graphql, we only authenticate the request if the auth header is present
		// if the auth header is not present, we assume that the request is not authenticated
		if authHeader != "" {

			if len(strings.Split(authHeader, "Bearer ")) > 1 {
				// get the token from the auth header
				token := strings.Split(authHeader, "Bearer ")[1]

				if reader, err := redis.Ral.Get(fmt.Sprintf("%s-%s", constant.UserRedisKey, token)); err == nil {

					var session typing.Session

					if err := function.Load(string(reader), &session); err == nil {

						u, err := repository.User().FindByMap(context.Background(), xiao.SQLMaps{
							WMaps: []xiao.SQLMap{
								{
									Map: map[string]interface{}{
										"id": session.ID,
									},
									JoinOperator:       xiao.And,
									ComparisonOperator: xiao.Equal,
								},
							},
						}, true)
						if err == nil {

							logger.GetLogger().Debug(fmt.Sprintf(`interceptor.Authenticate :: [%v] :: successfully authenticated token: %s :: found user id: %s`, r.Context().Value(typing.CtxTraceKey{}), authHeader, u.ID))

							// set the user in the request context
							r = r.WithContext(context.WithValue(r.Context(), typing.CtxSessionKey{}, u))

							intrusion := r.Context().Value(typing.CtxIntrusionKey{}).(*typing.Intrusion)
							intrusion.Session = typing.Session{
								ID:    u.ID,
								Email: u.Email,
								Role:  u.Role,
							}

						}
					}

				}
			}
		}
		next.ServeHTTP(w, r)
	})
}
