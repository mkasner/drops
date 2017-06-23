package drops

import (
	"fmt"
	"testing"
	"time"
)

func TestEvents(t *testing.T) {
	emitter := NewEventEmitter()

	const (
		Event1 Event = iota + 1
		Event2
		Event3
	)

	// setup handler
	handler := func(id int, count chan int) func(Event, interface{}, error) {
		hf := func(e Event, data interface{}, err error) {
			switch e {
			case Event1:
				fmt.Printf("handler%d got event1: %v\n", id, e)
			case Event2:
				fmt.Printf("handler%d got event2: %v\n", id, e)
			case Event3:
				fmt.Printf("handler%d got event3: %v\n", id, e)

			}
			count <- 1
		}
		return hf
	}

	handlersCount := 23
	countHandles := make(chan int, 100)
	for i := 0; i < handlersCount; i++ {
		emitter.Handler(EventHandlerFunc(handler(i, countHandles)))
	}

	count := 0
	go func() {
		for c := range countHandles {
			count = count + c
		}
	}()
	time.Sleep(2 * time.Second)

	emitter.Dispatch(Event1, nil, nil)
	emitter.Dispatch(Event2, nil, nil)
	emitter.Dispatch(Event3, nil, nil)
	time.Sleep(1 * time.Second)

	emitter.Close()
	close(countHandles)
	if count != handlersCount*3 {
		t.Fatalf("Not all handlers processed. expected %d got %d", handlersCount*3, count)
	}
}
