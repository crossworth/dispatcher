package dispatcher

import (
	"testing"
	"time"
)

type customEvent struct {
	ID   int
	Name string
}

func ptr[T any](input T) *T {
	return &input
}

func TestDispatcher_Example(t *testing.T) {
	dp := NewDispatcher()
	dp.Register(NewRegistration(func(input int) {
		println(input)
	}))

	dp.Register(NewRegistration(func(input *int) {
		println(input)
	}))

	dp.Register(NewRegistration(func(input string) {
		println(input)
	}))

	dp.Register(NewRegistration(func(input string) {
		println(input)
	}))

	dp.Register(NewRegistration(func(input customEvent) {
		println(input.ID)
		println(input.Name)
	}))

	dp.Dispatch(10)
	dp.Dispatch(ptr(100))
	dp.Dispatch("abc")
	dp.Dispatch(customEvent{
		ID:   10,
		Name: "Test",
	})
	time.Sleep(300 * time.Millisecond)
}
