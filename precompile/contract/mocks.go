// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/skychains/evm/precompile/contract (interfaces: BlockContext,AccessibleState)

// Package contract is a generated GoMock package.
package contract

import (
	big "math/big"
	reflect "reflect"

	snow "github.com/skychains/chain/snow"
	precompileconfig "github.com/skychains/evm/precompile/precompileconfig"
	common "github.com/ethereum/go-ethereum/common"
	gomock "go.uber.org/mock/gomock"
)

// MockBlockContext is a mock of BlockContext interface.
type MockBlockContext struct {
	ctrl     *gomock.Controller
	recorder *MockBlockContextMockRecorder
}

// MockBlockContextMockRecorder is the mock recorder for MockBlockContext.
type MockBlockContextMockRecorder struct {
	mock *MockBlockContext
}

// NewMockBlockContext creates a new mock instance.
func NewMockBlockContext(ctrl *gomock.Controller) *MockBlockContext {
	mock := &MockBlockContext{ctrl: ctrl}
	mock.recorder = &MockBlockContextMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBlockContext) EXPECT() *MockBlockContextMockRecorder {
	return m.recorder
}

// GetPredicateResults mocks base method.
func (m *MockBlockContext) GetPredicateResults(arg0 common.Hash, arg1 common.Address) []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPredicateResults", arg0, arg1)
	ret0, _ := ret[0].([]byte)
	return ret0
}

// GetPredicateResults indicates an expected call of GetPredicateResults.
func (mr *MockBlockContextMockRecorder) GetPredicateResults(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPredicateResults", reflect.TypeOf((*MockBlockContext)(nil).GetPredicateResults), arg0, arg1)
}

// Number mocks base method.
func (m *MockBlockContext) Number() *big.Int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Number")
	ret0, _ := ret[0].(*big.Int)
	return ret0
}

// Number indicates an expected call of Number.
func (mr *MockBlockContextMockRecorder) Number() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Number", reflect.TypeOf((*MockBlockContext)(nil).Number))
}

// Timestamp mocks base method.
func (m *MockBlockContext) Timestamp() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Timestamp")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// Timestamp indicates an expected call of Timestamp.
func (mr *MockBlockContextMockRecorder) Timestamp() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Timestamp", reflect.TypeOf((*MockBlockContext)(nil).Timestamp))
}

// MockAccessibleState is a mock of AccessibleState interface.
type MockAccessibleState struct {
	ctrl     *gomock.Controller
	recorder *MockAccessibleStateMockRecorder
}

// MockAccessibleStateMockRecorder is the mock recorder for MockAccessibleState.
type MockAccessibleStateMockRecorder struct {
	mock *MockAccessibleState
}

// NewMockAccessibleState creates a new mock instance.
func NewMockAccessibleState(ctrl *gomock.Controller) *MockAccessibleState {
	mock := &MockAccessibleState{ctrl: ctrl}
	mock.recorder = &MockAccessibleStateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccessibleState) EXPECT() *MockAccessibleStateMockRecorder {
	return m.recorder
}

// GetBlockContext mocks base method.
func (m *MockAccessibleState) GetBlockContext() BlockContext {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockContext")
	ret0, _ := ret[0].(BlockContext)
	return ret0
}

// GetBlockContext indicates an expected call of GetBlockContext.
func (mr *MockAccessibleStateMockRecorder) GetBlockContext() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockContext", reflect.TypeOf((*MockAccessibleState)(nil).GetBlockContext))
}

// GetChainConfig mocks base method.
func (m *MockAccessibleState) GetChainConfig() precompileconfig.ChainConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChainConfig")
	ret0, _ := ret[0].(precompileconfig.ChainConfig)
	return ret0
}

// GetChainConfig indicates an expected call of GetChainConfig.
func (mr *MockAccessibleStateMockRecorder) GetChainConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChainConfig", reflect.TypeOf((*MockAccessibleState)(nil).GetChainConfig))
}

// GetSnowContext mocks base method.
func (m *MockAccessibleState) GetSnowContext() *snow.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSnowContext")
	ret0, _ := ret[0].(*snow.Context)
	return ret0
}

// GetSnowContext indicates an expected call of GetSnowContext.
func (mr *MockAccessibleStateMockRecorder) GetSnowContext() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSnowContext", reflect.TypeOf((*MockAccessibleState)(nil).GetSnowContext))
}

// GetStateDB mocks base method.
func (m *MockAccessibleState) GetStateDB() StateDB {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStateDB")
	ret0, _ := ret[0].(StateDB)
	return ret0
}

// GetStateDB indicates an expected call of GetStateDB.
func (mr *MockAccessibleStateMockRecorder) GetStateDB() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStateDB", reflect.TypeOf((*MockAccessibleState)(nil).GetStateDB))
}
