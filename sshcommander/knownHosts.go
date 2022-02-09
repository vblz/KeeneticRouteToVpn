package sshcommander

import (
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"os"
	"path"
)

func getHostCallbackFunc(ignoreKnownHosts bool) (ssh.HostKeyCallback, error) {
	/* #nosec */
	if ignoreKnownHosts {
		return ssh.InsecureIgnoreHostKey(), nil
	}

	knownHostsFiles := getKnownHostLocations()

	return knownhosts.New(knownHostsFiles...)
}

func getKnownHostLocations() []string {
	knownHostsFiles := make([]string, 0, 2)
	if checkFile("/etc/ssh/ssh_known_hosts") {
		knownHostsFiles = append(knownHostsFiles, "/etc/ssh/ssh_known_hosts")
	}

	homeDir, err := os.UserHomeDir()
	if err == nil {
		homeKnownHostPath := path.Join(homeDir, ".ssh", "known_hosts")
		if checkFile(homeKnownHostPath) {
			knownHostsFiles = append(knownHostsFiles, homeKnownHostPath)
		}
	}

	return knownHostsFiles
}

func checkFile(filePath string) bool {
	s, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	return !s.IsDir()
}
