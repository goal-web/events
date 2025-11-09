package events

import (
	"github.com/goal-web/contracts"
)


type ServiceProvider struct {
}

func NewService() contracts.ServiceProvider {
	return &ServiceProvider{}
}

func (provider ServiceProvider) Stop() {

}

func (provider ServiceProvider) Start() error {
	return nil
}

func (provider ServiceProvider) Register(container contracts.Application) {
	container.Singleton("events", func(handler contracts.ExceptionHandler) contracts.EventDispatcher {
		return NewDispatcher(handler)
	})
}
