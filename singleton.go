package events

import (
	"sync"

	"github.com/goal-web/application"
	"github.com/goal-web/contracts"
)

var singleton contracts.EventDispatcher
var once sync.Once

func Default() contracts.EventDispatcher {
	once.Do(func() {
		singleton = application.Get("events").(contracts.EventDispatcher)
	})

	return singleton
}

func Dispatch(event contracts.Event) bool {
	Default().Dispatch(event)
	return true
}
