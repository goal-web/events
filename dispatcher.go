package events

import (
	"fmt"
	"github.com/goal-web/contracts"
	"github.com/qbhy/parallel"
	"sync"
)

func NewDispatcher(handler contracts.ExceptionHandler) contracts.EventDispatcher {
	return &EventDispatcher{
		eventListenersMap: sync.Map{},
		exceptionHandler:  handler,
	}
}

type EventDispatcher struct {
	eventListenersMap sync.Map

	// 依赖异常处理器
	exceptionHandler contracts.ExceptionHandler
}

func (dispatcher *EventDispatcher) Register(name string, listener contracts.EventListener) {
	dispatcher.eventListenersMap.Store(name, append(dispatcher.getListeners(name), listener))
}
func (dispatcher *EventDispatcher) getListeners(name string) []contracts.EventListener {
	if value, exists := dispatcher.eventListenersMap.Load(name); exists {
		return value.([]contracts.EventListener)
	}
	return nil
}

func (dispatcher *EventDispatcher) Dispatch(event contracts.Event) {
	if e, isSync := event.(contracts.SyncEvent); isSync && e.Sync() {
		// 同步执行事件
		dispatcher.handleEvent(event)
	} else {
		// 协程执行
		go func() {
			dispatcher.handleEvent(event)
		}()
	}
}

func (dispatcher *EventDispatcher) exceptionHandle(err interface{}, event contracts.Event) {
	if err != nil {
		dispatcher.exceptionHandler.Handle(EventException{
			error:  fmt.Errorf("%v", err),
			fields: nil,
			event:  event,
		})
	}
}

func (dispatcher *EventDispatcher) handleEvent(event contracts.Event) {
	defer func() {
		dispatcher.exceptionHandle(recover(), event)
	}()

	listeners := dispatcher.getListeners(event.Event())
	parallelInstance := parallel.NewParallel(len(listeners))

	for _, listener := range listeners {
		_ = parallelInstance.Add(func() interface{} {
			listener.Handle(event)
			return nil
		})
	}

	for _, result := range parallelInstance.Wait() {
		dispatcher.exceptionHandle(result, event)
	}
}
