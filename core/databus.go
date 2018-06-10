package core

import (
	"github.com/asaskevich/EventBus"
)

type DataBus struct {
	EventBus.Bus
}

func NewDataBus() *DataBus {
	return &DataBus{
		Bus: EventBus.New(),
	}
}
