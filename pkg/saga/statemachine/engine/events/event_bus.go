package events

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/seata/seata-go/pkg/saga/statemachine/constant"
	"github.com/seata/seata-go/pkg/saga/statemachine/engine/process_ctrl"
	"github.com/seata/seata-go/pkg/util/collection"
	"github.com/seata/seata-go/pkg/util/log"
)

type EventBus interface {
	Offer(ctx context.Context, event Event) (bool, error)

	EventConsumerList(event Event) []EventConsumer

	RegisterEventConsumer(consumer EventConsumer)
}

type BaseEventBus struct {
	eventConsumerList []EventConsumer
}

func (b *BaseEventBus) RegisterEventConsumer(consumer EventConsumer) {
	if b.eventConsumerList == nil {
		b.eventConsumerList = make([]EventConsumer, 0)
	}
	b.eventConsumerList = append(b.eventConsumerList, consumer)
}

func (b *BaseEventBus) EventConsumerList(event Event) []EventConsumer {
	var acceptedConsumerList = make([]EventConsumer, 0)
	for i := range b.eventConsumerList {
		eventConsumer := b.eventConsumerList[i]
		if eventConsumer.Accept(event) {
			acceptedConsumerList = append(acceptedConsumerList, eventConsumer)
		}
	}
	return acceptedConsumerList
}

type DirectEventBus struct {
	BaseEventBus
}

func (d DirectEventBus) Offer(ctx context.Context, event Event) (bool, error) {
	eventConsumerList := d.EventConsumerList(event)
	if len(eventConsumerList) == 0 {
		log.Debugf("cannot find event handler by type: %T", event)
		return false, nil
	}

	isFirstEvent := true
	processContext, ok := event.(process_ctrl.ProcessContext)
	if !ok {
		log.Errorf("event %T is illegal, required process_ctrl.ProcessContext", event)
		return false, nil
	}

	stack := processContext.GetVariable(constant.VarNameSyncExeStack).(*collection.Stack)
	if stack == nil {
		stack = collection.NewStack()
		processContext.SetVariable(constant.VarNameSyncExeStack, stack)
		isFirstEvent = true
	}

	stack.Push(processContext)
	if isFirstEvent {
		for stack.Len() > 0 {
			currentContext := stack.Pop().(process_ctrl.ProcessContext)
			for _, eventConsumer := range eventConsumerList {
				err := eventConsumer.Process(ctx, currentContext)
				if err != nil {
					log.Errorf("process event %T error: %s", event, err.Error())
					return false, err
				}
			}
		}
	}

	return true, nil
}

type AsyncEventBus struct {
	BaseEventBus
}

func (a AsyncEventBus) Offer(ctx context.Context, event Event) (bool, error) {
	eventConsumerList := a.EventConsumerList(event)
	if len(eventConsumerList) == 0 {
		errStr := fmt.Sprintf("cannot find event handler by type: %T", event)
		log.Errorf(errStr)
		return false, errors.New(errStr)
	}

	processContext, ok := event.(process_ctrl.ProcessContext)
	if !ok {
		errStr := fmt.Sprintf("event %T is illegal, required process_ctrl.ProcessContext", event)
		log.Errorf(errStr)
		return false, errors.New(errStr)
	}

	for _, eventConsumer := range eventConsumerList {
		go func() {
			err := eventConsumer.Process(ctx, processContext)
			if err != nil {
				log.Errorf("process event %T error: %s", event, err.Error())
			}
		}()
	}

	return true, nil
}
