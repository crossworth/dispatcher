package dispatcher

import (
	"fmt"
)

// Dispatcher is a generic dispatcher.
type Dispatcher[T any] struct {
	handlers map[string][]any
}

// NewDispatcher creates a new *Dispatcher[any].
func NewDispatcher() *Dispatcher[any] {
	return &Dispatcher[any]{
		handlers: map[string][]any{},
	}
}

type registration struct {
	handler func(input any)
	typ     any
}

// NewRegistration creates a new dispatcher registration.
func NewRegistration[T any](handler func(input T)) registration {
	return registration{
		handler: func(input any) {
			handler(input.(T))
		},
		typ: *new(T),
	}
}

// Register registers the provided registration
// on the dynamic dispatcher.
func (d Dispatcher[T]) Register(registration registration) {
	typ := fmt.Sprintf("%T", registration.typ)
	d.handlers[typ] = append(d.handlers[typ], registration.handler)
}

// Dispatch dispatches the provided event
// for all the registrations.
func (d Dispatcher[T]) Dispatch(event T) {
	typ := fmt.Sprintf("%T", event)
	handlers, found := d.handlers[typ]
	if !found {
		return
	}

	for _, handler := range handlers {
		go handler.(func(input T))(event)
	}
}
