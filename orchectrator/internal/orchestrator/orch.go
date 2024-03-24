package orchestrator

import (
	"encoding/json"
	"errors"
	"github.com/ilyazat/task_tracker/orchestrator/internal/model"
	amqp "github.com/rabbitmq/amqp091-go"
	"strings"
)

type Orchestrator struct {
}

// domain.service.handler
func NewOrchestrator() *Orchestrator {
	return &Orchestrator{}
}

func (o *Orchestrator) Dispatch(msg amqp.Delivery) error {
	var event model.Event

	if err := json.Unmarshal(msg.Body, &event); err != nil {
		return err
	}
	parts := strings.Split(event.Name, ".")

	if len(parts) != 3 {
		return errors.New("bad name of the event.Name")
	}

	switch event.Name {
	case "fuck":

	}
}
