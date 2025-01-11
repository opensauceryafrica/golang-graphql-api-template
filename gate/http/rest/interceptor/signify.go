package interceptor

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"cendit.io/garage/config"
	"cendit.io/garage/primer/gql"
	"cendit.io/garage/zksnark"
	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// Signify handles payload signature verification
func Signify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// get the signature
		signature := r.Header.Get("X-Request-Witness")

		// read the body
		body, _ := io.ReadAll(r.Body)

		// restore the body
		r.Body = io.NopCloser(strings.NewReader(string(body)))

		// if the signature is not present, return an error
		if signature == "" && (!(r.Method == http.MethodPost && strings.Contains(string(body), "IntrospectionQuery")) && !(r.Method == http.MethodGet && r.URL.Query().Get("witness") == "playground") && !(r.Method == http.MethodGet && r.URL.Path == "/favicon.ico")) {

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"errors": []gqlerror.Error{
					{
						Path:    graphql.GetPath(r.Context()),
						Message: "Witness is missing!",
						Extensions: map[string]interface{}{
							"code":   gql.CodeForStatus[http.StatusBadRequest],
							"status": http.StatusBadRequest,
						},
					},
				},
			})
			return
		}

		// if the signature is present, verify it
		// bypass the signature verification if the signature matches a specific value
		if signature != "" && signature != config.Env.AppSecret {

			// verify the signature
			// basically take the body and try to verify it with the signature
			if !zksnark.Witness(signature, string(body)) {

				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"errors": []gqlerror.Error{
						{
							Path:    graphql.GetPath(r.Context()),
							Message: "Witness is not genuine!",
							Extensions: map[string]interface{}{
								"code":   gql.CodeForStatus[http.StatusBadRequest],
								"status": http.StatusBadRequest,
							},
						},
					},
				})
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
