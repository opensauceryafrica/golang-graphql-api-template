package interceptor

import (
	"net/http"
)

// RequestTracing handles request tracing
func RequestTracing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// @TODO: tracing logic
		// create a trace id and add it to the request context
		// log the first event of the trace to BigQuery
		// further event logs can use the trace id to link to the first event
		// so, should an audit be required in the future, we can trace the request
		// from the first event to the last event, first by the time of occurence
		// then further by the trace id
		// BigQuery provide analytical capabilities in standard SQL and it is
		// relatively cost effective to store data in it. Retrieval is a
		// tiny bit more expensive but that's not an ish as we don't expect to
		// be retrieving data from BigQuery often

		next.ServeHTTP(w, r)
	})
}
