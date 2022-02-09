// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package keeneticmanager

import (
	"sync"
)

// Ensure, that ExecutorMock does implement Executor.
// If this is not the case, regenerate this file with moq.
var _ Executor = &ExecutorMock{}

// ExecutorMock is a mock implementation of Executor.
//
// 	func TestSomethingThatUsesExecutor(t *testing.T) {
//
// 		// make and configure a mocked Executor
// 		mockedExecutor := &ExecutorMock{
// 			RunCommandFunc: func(command string) (string, error) {
// 				panic("mock out the RunCommand method")
// 			},
// 		}
//
// 		// use mockedExecutor in code that requires Executor
// 		// and then make assertions.
//
// 	}
type ExecutorMock struct {
	// RunCommandFunc mocks the RunCommand method.
	RunCommandFunc func(command string) (string, error)

	// calls tracks calls to the methods.
	calls struct {
		// RunCommand holds details about calls to the RunCommand method.
		RunCommand []struct {
			// Command is the command argument value.
			Command string
		}
	}
	lockRunCommand sync.RWMutex
}

// RunCommand calls RunCommandFunc.
func (mock *ExecutorMock) RunCommand(command string) (string, error) {
	if mock.RunCommandFunc == nil {
		panic("ExecutorMock.RunCommandFunc: method is nil but Executor.RunCommand was just called")
	}
	callInfo := struct {
		Command string
	}{
		Command: command,
	}
	mock.lockRunCommand.Lock()
	mock.calls.RunCommand = append(mock.calls.RunCommand, callInfo)
	mock.lockRunCommand.Unlock()
	return mock.RunCommandFunc(command)
}

// RunCommandCalls gets all the calls that were made to RunCommand.
// Check the length with:
//     len(mockedExecutor.RunCommandCalls())
func (mock *ExecutorMock) RunCommandCalls() []struct {
	Command string
} {
	var calls []struct {
		Command string
	}
	mock.lockRunCommand.RLock()
	calls = mock.calls.RunCommand
	mock.lockRunCommand.RUnlock()
	return calls
}
