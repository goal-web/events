package events

import "github.com/goal-web/contracts"

var dispatcher contracts.EventDispatcher

type ServiceProvider struct {
}

func Dispatch(event contracts.Event) bool {
	if dispatcher != nil {
		dispatcher.Dispatch(event)
		return true
	}
	return false
}

func (this ServiceProvider) Stop() {

}

func (this ServiceProvider) Start() error {
	return nil
}

func (provider ServiceProvider) Register(container contracts.Application) {
	container.Singleton("events", func(handler contracts.ExceptionHandler) contracts.EventDispatcher {
		dispatcher = NewDispatcher(handler)
		return dispatcher
	})
}
