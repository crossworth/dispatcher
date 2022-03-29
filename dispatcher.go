package dispatcher

import (
	"fmt"
	"strings"
	"sync"
)

type handlerRecord struct {
	handlerFunc     any
	executionPolicy string
}

// Dispatcher is a generic dispatcher.
type Dispatcher[T any] struct {
	mu       sync.RWMutex
	handlers map[string]handlerRecord
	counter  int64
}

// NewDispatcher creates a new *Dispatcher[any].
func NewDispatcher() *Dispatcher[any] {
	return &Dispatcher[any]{
		handlers: map[string]handlerRecord{},
	}
}

// Handler holds information about the handler and type
// to be used to determine the correct dispatcher.
type Handler interface {
	handlerType() any
	handlerFunc() func(input any)
	executionPolicy() string
}

type handler struct {
	typ     any
	handler func(input any)
	policy  string
}

func (r handler) handlerType() any {
	return r.typ
}

func (r handler) handlerFunc() func(input any) {
	return r.handler
}

func (r handler) executionPolicy() string {
	return r.policy
}

// HandlerFunc creates a new dispatcher Handler.
func HandlerFunc[T any](handlerFunc func(input T)) Handler {
	return handler{
		typ: *new(T),
		handler: func(input any) {
			handlerFunc(input.(T))
		},
		policy: "always",
	}
}

// HandlerFuncOnce creates a new dispatcher Handler that only executes once.
func HandlerFuncOnce[T any](handlerFunc func(input T)) Handler {
	return handler{
		typ: *new(T),
		handler: func(input any) {
			handlerFunc(input.(T))
		},
		policy: "once",
	}
}

// Register registers the provided handler
// on the dynamic dispatcher.
func (d *Dispatcher[T]) Register(handler Handler) string {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.counter++

	handlerID := fmt.Sprintf("%T_%d", handler.handlerType(), d.counter)

	d.handlers[handlerID] = handlerRecord{
		handlerFunc:     handler.handlerFunc(),
		executionPolicy: handler.executionPolicy(),
	}

	return handlerID
}

// Unregister unregister the provided handlerID.
func (d *Dispatcher[T]) Unregister(handlerID string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.handlers, handlerID)
}

// Dispatch dispatches the provided event
// for all the handlers for the given type.
func (d *Dispatcher[T]) Dispatch(event T) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	typ := fmt.Sprintf("%T", event)
	for key, handler := range d.handlers {
		keyType := extractType(key)

		if keyType != typ {
			continue
		}

		handler.handlerFunc.(func(input T))(event)

		if handler.executionPolicy == "once" {
			delete(d.handlers, key)
		}
	}
}

func extractType(handlerID string) string {
	parts := strings.Split(handlerID, "_")
	return parts[0]
}
