package events

import (
	"github.com/goal-web/contracts"
)

type EventException struct {
	error
	Event    contracts.Event
	Previous contracts.Exception
}

func (e EventException) GetPrevious() contracts.Exception {
	return nil
}
