package dispatcher

import (
	"sync"
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
	dp.Register(HandlerFunc(func(input int) {
		println(input)
	}))

	dp.Register(HandlerFunc(func(input *int) {
		println(input)
	}))

	dp.Register(HandlerFunc(func(input string) {
		println(input)
	}))

	dp.Register(HandlerFunc(func(input string) {
		println(input)
	}))

	dp.Register(HandlerFunc(func(input customEvent) {
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

func TestDispatcher_Dispatch(t *testing.T) {
	dp := NewDispatcher()

	calls := int32(0)
	callsMu := sync.Mutex{}

	dp.Register(HandlerFunc(func(input int) {
		callsMu.Lock()
		calls++
		callsMu.Unlock()

		if input != 10 {
			t.Fail()
		}
	}))
	dp.Dispatch(10)

	time.Sleep(100 * time.Millisecond)

	ensureCalls := func() {
		callsMu.Lock()
		defer callsMu.Unlock()
		if calls != 1 {
			t.Fail()
		}
	}

	ensureCalls()

	dp.Dispatch(50.0)

	time.Sleep(100 * time.Millisecond)
	ensureCalls()
}

func TestDispatcher_DispatchOnce(t *testing.T) {
	dp := NewDispatcher()

	calls := int32(0)
	callsMu := sync.Mutex{}

	dp.Register(HandlerFuncOnce(func(input int) {
		callsMu.Lock()
		calls++
		callsMu.Unlock()

		if input != 10 {
			t.Fail()
		}
	}))
	dp.Dispatch(10)
	dp.Dispatch(10)
	dp.Dispatch(10)
	dp.Dispatch(10)
	dp.Dispatch(10)

	time.Sleep(100 * time.Millisecond)

	callsMu.Lock()
	defer callsMu.Unlock()
	if calls != 1 {
		t.Fail()
	}
}

func TestDispatcher_Unregister(t *testing.T) {
	dp := NewDispatcher()
	calls := int32(0)
	callsMu := sync.Mutex{}

	handlerID := dp.Register(HandlerFuncOnce(func(input int) {
		callsMu.Lock()
		calls++
		callsMu.Unlock()

		if input != 10 {
			t.Fail()
		}
	}))

	time.Sleep(1 * time.Millisecond)

	dp.Dispatch(10)

	time.Sleep(5 * time.Millisecond)

	dp.Unregister(handlerID)

	time.Sleep(5 * time.Millisecond)

	dp.Dispatch(10)

	time.Sleep(5 * time.Millisecond)

	callsMu.Lock()
	defer callsMu.Unlock()
	if calls != 1 {
		t.Fatalf("got wrong number of calls, %d", calls)
	}
}
