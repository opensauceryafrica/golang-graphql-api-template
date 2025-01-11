package interceptor

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"cendit.io/garage/function"
	"cendit.io/garage/logger"
	"cendit.io/garage/primer/constant"
	"cendit.io/garage/primer/typing"
	"cendit.io/garage/redis"
)

type WrappedResponseWriter struct {
	http.ResponseWriter
	status int
	body   []byte
}

func (w *WrappedResponseWriter) WriteHeader(status int) {
	w.status = status
}

func (w *WrappedResponseWriter) Write(b []byte) (int, error) {
	w.body = b
	return len(b), nil
}

// Intruder handles request/response intrusion
func Intruder(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// wrap the response writer to hijack write control
		ww := &WrappedResponseWriter{ResponseWriter: w}

		next.ServeHTTP(ww, r)

		// get context from request
		ctx := r.Context()

		// get status code from context
		instrusion, ok := ctx.Value(typing.CtxIntrusionKey{}).(*typing.Intrusion)

		// get auth header
		authHeader := r.Header.Get("Authorization")
		token := ""

		if authHeader != "" {

			if len(strings.Split(authHeader, "Bearer ")) > 1 {
				// get the token from the auth header
				token = strings.Split(authHeader, "Bearer ")[1]
			}
		}

		// get user from context
		session := instrusion.Session

		if ok && session.ID != "" {

			// generate a new token and set it in the response header
			newToken := function.GenerateUUID()
			w.Header().Set("X-Authorization", newToken)

			// set the token in the redis (expire in 5 minutes)
			if err := redis.Ral.Set(fmt.Sprintf("%s-%s", constant.UserRedisKey, newToken), function.Bite(session), 5*time.Minute); err != nil {
				logger.GetLogger().Error(fmt.Sprintf(`interceptor.Authenticate :: [%v] :: failed to set token in redis: %s`, ctx.Value(typing.CtxTraceKey{}), err.Error()))
			} else {
				// remove the old token from the redis
				if err := redis.Ral.Del(fmt.Sprintf("%s-%s", constant.UserRedisKey, token)); err != nil {
					logger.GetLogger().Error(fmt.Sprintf(`interceptor.Authenticate :: [%v] :: failed to delete old token from redis: %s`, ctx.Value(typing.CtxTraceKey{}), err.Error()))
				}
			}

		}

		if ok && instrusion.Status != 0 {
			ww.ResponseWriter.WriteHeader(instrusion.Status)

			// write the response body
			ww.ResponseWriter.Write(ww.body)
		} else {
			if ww.status != 0 {
				ww.ResponseWriter.WriteHeader(ww.status)
			}
			ww.ResponseWriter.Write(ww.body)
		}
	})
}
