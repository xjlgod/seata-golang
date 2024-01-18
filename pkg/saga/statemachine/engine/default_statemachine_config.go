package engine

import (
	"github.com/seata/seata-go/pkg/saga/statemachine/engine/events"
	"github.com/seata/seata-go/pkg/saga/statemachine/engine/expr"
	"github.com/seata/seata-go/pkg/saga/statemachine/engine/invoker"
	"github.com/seata/seata-go/pkg/saga/statemachine/engine/sequence"
	"github.com/seata/seata-go/pkg/saga/statemachine/engine/status_decision"
	"github.com/seata/seata-go/pkg/saga/statemachine/engine/store"
	"sync"
)

const (
	DefaultTransOperTimeout     = 60000 * 30
	DefaultServiceInvokeTimeout = 60000 * 5
)

type DefaultStateMachineConfig struct {
	// Configuration
	transOperationTimeout int
	serviceInvokeTimeout  int
	charset               string
	defaultTenantId       string

	// Components

	// Event publisher
	syncProcessCtrlEventPublisher  events.EventPublisher
	asyncProcessCtrlEventPublisher events.EventPublisher

	// Store related components
	stateLogRepository     store.StateLogRepository
	stateLogStore          store.StateLogStore
	stateLangStore         store.StateLangStore
	stateMachineRepository store.StateMachineRepository

	// Expression related components
	expressionFactoryManager expr.ExpressionFactoryManager
	expressionResolver       expr.ExpressionResolver

	// Invoker related components
	serviceInvokerManager invoker.ServiceInvokerManager
	scriptInvokerManager  invoker.ScriptInvokerManager

	// Other components
	statusDecisionStrategy status_decision.StatusDecisionStrategy
	seqGenerator           sequence.SeqGenerator
	componentLock          *sync.Mutex
}

func (c *DefaultStateMachineConfig) ComponentLock() *sync.Mutex {
	return c.componentLock
}

func (c *DefaultStateMachineConfig) SetComponentLock(componentLock *sync.Mutex) {
	c.componentLock = componentLock
}

func (c *DefaultStateMachineConfig) SetTransOperationTimeout(transOperationTimeout int) {
	c.transOperationTimeout = transOperationTimeout
}

func (c *DefaultStateMachineConfig) SetServiceInvokeTimeout(serviceInvokeTimeout int) {
	c.serviceInvokeTimeout = serviceInvokeTimeout
}

func (c *DefaultStateMachineConfig) SetCharset(charset string) {
	c.charset = charset
}

func (c *DefaultStateMachineConfig) SetDefaultTenantId(defaultTenantId string) {
	c.defaultTenantId = defaultTenantId
}

func (c *DefaultStateMachineConfig) SetSyncProcessCtrlEventPublisher(syncProcessCtrlEventPublisher events.EventPublisher) {
	c.syncProcessCtrlEventPublisher = syncProcessCtrlEventPublisher
}

func (c *DefaultStateMachineConfig) SetAsyncProcessCtrlEventPublisher(asyncProcessCtrlEventPublisher events.EventPublisher) {
	c.asyncProcessCtrlEventPublisher = asyncProcessCtrlEventPublisher
}

func (c *DefaultStateMachineConfig) SetStateLogRepository(stateLogRepository store.StateLogRepository) {
	c.stateLogRepository = stateLogRepository
}

func (c *DefaultStateMachineConfig) SetStateLogStore(stateLogStore store.StateLogStore) {
	c.stateLogStore = stateLogStore
}

func (c *DefaultStateMachineConfig) SetStateLangStore(stateLangStore store.StateLangStore) {
	c.stateLangStore = stateLangStore
}

func (c *DefaultStateMachineConfig) SetStateMachineRepository(stateMachineRepository store.StateMachineRepository) {
	c.stateMachineRepository = stateMachineRepository
}

func (c *DefaultStateMachineConfig) SetExpressionFactoryManager(expressionFactoryManager expr.ExpressionFactoryManager) {
	c.expressionFactoryManager = expressionFactoryManager
}

func (c *DefaultStateMachineConfig) SetExpressionResolver(expressionResolver expr.ExpressionResolver) {
	c.expressionResolver = expressionResolver
}

func (c *DefaultStateMachineConfig) SetServiceInvokerManager(serviceInvokerManager invoker.ServiceInvokerManager) {
	c.serviceInvokerManager = serviceInvokerManager
}

func (c *DefaultStateMachineConfig) SetScriptInvokerManager(scriptInvokerManager invoker.ScriptInvokerManager) {
	c.scriptInvokerManager = scriptInvokerManager
}

func (c *DefaultStateMachineConfig) SetStatusDecisionStrategy(statusDecisionStrategy status_decision.StatusDecisionStrategy) {
	c.statusDecisionStrategy = statusDecisionStrategy
}

func (c *DefaultStateMachineConfig) SetSeqGenerator(seqGenerator sequence.SeqGenerator) {
	c.seqGenerator = seqGenerator
}

func (c *DefaultStateMachineConfig) StateLogRepository() store.StateLogRepository {
	return c.stateLogRepository
}

func (c *DefaultStateMachineConfig) StateMachineRepository() store.StateMachineRepository {
	return c.stateMachineRepository
}

func (c *DefaultStateMachineConfig) StateLogStore() store.StateLogStore {
	return c.stateLogStore
}

func (c *DefaultStateMachineConfig) StateLangStore() store.StateLangStore {
	return c.stateLangStore
}

func (c *DefaultStateMachineConfig) ExpressionFactoryManager() expr.ExpressionFactoryManager {
	return c.expressionFactoryManager
}

func (c *DefaultStateMachineConfig) ExpressionResolver() expr.ExpressionResolver {
	return c.expressionResolver
}

func (c *DefaultStateMachineConfig) SeqGenerator() sequence.SeqGenerator {
	return c.seqGenerator
}

func (c *DefaultStateMachineConfig) StatusDecisionStrategy() status_decision.StatusDecisionStrategy {
	return c.statusDecisionStrategy
}

func (c *DefaultStateMachineConfig) EventPublisher() events.EventPublisher {
	return c.syncProcessCtrlEventPublisher
}

func (c *DefaultStateMachineConfig) AsyncEventPublisher() events.EventPublisher {
	return c.asyncProcessCtrlEventPublisher
}

func (c *DefaultStateMachineConfig) ServiceInvokerManager() invoker.ServiceInvokerManager {
	return c.serviceInvokerManager
}

func (c *DefaultStateMachineConfig) ScriptInvokerManager() invoker.ScriptInvokerManager {
	return c.scriptInvokerManager
}

func (c *DefaultStateMachineConfig) CharSet() string {
	return c.charset
}

func (c *DefaultStateMachineConfig) SetCharSet(charset string) {
	c.charset = charset
}

func (c *DefaultStateMachineConfig) DefaultTenantId() string {
	return c.defaultTenantId
}

func (c *DefaultStateMachineConfig) TransOperationTimeout() int {
	return c.transOperationTimeout
}

func (c *DefaultStateMachineConfig) ServiceInvokeTimeout() int {
	return c.serviceInvokeTimeout
}

func NewDefaultStateMachineConfig() *DefaultStateMachineConfig {
	c := &DefaultStateMachineConfig{
		transOperationTimeout: DefaultTransOperTimeout,
		serviceInvokeTimeout:  DefaultServiceInvokeTimeout,
		charset:               "UTF-8",
		defaultTenantId:       "000001",
		componentLock:         &sync.Mutex{},
	}
	return c
}
