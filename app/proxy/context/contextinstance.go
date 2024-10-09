package contextproxy

import (
	"context"
)

// ContextInstanceInterface is an interface for context.Context.
type ContextInstanceInterface interface {
	Value(key interface{}) interface{}
}

// ContextInstance is a struct that implements ContextInstanceInterface.
type ContextInstance struct {
	FieldContext context.Context
}

func (c *ContextInstance) Value(key interface{}) interface{} {
	return c.FieldContext.Value(key)
}
