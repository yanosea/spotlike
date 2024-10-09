package contextproxy

import (
	"context"
)

// Context is an interface for context.
type Context interface {
	Background() *ContextInstance
}

// ContextProxy is a struct that implements Context.
type ContextProxy struct{}

// New is a constructor for ContextProxy.
func New() *ContextProxy {
	return &ContextProxy{}
}

// Background is a proxy for context.Background().
func (*ContextProxy) Background() *ContextInstance {
	return &ContextInstance{FieldContext: context.Background()}
}
