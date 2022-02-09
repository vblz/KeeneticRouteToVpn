package keeneticmanager

//go:generate moq -out keeneticManager_mock_test.go . Executor

import (
	"fmt"
	"strings"
)

// Executor executes commands at Keenetic
type Executor interface {
	RunCommand(command string) (string, error)
}

type manager struct {
	exc Executor
}

// NewKeeneticManager creates and returns a new instance of keenetic manager
func NewKeeneticManager(executor Executor) *manager {
	return &manager{
		exc: executor,
	}
}

func (m *manager) CleanExistingGeneratedIPRules() error {
	settings, err := m.GetRunningConfiguration()
	if err != nil {
		return fmt.Errorf("can't get config: %w", err)
	}

	generatedRules := make([]string, 0)

	for _, line := range strings.Split(settings, "\n") {
		trimmed := strings.TrimSpace(line)
		if isGeneratedIPRule(trimmed) {
			generatedRules = append(generatedRules, trimmed)
		}
	}

	for _, oldRule := range generatedRules {
		cleaned := strings.TrimSuffix(oldRule, ipRouteSuffix)
		_, err := m.exc.RunCommand("no " + cleaned)
		if err != nil {
			return fmt.Errorf("can't remove generated route: %w", err)
		}
	}

	return nil
}

func (m *manager) AddVpnRoute(host, interfaceName string) error {
	command := makeIPRouteCommand(host, interfaceName)
	_, err := m.exc.RunCommand(command)
	return err
}

func (m *manager) SaveSettings() error {
	_, err := m.exc.RunCommand("copy running-config startup-config")
	return err
}

func (m *manager) GetRunningConfiguration() (string, error) {
	return m.exc.RunCommand("more running-config")
}

func isGeneratedIPRule(trimmedConfigLine string) bool {
	return strings.HasPrefix(trimmedConfigLine, ipRouteCommand) && strings.HasSuffix(trimmedConfigLine, ipRulesComment)
}
