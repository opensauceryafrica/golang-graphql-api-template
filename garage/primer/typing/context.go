package typing

// CtxTraceKey is the context key for tracing.
type CtxTraceKey struct{}

// CtxSessionKey is the context key for session.
type CtxSessionKey struct{}

// CtxErrorKey is the context key for error.
type CtxErrorKey struct{}

// CtxStatusKey is the context key for status.
type CtxStatusKey struct{}

// CtxIntrusionKey is the context key for intrusion.
type CtxIntrusionKey struct{}

type Intrusion struct {
	Session Session
	Status  int
}
