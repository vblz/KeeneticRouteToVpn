package options

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/term"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

var tempList = []string{
	"\t",
	"10.0.0.1\r",
	"",
	"192.168.1.1 # comment",
	"4.4.0.0/24",
	"# comment",
	"\r   localhost",
	"",
}

var cleanedList = []string{
	"10.0.0.1",
	"192.168.1.1",
	"4.4.0.0/24",
	"127.0.0.1",
}

func TestParse_Defaults(t *testing.T) {
	hostFileName := tempHostsFile(t)
	defer os.Remove(hostFileName)

	args := []string{
		"mockExecutableName",
		"--password-stdin",
		hostFileName,
	}

	tearDown := mockStdinWithData(t, "testPassword12312")
	defer tearDown()

	opts, err := Parse(args)
	require.NoError(t, err)

	assert.Equal(t, uint16(22), opts.Port)
	assert.Equal(t, "192.168.1.1", opts.Host)
	assert.Equal(t, "testPassword12312", opts.Password)
	assert.Equal(t, "admin", opts.Username)
	assert.Equal(t, "Wireguard0", opts.InterfaceName)
	assert.Equal(t, cleanedList, opts.HostsList)
	assert.False(t, opts.IgnoreHostChecking)
}

func TestParse_Full(t *testing.T) {
	hostFileName := tempHostsFile(t)
	defer os.Remove(hostFileName)

	args := []string{
		"mockExecutableName",
		"-u",
		"testDevelopment",
		"-p",
		"2222",
		"-H",
		"10.0.0.1",
		"-i",
		"OpenVPN0",
		"--password-stdin",
		"--insecure-ignore-host-checking",
		hostFileName,
	}

	tearDown := mockStdinWithData(t, "testPassword")
	defer tearDown()

	opts, err := Parse(args)
	require.NoError(t, err)

	assert.Equal(t, uint16(2222), opts.Port)
	assert.Equal(t, "10.0.0.1", opts.Host)
	assert.Equal(t, "testPassword", opts.Password)
	assert.Equal(t, "testDevelopment", opts.Username)
	assert.Equal(t, "OpenVPN0", opts.InterfaceName)
	assert.Equal(t, cleanedList, opts.HostsList)
	assert.True(t, opts.IgnoreHostChecking)
}

func TestParse_FilePath(t *testing.T) {
	args := []string {
		"mockExecutableName",
	}

	tearDown := mockStdinWithData(t, "Paasswwooorrd")
	defer tearDown()

	_, err := Parse(args)
	require.Error(t, err)
}

func TestParse_FileNotExist(t *testing.T) {
	args := []string{
		"mockExecutableName",
		"not_really_file.txt",
	}

	tearDown := mockStdinWithData(t, "testPassword")
	defer tearDown()

	_, err := Parse(args)
	require.Error(t, err)
}

func TestSecureReadPasswordFunc(t *testing.T) {
	assert.Equal(t,
		reflect.ValueOf(secureReadPasswordFunc).Pointer(),
		reflect.ValueOf(term.ReadPassword).Pointer())
}

func TestReadPassword(t *testing.T) {
	const password = "testPassword"

	oldValue := secureReadPasswordFunc
	defer func() {
		secureReadPasswordFunc = oldValue
	}()

	calls := 0
	secureReadPasswordFunc = func(fd int) ([]byte, error) {
		assert.Equal(t, fd, 0)
		calls++
		return []byte(password), nil
	}

	hostFileName := tempHostsFile(t)
	defer os.Remove(hostFileName)

	args := []string{
		"mockExecutableName",
		hostFileName,
	}

	opts, err := Parse(args)
	assert.Equal(t, password, opts.Password)
	assert.NoError(t, err)
	assert.Equal(t, 1, calls)
}

// mockStdinWithData changes stdin for file with data. Returns a tear down func, restoring stdIn;
func mockStdinWithData(t *testing.T, data string) func() {
	original := os.Stdin

	read, write, err := os.Pipe()
	require.NoError(t, err)
	_, err = write.WriteString(data)
	require.NoError(t, err)
	err = write.Close()
	require.NoError(t, err)

	os.Stdin = read

	return func() {
		os.Stdin = original
	}
}

// tempHostsFile makes a temp file with list of hosts.
// Returns full path to the file. File is closed. It is the caller's responsibility
// to remove the file when no longer needed.
func tempHostsFile(t *testing.T) string {
	f, err := ioutil.TempFile("", "")
	require.NoError(t, err)

	_, err = f.WriteString(strings.Join(tempList, "\n"))
	require.NoError(t, err)

	err = f.Close()
	require.NoError(t, err)

	return f.Name()
}
