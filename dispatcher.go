package dispatcher

import (
	"fmt"
	"sync"
)

// Dispatcher is a generic dispatcher.
type Dispatcher[T any] struct {
	mu       sync.RWMutex
	handlers map[string][]any
}

// NewDispatcher creates a new *Dispatcher[any].
func NewDispatcher() *Dispatcher[any] {
	return &Dispatcher[any]{
		handlers: map[string][]any{},
	}
}

// Handler holds information about the handler and type
// to be used to determine the correct dispatcher.
type Handler interface {
	handlerType() any
	handlerFunc() func(input any)
}

type handler struct {
	typ     any
	handler func(input any)
}

func (r handler) handlerType() any {
	return r.typ
}

func (r handler) handlerFunc() func(input any) {
	return r.handler
}

// HandlerFunc creates a new dispatcher Handler.
func HandlerFunc[T any](handlerFunc func(input T)) Handler {
	return handler{
		typ: *new(T),
		handler: func(input any) {
			handlerFunc(input.(T))
		},
	}
}

// Register registers the provided handler
// on the dynamic dispatcher.
func (d *Dispatcher[T]) Register(handler Handler) {
	d.mu.Lock()
	defer d.mu.Unlock()

	typ := fmt.Sprintf("%T", handler.handlerType())
	d.handlers[typ] = append(d.handlers[typ], handler.handlerFunc())
}

// Dispatch dispatches the provided event
// for all the handlers for the given type.
func (d *Dispatcher[T]) Dispatch(event T) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	typ := fmt.Sprintf("%T", event)
	handlers, found := d.handlers[typ]
	if !found {
		return
	}

	for _, handler := range handlers {
		handler.(func(input T))(event)
	}
}
