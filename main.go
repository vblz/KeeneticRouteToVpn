package main

import (
	"fmt"
	"github.com/go-pkgz/lgr"
	"github.com/vblz/KeeneticRouteToVpn/keeneticmanager"
	"github.com/vblz/KeeneticRouteToVpn/options"
	"github.com/vblz/KeeneticRouteToVpn/sshcommander"
	"os"
)

func main() {
	opts, err := options.Parse(os.Args)
	if err != nil {
		lgr.Printf("[ERROR] %v", err)
		os.Exit(1)
	}

	sshClient, err := sshcommander.NewCommander(opts.Host, opts.Port, opts.Username, opts.Password, opts.IgnoreHostChecking)
	if err != nil {
		lgr.Fatalf("%v", err)
		os.Exit(2)
	}

	err = sshClient.Dial()
	if err != nil {
		lgr.Fatalf("error while dialing: %v", err)
		os.Exit(2)
	}

	err = setupHostRules(sshClient, opts.HostsList, opts.InterfaceName)
	if err != nil {
		lgr.Fatalf("%v", err)
		os.Exit(3)
	}

	err = sshClient.Close()
	if err != nil {
		lgr.Printf("[ERROR] can't close ssh connection: %s", err)
	}

	lgr.Printf("[INFO] finished")
}

func setupHostRules(sshClient keeneticmanager.Executor, listHosts []string, interfaceName string) error {
	manager := keeneticmanager.NewKeeneticManager(sshClient)
	err := manager.CleanExistingGeneratedIPRules()

	if err != nil {
		return fmt.Errorf("can't clean up existing rules: %w", err)
	}

	for _, host := range listHosts {
		err = manager.AddVpnRoute(host, interfaceName)
		if err != nil {
			return fmt.Errorf("can't setup the rule for host '%s': %w", host, err)
		}
	}

	err = manager.SaveSettings()
	if err != nil {
		return fmt.Errorf("can't save settings: %w", err)
	}

	return nil
}
