## Generic dynamic dispatcher

```go
package main

import (
	"time"

	"github.com/crossworth/dispatcher"
)

type customEvent struct {
	ID   int
	Name string
}

func ptr[T any](input T) *T {
	return &input
}

func main() {
	dp := dispatcher.NewDispatcher()

	dp.Register(dispatcher.HandlerFunc(func(input int) {
		println(input)
	}))

	dp.Register(dispatcher.HandlerFunc(func(input *int) {
		println(input)
	}))

	dp.Register(dispatcher.HandlerFunc(func(input string) {
		println(input)
	}))

	dp.Register(dispatcher.HandlerFunc(func(input string) {
		println(input)
	}))

	dp.Register(dispatcher.HandlerFunc(func(input customEvent) {
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

// Result:
// 10
// abc
// 0xc000018400
// abc
// 10
// Test
```