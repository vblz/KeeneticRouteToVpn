package keeneticmanager

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

var runningConfGeneratedRules = []interface{} {
	"ip route 10.1.1.1 Wireguard0 auto !KeeneticManager ip Rule",
	"ip route 4.4.4.4 Wireguard0 auto !KeeneticManager ip Rule",
	"ip route 4.4.0.0/16 AnotherInterface auto !KeeneticManager ip Rule",
}

var runningConf = fmt.Sprintf(`test
!@#!@#
!@#@!#!@#
%s
  %s
  ip route 4.4.4.4 Wireguard0 auto !another commend
ip route 5.5.5.5 Wireguard0 auto
	%s
asdasd
qweqwe
asdasd
`, runningConfGeneratedRules...)

func Test_manager_CleanExistingGeneratedIpRules(t *testing.T) {
	executorMock := &ExecutorMock{
		RunCommandFunc: func(command string) (string, error) {
			if command == "more running-config" {
				return runningConf, nil
			}
			return "OK", nil
		},
	}

	manager := NewKeeneticManager(executorMock)
	err := manager.CleanExistingGeneratedIPRules()
	assert.NoError(t, err)
	assert.Equal(t, len(runningConfGeneratedRules) + 1, len(executorMock.RunCommandCalls()))
	assert.Equal(t, "more running-config", executorMock.RunCommandCalls()[0].Command)
	for i, rule := range []string {
		"no ip route 10.1.1.1 Wireguard0",
		"no ip route 4.4.4.4 Wireguard0",
		"no ip route 4.4.0.0/16 AnotherInterface",
	} {
		assert.Equal(t, rule, executorMock.RunCommandCalls()[i+1].Command)
	}
}

func TestManager_GetRunningConfiguration(t *testing.T) {
	testData := []struct{
		stdOut string
		err error
	} {
		{ "test", nil },
		{ "", nil},
		{ "", errors.New("test")},
		{ "test", errors.New("test")},
	}

	for i, tt := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			executorMock := &ExecutorMock{
				RunCommandFunc: func(command string) (string, error) {
					require.Equal(t, "more running-config", command)
					return tt.stdOut, tt.err
				},
			}

			manager := NewKeeneticManager(executorMock)
			res, err := manager.GetRunningConfiguration()
			assert.Equal(t, tt.stdOut, res)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestManager_SaveSettings(t *testing.T) {
	testData := []struct{
		stdOut string
		err error
	} {
		{ "test", nil },
		{ "", nil},
		{ "", errors.New("test")},
		{ "test", errors.New("test")},
	}

	for i, tt := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			executorMock := &ExecutorMock{
				RunCommandFunc: func(command string) (string, error) {
					require.Equal(t, "copy running-config startup-config", command)
					return tt.stdOut, tt.err
				},
			}

			manager := NewKeeneticManager(executorMock)
			err := manager.SaveSettings()
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestManager_AddVpnRoute(t *testing.T) {
	const host = "10.0.0.0/8"
	const interfaceName = "TEST_INTERFACE_NAME5"
	expectedCommand := "ip route 10.0.0.0/8 TEST_INTERFACE_NAME5 auto !KeeneticManager ip Rule"

	testData := []struct{
		stdOut string
		err error
	} {
		{ "test", nil },
		{ "", nil},
		{ "", errors.New("test")},
		{ "test", errors.New("test")},
	}

	for i, tt := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			executorMock := &ExecutorMock{
				RunCommandFunc: func(command string) (string, error) {
					require.Equal(t, expectedCommand, command)
					return tt.stdOut, tt.err
				},
			}

			manager := NewKeeneticManager(executorMock)
			err := manager.AddVpnRoute(host, interfaceName)
			assert.Equal(t, tt.err, err)
		})
	}
}
