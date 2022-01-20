package events

import (
	"github.com/goal-web/contracts"
)

type EventException struct {
	error
	fields contracts.Fields
	event  contracts.Event
}

func (e EventException) Error() string {
	return e.error.Error()
}

func (e EventException) Fields() contracts.Fields {
	return e.fields
}
