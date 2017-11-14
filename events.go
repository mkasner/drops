package drops

import (
	"sync"
)

type Event uint

type event struct {
	Event Event
	Data  interface{}
	Err   error
}

type handler struct {
	h EventHandler
}

type EventEmitter struct {
	eventCh      chan *event
	handlerCh    chan EventHandler
	handlers     []handler
	wg           sync.WaitGroup // Emitter process wait group
	closeCh      chan struct{}
	recvCh       chan struct{}
	stopDispatch bool
}

func NewEventEmitter() *EventEmitter {
	ee := &EventEmitter{
		eventCh:   make(chan *event, 100),
		handlerCh: make(chan EventHandler),
		handlers:  make([]handler, 0),
		closeCh:   make(chan struct{}),
		recvCh:    make(chan struct{}),
	}
	go ee.register()
	go ee.handle()

	return ee
}

func (t *EventEmitter) Listen() {

}

func (t *EventEmitter) register() {
	for {
		select {
		case <-t.closeCh:
			return
		case h := <-t.handlerCh:
			t.handlers = append(t.handlers, handler{h: h})
		}
	}
}

func (t *EventEmitter) handle() {
	for {
		select {
		case <-t.closeCh:
			return
		case e := <-t.eventCh:
			t.wg.Add(len(t.handlers))
			for _, h := range t.handlers {
				go func(eh EventHandler) {
					eh.Handle(e.Event, e.Data, e.Err)
					t.wg.Done()
				}(h.h)
			}
			t.recvCh <- struct{}{}
		}
	}
}

func (t *EventEmitter) Dispatch(e Event, data interface{}, err error) {
	if t.stopDispatch {
		return
	}
	t.eventCh <- &event{
		Event: e,
		Data:  data,
		Err:   err,
	}
	<-t.recvCh
}

func (t *EventEmitter) Handler(handler EventHandler) {
	t.handlerCh <- handler
}

func (t *EventEmitter) Close() {
	t.stopDispatch = true
	t.wg.Wait() // waiting for handlers to finish
	t.closeCh <- struct{}{}
	t.closeCh <- struct{}{}
}

type EventHandler interface {
	Handle(Event, interface{}, error)
}

type EventHandlerFunc func(Event, interface{}, error)

func (t EventHandlerFunc) Handle(e Event, data interface{}, err error) {
	t(e, data, err)
}
