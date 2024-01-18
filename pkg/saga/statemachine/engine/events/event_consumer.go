package events

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/seata/seata-go/pkg/saga/statemachine/engine/process_ctrl"
)

type EventConsumer interface {
	Accept(event Event) bool

	Process(ctx context.Context, event Event) error
}

type ProcessCtrlEventConsumer struct {
	processController process_ctrl.ProcessController
}

func (p ProcessCtrlEventConsumer) Accept(event Event) bool {
	if event == nil {
		return false
	}

	_, ok := event.(process_ctrl.ProcessContext)
	return ok
}

func (p ProcessCtrlEventConsumer) Process(ctx context.Context, event Event) error {
	processContext, ok := event.(process_ctrl.ProcessContext)
	if !ok {
		return errors.New(fmt.Sprint("event %T is illegal, required process_ctrl.ProcessContext", event))
	}
	return p.processController.Process(ctx, processContext)
}
