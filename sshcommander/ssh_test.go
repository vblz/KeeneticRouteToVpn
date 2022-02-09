package sshcommander

import (
	sshServer "github.com/gliderlabs/ssh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
	"net"
	"testing"
	"time"
)

const testUsername = "testUserName"
const testPassword = "testPasswrod"

func TestNewCommander(t *testing.T) {
	const host = "192.168.1.100"
	const port = uint16(444)
	const ignoreKnownHosts = false

	c, err := NewCommander(host, port, testUsername, testPassword, ignoreKnownHosts)
	require.NoError(t, err)


	assert.Equal(t, host, c.host)
	assert.Equal(t, port, c.port)
	assert.Equal(t, testUsername, c.config.User)
	assert.Equal(t, 1, len(c.config.Auth))

	addr, _ := net.ResolveIPAddr("", "10.1.1.1")

	assert.Error(t, c.config.HostKeyCallback("test", addr, publickKeyMock{}))
}

func TestNewCommander_IgnoreKnownHosts(t *testing.T) {
	const ignoreKnownHosts = true

	c, err := NewCommander("", 0, "", "", ignoreKnownHosts)
	require.NoError(t, err)

	assert.NoError(t, c.config.HostKeyCallback("test", nil, nil))
}

func TestSshCommander_RunCommand(t *testing.T) {
	c, err := NewCommander("127.0.0.1", 4222, testUsername, testPassword, true)
	require.NoError(t, err)

	tearDown := makeSSHServer()
	defer tearDown()

	err = c.Dial()
	require.NoError(t, err)

	testData := []struct{
		name string
		command string
		expectedStdOut string
		expectError string
	} {
		{"no error, stdout", "out", "test out", ""},
		{"no error, no stdout", "", "", ""},
		{"error from exitCode", "out code", "test out", "Process exited with status 16"},
		{"error from exitCode with stdErr", "out code err", "test out", "Process exited with status 16"},
		{"error from with stdErr", "out err", "test out", "test error"},
		{"error from with stdErr, no stdout", "err", "", "test error"},
	}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			var stdOut string
			stdOut, err = c.RunCommand(tt.command)
			assert.Equal(t, tt.expectedStdOut, stdOut)
			if tt.expectError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectError)
			}
		})
	}

	err = c.Close()
	assert.NoError(t, err)
}

func TestSshCommander_RunCommandWithStdErr(t *testing.T) {
	c, err := NewCommander("127.0.0.1", 4222, testUsername, testPassword, true)
	require.NoError(t, err)

	tearDown := makeSSHServer()
	defer tearDown()

	err = c.Dial()
	require.NoError(t, err)

	testData := []struct{
		name string
		command string
		expectedStdOut string
		expectedStdErr string
		expectError string
	} {
		{"no error, stdout", "out", "test out",  "",""},
		{"no error, no stdout", "", "", "", ""},
		{"error from exitCode", "out code", "test out", "", "Process exited with status 16"},
		{"error from exitCode with stdErr", "out code err", "test out", "test error","Process exited with status 16"},
		{"stdErr", "out err", "test out", "test error", ""},
		{"stdErr, no stdout", "err", "", "test error", ""},
		{"error from code, with stdErr, no stdout", "err code", "", "test error", "Process exited with status 16"},
	}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			var stdOut, stdErr string
			stdOut, stdErr, err = c.RunCommandWithStdErr(tt.command)
			assert.Equal(t, tt.expectedStdOut, stdOut)
			assert.Equal(t, tt.expectedStdErr, stdErr)
			if tt.expectError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectError)
			}
		})
	}

	err = c.Close()
	assert.NoError(t, err)
}

func makeSSHServer() func() error {
	srv := &sshServer.Server{Addr: "127.0.0.1:4222", Handler: func(session sshServer.Session) {
		stdErr := ""
		stdOut := ""
		exitCode := 0

		for _, c := range session.Command() {
			switch c {
			case "err":
				stdErr = "test error"
			case "out":
				stdOut = "test out"
			case "code":
				exitCode = 16
			}
		}

		if stdErr != "" {
			_, err := session.Stderr().Write([]byte(stdErr))
			if err != nil {
				panic(err)
			}
		}

		if stdOut != "" {
			_, err := session.Write([]byte(stdOut))
			if err != nil {
				panic(err)
			}
		}

		err := session.Exit(exitCode)
		if err != nil {
			panic(err)
		}
	}}

	srv.PasswordHandler = func(ctx sshServer.Context, password string) bool {
		return ctx.User() == testUsername && password == testPassword
	}

	go srv.ListenAndServe()

	time.Sleep(time.Millisecond)
	return srv.Close
}

type publickKeyMock struct {}

func (p publickKeyMock) Type() string {
	return "mock type"
}

func (p publickKeyMock) Marshal() []byte {
	return []byte("mock marshal")
}

func (p publickKeyMock) Verify(data []byte, sig *ssh.Signature) error {
	return nil
}
