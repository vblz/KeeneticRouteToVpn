package options

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"golang.org/x/term"
	"io"
	"os"
	"strings"
)

const stdInFileDescriptor = 0
var secureReadPasswordFunc = term.ReadPassword

type options struct {
	Username  string `default:"admin" short:"u" long:"username" env:"USERNAME" description:"username"`
	Host      string `default:"192.168.1.1" short:"H" long:"host" env:"HOST" description:"host to connect to"`
	Port      uint16 `default:"22" short:"p" long:"port" env:"PORT" description:"port to connect to"`
	Interface string `default:"Wireguard0" short:"i" long:"interface" env:"INTERFACE" description:"interface name"`

	PasswordFromStdin bool `long:"password-stdin" env:"PASSWORD_STDIN" description:"take password from stdin"`

	IgnoreHostChecking bool `long:"insecure-ignore-host-checking" env:"INSECURE_IGNORE_HOST_CHECKING" description:"ignore known_hosts checking"`

	Args struct {
		PathToList string
	} `positional-args:"yes" required:"yes" description:"path to file with the list of hosts"`
}

// ParsedOptions contains all options prepared and parsed
type ParsedOptions struct {
	Username      string
	Password      string
	Host          string
	Port          uint16
	InterfaceName string

	IgnoreHostChecking bool

	HostsList []string
}

// Parse options from passed args, prompts for password if needed, and reads hosts file contents
func Parse(args []string) (ParsedOptions, error) {
	args = args[1:]
	var opts options
	if _, err := flags.ParseArgs(&opts, args); err != nil {
		if err.(*flags.Error).Type == flags.ErrHelp {
			os.Exit(0)
		}
		return ParsedOptions{}, fmt.Errorf("can't parse options: %w", err)
	}

	list, err := loadHostsList(opts.Args.PathToList)
	if err != nil {
		return ParsedOptions{}, fmt.Errorf("can't load list of hosts: %w", err)
	}

	password, err := getPassword(opts.PasswordFromStdin)
	if err != nil {
		return ParsedOptions{}, err
	}

	return ParsedOptions{
		Username:      opts.Username,
		Password:      password,
		Host:          opts.Host,
		Port:          opts.Port,
		InterfaceName: opts.Interface,

		IgnoreHostChecking: opts.IgnoreHostChecking,

		HostsList: list,
	}, nil
}

func getPassword(fromStdin bool) (string, error) {
	if fromStdin {
		contents, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("can't read password from stdin: %w", err)
		}

		return strings.TrimSpace(string(contents)), nil
	}

	fmt.Print("Password:")
	password, err := secureReadPasswordFunc(stdInFileDescriptor)
	fmt.Println()
	if err != nil {
		return "", fmt.Errorf("can't read password from prompt: %w", err)
	}

	return string(password), nil
}
