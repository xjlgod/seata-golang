package utils

import (
	"context"
	"errors"
	"github.com/seata/seata-go/pkg/saga/statemachine/constant"
	"github.com/seata/seata-go/pkg/saga/statemachine/engine"
	"github.com/seata/seata-go/pkg/saga/statemachine/engine/process_ctrl"
	"github.com/seata/seata-go/pkg/saga/statemachine/statelang"
	"github.com/seata/seata-go/pkg/saga/statemachine/statelang/state"
	"github.com/seata/seata-go/pkg/util/log"
	"golang.org/x/sync/semaphore"
	"reflect"
	"strings"
	"sync"
	"time"
)

func EndStateMachine(ctx context.Context, processContext process_ctrl.ProcessContext) error {
	if processContext.HasVariable(constant.VarNameIsLoopState) {
		if processContext.HasVariable(constant.LoopSemaphore) {
			weighted, ok := processContext.GetVariable(constant.LoopSemaphore).(semaphore.Weighted)
			if !ok {
				return errors.New("semaphore type is not weighted")
			}
			weighted.Release(1)
		}
	}

	stateMachineInstance, ok := processContext.GetVariable(constant.VarNameStateMachineInst).(statelang.StateMachineInstance)
	if !ok {
		return errors.New("state machine instance type is not statelang.StateMachineInstance")
	}

	stateMachineInstance.SetEndTime(time.Now())

	exp, ok := processContext.GetVariable(constant.VarNameCurrentException).(error)
	if !ok {
		return errors.New("exception type is not error")
	}

	if exp != nil {
		stateMachineInstance.SetException(exp)
		log.Debugf("Exception Occurred: %s", exp)
	}

	stateMachineConfig, ok := processContext.GetVariable(constant.VarNameStateMachineConfig).(engine.StateMachineConfig)

	if err := stateMachineConfig.StatusDecisionStrategy().DecideOnEndState(ctx, processContext, stateMachineInstance, exp); err != nil {
		return err
	}

	contextParams, ok := processContext.GetVariable(constant.VarNameStateMachineContext).(map[string]interface{})
	if !ok {
		return errors.New("state machine context type is not map[string]interface{}")
	}
	endParams := stateMachineInstance.EndParams()
	for k, v := range contextParams {
		endParams[k] = v
	}
	stateMachineInstance.SetEndParams(endParams)

	stateInstruction, ok := processContext.GetInstruction().(process_ctrl.StateInstruction)
	if !ok {
		return errors.New("state instruction type is not process_ctrl.StateInstruction")
	}
	stateInstruction.SetEnd(true)

	stateMachineInstance.SetRunning(false)
	stateMachineInstance.SetEndTime(time.Now())

	if stateMachineInstance.StateMachine().IsPersist() && stateMachineConfig.StateLangStore() != nil {
		err := stateMachineConfig.StateLogStore().RecordStateMachineFinished(ctx, stateMachineInstance, processContext)
		if err != nil {
			return err
		}
	}

	callBack, ok := processContext.GetVariable(constant.VarNameAsyncCallback).(engine.CallBack)
	if ok {
		if exp != nil {
			callBack.OnError(ctx, processContext, stateMachineInstance, exp)
		} else {
			callBack.OnFinished(ctx, processContext, stateMachineInstance)
		}
	}

	return nil
}

func HandleException(processContext process_ctrl.ProcessContext, abstractTaskState *state.AbstractTaskState, err error) {
	catches := abstractTaskState.Catches()
	if catches != nil && len(catches) != 0 {
		for _, exceptionMatch := range catches {
			exceptions := exceptionMatch.Exceptions()
			exceptionTypes := exceptionMatch.ExceptionTypes()
			if exceptions != nil && len(exceptions) != 0 {
				if exceptionTypes == nil {
					lock := processContext.GetVariable(constant.VarNameProcessContextMutexLock).(*sync.Mutex)
					lock.Lock()
					defer lock.Unlock()
					error := errors.New("")
					for i := 0; i < len(exceptions); i++ {
						exceptionTypes = append(exceptionTypes, reflect.TypeOf(error))
					}
				}

				exceptionMatch.SetExceptionTypes(exceptionTypes)
			}

			for i, _ := range exceptionTypes {
				if reflect.TypeOf(err) == exceptionTypes[i] {
					// HACK: we can not get error type in config file during runtime, so we use exception str
					if strings.Contains(err.Error(), exceptions[i]) {
						hierarchicalProcessContext := processContext.(process_ctrl.HierarchicalProcessContext)
						hierarchicalProcessContext.SetVariable(constant.VarNameCurrentExceptionRoute, exceptionMatch.Next())
						return
					}
				}
			}
		}
	}

	log.Error("Task execution failed and no catches configured")
	hierarchicalProcessContext := processContext.(process_ctrl.HierarchicalProcessContext)
	hierarchicalProcessContext.SetVariable(constant.VarNameIsExceptionNotCatch, true)
}

// GetOriginStateName get origin state name without suffix like fork
func GetOriginStateName(stateInstance statelang.StateInstance) string {
	stateName := stateInstance.Name()
	if stateName != "" {
		end := strings.LastIndex(stateName, constant.LoopStateNamePattern)
		if end > -1 {
			return stateName[:end+1]
		}
	}
	return stateName
}

// IsTimeout test if is timeout
func IsTimeout(gmtUpdated time.Time, timeoutMillis int) bool {
	if timeoutMillis < 0 {
		return false
	}
	return time.Now().Unix()-gmtUpdated.Unix() > int64(timeoutMillis)
}
