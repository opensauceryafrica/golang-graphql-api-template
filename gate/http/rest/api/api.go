package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"cendit.io/garage/config"

	"cendit.io/garage/logger"
	"cendit.io/garage/primer/constant"
	"cendit.io/gate/http/graphql"
	"cendit.io/gate/http/graphql/resolver"
	"cendit.io/gate/http/rest/interceptor"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/cors"
)

type API struct {
	Server    *http.Server
	Variables *config.Variable
	Deps      *config.Dependencies
}

// SetupServerHandler ...
func (a *API) SetupServerHandler() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(60 * time.Second))
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodPatch},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", config.HeaderRequestID, config.HeaderRequestSource},
		ExposedHeaders:   []string{"Cendit"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(interceptor.RequestTracing)

	mux.Use(interceptor.Signify)

	mux.Use(interceptor.Authenticate)

	mux.Mount("/health", HealthRoute())

	srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: &resolver.Resolver{}}))

	mux.Handle("/graph/entrypoint", interceptor.Intruder(srv))

	// only enable introspection and playground in dev mode
	if os.Getenv("ENVIRONMENT") != "production" {
		mux.Handle("/dev/graphql", playground.ApolloSandboxHandler("Cendit GraphQL playground", "/cendit/graph/entrypoint"))
		mux.Handle("/graphql", playground.ApolloSandboxHandler("Cendit GraphQL playground", "/graph/entrypoint"))
	}

	return mux
}

// Serve starts the service api
func (a *API) Serve() error {
	a.Server = &http.Server{
		Addr:           fmt.Sprintf(":%d", a.Variables.Port),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		Handler:        a.SetupServerHandler(),
		MaxHeaderBytes: 1 << 20,
	}

	logger.GetLogger().Info("[API]: Starting ...")
	return a.Server.ListenAndServe()
}

// Shutdown stops the service api
func (a *API) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return a.Server.Shutdown(ctx)
}

// Handler wraps our http handlers so we can execute some actions before and after a handler is run
type Handler func(w http.ResponseWriter, r *http.Request) *ServerResponse

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := h(w, r)

	var responseBytes []byte
	var marshalErr error

	responseBytes, marshalErr = json.Marshal(response)
	if marshalErr != nil {
		writeErrorResponse(w, constant.Error, "Error decoding resposne body")
	}
	writeJSONResponse(w, response.StatusCode, responseBytes)
}
