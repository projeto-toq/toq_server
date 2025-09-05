package derrors

import (
	"errors"
)

// KindError is an error that carries a Kind and optional details/public message.
type KindError interface {
	error
	Kind() Kind
	Details() any
	PublicMessage() string
	Unwrap() error
}

// E is a concrete implementation of KindError.
type E struct {
	kind    Kind
	msg     string
	public  string
	details any
	cause   error
}

func (e *E) Error() string {
	if e.msg != "" {
		return e.msg
	}
	if e.public != "" {
		return e.public
	}
	return "error"
}
func (e *E) Kind() Kind   { return e.kind }
func (e *E) Details() any { return e.details }
func (e *E) PublicMessage() string {
	if e.public != "" {
		return e.public
	}
	return e.msg
}
func (e *E) Unwrap() error { return e.cause }

// New creates a new KindError.
func New(k Kind, msg string, opts ...Option) *E {
	e := &E{kind: k, msg: msg}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// Wrap wraps a cause with a new KindError.
func Wrap(cause error, k Kind, msg string, opts ...Option) *E {
	e := &E{kind: k, msg: msg, cause: cause}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// AsKind tries to extract KindError and return its kind; returns false if not KindError.
func AsKind(err error) (Kind, bool) {
	var ke KindError
	if errors.As(err, &ke) {
		return ke.Kind(), true
	}
	return KindInfra, false
}

// Option configures an error instance.
type Option func(*E)

// WithPublicMessage sets a safe public message for clients.
func WithPublicMessage(public string) Option { return func(e *E) { e.public = public } }

// WithDetails attaches optional details (e.g., field validation information).
func WithDetails(details any) Option { return func(e *E) { e.details = details } }

// Convenience constructors
func Infra(msg string, cause error, opts ...Option) *E { return Wrap(cause, KindInfra, msg, opts...) }
func Auth(msg string, opts ...Option) *E               { return New(KindAuth, msg, opts...) }
func Forbidden(msg string, opts ...Option) *E          { return New(KindForbidden, msg, opts...) }
func Conflict(msg string, opts ...Option) *E           { return New(KindConflict, msg, opts...) }
func Gone(msg string, opts ...Option) *E               { return New(KindGone, msg, opts...) }
func Unprocessable(msg string, details any, opts ...Option) *E {
	return New(KindUnprocessable, msg, append(opts, WithDetails(details))...)
}
func NotFound(msg string, opts ...Option) *E   { return New(KindNotFound, msg, opts...) }
func BadRequest(msg string, opts ...Option) *E { return New(KindBadRequest, msg, opts...) }
func Validation(msg string, details any, opts ...Option) *E {
	return New(KindValidation, msg, append(opts, WithDetails(details))...)
}
func Locked(msg string, opts ...Option) *E          { return New(KindLocked, msg, opts...) }
func TooManyRequests(msg string, opts ...Option) *E { return New(KindTooManyRequests, msg, opts...) }
