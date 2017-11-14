package drops

import (
	"fmt"
	"testing"
)

func TestEvents(t *testing.T) {
	emitter := NewEventEmitter()

	const (
		Event1 Event = iota + 1
		Event2
		Event3
	)

	count := 0
	// setup handler
	handler := func(id int) func(Event, interface{}, error) {
		hf := func(e Event, data interface{}, err error) {
			switch e {
			case Event1:
				fmt.Printf("handler%d got event1: %v\n", id, e)
			case Event2:
				fmt.Printf("handler%d got event2: %v\n", id, e)
			case Event3:
				fmt.Printf("handler%d got event3: %v\n", id, e)

			}
			count++
		}
		return hf
	}

	handlersCount := 23
	for i := 0; i < handlersCount; i++ {
		emitter.Handler(EventHandlerFunc(handler(i)))
	}

	emitter.Dispatch(Event1, nil, nil)

	emitter.Close()
	if count != handlersCount {
		t.Fatalf("Not all handlers processed. expected %d got %d", handlersCount, count)
	}
}

func TestEventsMultiple(t *testing.T) {
	emitter := NewEventEmitter()

	const (
		Event1 Event = iota + 1
		Event2
		Event3
	)

	count := 0
	// setup handler
	handler := func(id int) func(Event, interface{}, error) {
		hf := func(e Event, data interface{}, err error) {
			switch e {
			case Event1:
				fmt.Printf("handler%d got event1: %v\n", id, e)
			case Event2:
				fmt.Printf("handler%d got event2: %v\n", id, e)
			case Event3:
				fmt.Printf("handler%d got event3: %v\n", id, e)

			}
			count++
		}
		return hf
	}

	handlersCount := 23
	for i := 0; i < handlersCount; i++ {
		emitter.Handler(EventHandlerFunc(handler(i)))
	}

	emitter.Dispatch(Event1, nil, nil)
	emitter.Dispatch(Event2, nil, nil)
	emitter.Dispatch(Event3, nil, nil)

	emitter.Close()
	fmt.Println("Close")
	if count != handlersCount*3 {
		t.Fatalf("Not all handlers processed. expected %d got %d", handlersCount*3, count)
	}
}

func TestEventsNoHandlers(t *testing.T) {
	emitter := NewEventEmitter()

	emitter.Close()
}

func TestEventsOneHandler(t *testing.T) {
	emitter := NewEventEmitter()

	const (
		Event1 Event = iota + 1
		Event2
		Event3
		Event4
		Event5
		Event6
	)

	// setup handler
	hf := func(e Event, data interface{}, err error) {
		fmt.Printf("handler got event: %064b\n", e)
	}

	emitter.Handler(EventHandlerFunc(hf))
	emitter.Dispatch(Event6, nil, nil)
	emitter.Close()
}

func TestEventsStopDispatch(t *testing.T) {
	emitter := NewEventEmitter()

	const (
		Event1 Event = iota + 1
		Event2
		Event3
		Event4
		Event5
		Event6
	)

	// setup handler
	hf := func(e Event, data interface{}, err error) {
		fmt.Printf("handler got event: %064b\n", e)
	}

	emitter.Handler(EventHandlerFunc(hf))
	emitter.Dispatch(Event6, nil, nil)
	emitter.Close()
	emitter.Dispatch(Event3, nil, nil)
}
