package sshcommander

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
)

// do not confuse with ssh.session
type sshCommander struct {
	host string
	port uint16

	client *ssh.Client
	config *ssh.ClientConfig
}

// NewCommander creates and return new instance of ssh commander
func NewCommander(host string, port uint16, username, password string, ignoreKnownHosts bool) (*sshCommander, error) {
	hostCallbackFunc, err := getHostCallbackFunc(ignoreKnownHosts)
	if err != nil {
		return nil, fmt.Errorf("can't init host fuction: %w", err)
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: hostCallbackFunc,
	}

	return &sshCommander{
		host:   host,
		port:   port,
		config: config,
	}, nil
}

func (s *sshCommander) Close() error {
	if s.client == nil {
		return fmt.Errorf("client wasn't opened")
	}
	defer func() {
		s.client = nil
	}()

	return s.client.Close()
}

func (s *sshCommander) Dial() error {
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", s.host, s.port), s.config)

	if err == nil {
		s.client = client
	}

	return err
}

func (s *sshCommander) RunCommand(command string) (string, error) {
	stdOut, stdErr, err := s.RunCommandWithStdErr(command)
	if err != nil {
		return stdOut, err
	}

	if stdErr != "" {
		return stdOut, StdErrError{StdErrorText: stdErr}
	}

	return stdOut, nil
}

func (s *sshCommander) RunCommandWithStdErr(command string) (stdOut, stdErr string, err error) {
	if s.client == nil {
		return "", "", fmt.Errorf("you should dial first")
	}
	session, err := s.client.NewSession()
	if err != nil {
		return "", "", err
	}
	defer session.Close()

	var stdOutBuf, stdErrBuf bytes.Buffer
	session.Stdout = &stdOutBuf
	session.Stderr = &stdErrBuf

	err = session.Run(command)

	return stdOutBuf.String(), stdErrBuf.String(), err
}
