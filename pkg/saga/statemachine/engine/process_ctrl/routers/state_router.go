package routers

import (
	"context"
	"github.com/seata/seata-go/pkg/saga/statemachine/constant"
	"github.com/seata/seata-go/pkg/saga/statemachine/engine/exception"
	"github.com/seata/seata-go/pkg/saga/statemachine/engine/process_ctrl"
	"github.com/seata/seata-go/pkg/saga/statemachine/engine/process_ctrl/utils"
	"github.com/seata/seata-go/pkg/saga/statemachine/statelang"
	sagaState "github.com/seata/seata-go/pkg/saga/statemachine/statelang/state"
	seataErrors "github.com/seata/seata-go/pkg/util/errors"
	"github.com/seata/seata-go/pkg/util/log"
)

type EndStateRouter struct {
}

func (e EndStateRouter) Route(ctx context.Context, processContext process_ctrl.ProcessContext, state statelang.State) (process_ctrl.Instruction, error) {
	return nil, nil
}

type TaskStateRouter struct {
}

func (t TaskStateRouter) Route(ctx context.Context, processContext process_ctrl.ProcessContext, state statelang.State) (process_ctrl.Instruction, error) {
	stateInstruction, _ := processContext.GetInstruction().(process_ctrl.StateInstruction)
	if stateInstruction.End() {
		log.Infof("StateInstruction is ended, Stop the StateMachine executing. StateMachine[%s] Current State[%s]",
			stateInstruction.StateMachineName(), stateInstruction.StateName())
	}

	// check if in loop async condition
	isLoop, ok := processContext.GetVariable(constant.VarNameIsLoopState).(bool)
	if ok && isLoop {
		log.Infof("StateMachine[%s] Current State[%s] is in loop async condition, skip route processing.",
			stateInstruction.StateMachineName(), stateInstruction.StateName())
		return nil, nil
	}

	// The current CompensationTriggerState can mark the compensation process is started and perform compensation
	// route processing.
	compensationTriggerState, ok := processContext.GetVariable(constant.VarNameCurrentCompensateTriggerState).(statelang.State)
	if ok {
		return t.compensateRoute(ctx, processContext, compensationTriggerState)
	}

	// There is an exception route, indicating that an exception is thrown, and the exception route is prioritized.
	next := processContext.GetVariable(constant.VarNameCurrentExceptionRoute).(string)

	if next != "" {
		processContext.RemoveVariable(constant.VarNameCurrentExceptionRoute)
	} else {
		next = state.Next()
	}

	// If next is empty, the state selected by the Choice state was taken.
	if next == "" && processContext.HasVariable(constant.VarNameCurrentChoice) {
		next = processContext.GetVariable(constant.VarNameCurrentChoice).(string)
		processContext.RemoveVariable(constant.VarNameCurrentChoice)
	}

	if next == "" {
		return nil, nil
	}

	stateMachine := state.StateMachine()
	nextState := stateMachine.State(next)
	if nextState == nil {
		return nil, exception.NewEngineExecutionException(seataErrors.ObjectNotExists,
			"Next state["+next+"] is not exits", nil)
	}

	stateInstruction.SetStateName(next)

	if nil != utils.GetLoopConfig(ctx, processContext, nextState) {
		stateInstruction.SetTemporaryState(sagaState.NewLoopStartStateImpl())
	}

	return stateInstruction, nil
}

func (t *TaskStateRouter) compensateRoute(ctx context.Context, processContext process_ctrl.ProcessContext, state statelang.State) (process_ctrl.Instruction, error) {
	//If there is already a compensation state that has been executed,
	// it is judged whether it is wrong or unsuccessful,
	// and the compensation process is interrupted.
	// TODO implement me
	return nil, nil
}
