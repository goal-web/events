package events

import (
	"fmt"
	"github.com/goal-web/contracts"
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
	// 处理异常
	defer func() {
		if err := recover(); err != nil {
			go func() {
				dispatcher.exceptionHandler.Handle(EventException{
					error:  fmt.Errorf("%v", err),
					fields: nil,
					event:  event,
				})
			}()
		}
	}()

	if _, isSync := event.(contracts.SyncEvent); isSync {
		// 同步执行事件
		for _, listener := range dispatcher.getListeners(event.Event()) {
			listener.Handle(event)
		}
	} else {
		// 协程执行
		go func() {
			for _, listener := range dispatcher.getListeners(event.Event()) {
				listener.Handle(event)
			}
		}()
	}
}
