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
	h          EventHandler
	stop, done chan struct{}
}

type EventEmitter struct {
	eventCh   chan *event
	handlerCh chan EventHandler
	handlers  []handler
	wg        sync.WaitGroup
	closeCh   chan struct{}
}

func NewEventEmitter() *EventEmitter {
	ee := &EventEmitter{
		eventCh:   make(chan *event, 100),
		handlerCh: make(chan EventHandler, 1),
		handlers:  make([]handler, 0),
		closeCh:   make(chan struct{}, 1),
	}
	ee.wg.Add(2)
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
			t.closeCh <- struct{}{}
			t.wg.Done()
			return
		case h := <-t.handlerCh:
			t.handlers = append(t.handlers, handler{h: h, stop: make(chan struct{}), done: make(chan struct{})})
		}
	}
}

func (t *EventEmitter) handle() {
	for {
		select {
		case <-t.closeCh:
			t.closeCh <- struct{}{}
			t.wg.Done()
			return
		case e := <-t.eventCh:
			for _, h := range t.handlers {
				go func(h EventHandler) {
					h.Handle(e.Event, e.Data, e.Err)
				}(h.h)
			}
		}
	}
}

func (t *EventEmitter) Dispatch(e Event, data interface{}, err error) {
	go func() {
		t.eventCh <- &event{
			Event: e,
			Data:  data,
			Err:   err,
		}
	}()
}

func (t *EventEmitter) Handler(handler EventHandler) {
	t.handlerCh <- handler
}

func (t *EventEmitter) Close() {
	var wg sync.WaitGroup
	for _, h := range t.handlers {
		wg.Add(1)
		go func() {
			go h.h.Stop(h.stop, h.done)
			h.stop <- struct{}{}
			<-h.done
			wg.Done()
		}()
		wg.Wait()
	}
	// closing event emmiter processes for register and dispatch
	t.closeCh <- struct{}{}
	t.wg.Wait()
	close(t.closeCh)
}

type EventHandler interface {
	Handle(Event, interface{}, error)
	Stop(chan struct{}, chan struct{})
}

type EventHandlerFunc func(Event, interface{}, error)

func (t EventHandlerFunc) Handle(e Event, data interface{}, err error) {
	t(e, data, err)
}

func (t EventHandlerFunc) Stop(stop, done chan struct{}) {
	<-stop
	done <- struct{}{}
}
