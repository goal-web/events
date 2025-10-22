package events

import (
	"sync"

	"github.com/goal-web/application"
	"github.com/goal-web/contracts"
)

var dispatcher contracts.EventDispatcher

var once sync.Once

type ServiceProvider struct {
}

func NewService() contracts.ServiceProvider {
	return &ServiceProvider{}
}

func Dispatch(event contracts.Event) bool {
	once.Do(func() {
		if dispatcher == nil {
			dispatcher = application.Get("evens").(contracts.EventDispatcher)
		}
	})
	dispatcher.Dispatch(event)
	return true
}

func (provider ServiceProvider) Stop() {

}

func (provider ServiceProvider) Start() error {
	return nil
}

func (provider ServiceProvider) Register(container contracts.Application) {
	container.Singleton("events", func(handler contracts.ExceptionHandler) contracts.EventDispatcher {
		dispatcher = NewDispatcher(handler)
		return dispatcher
	})
}
