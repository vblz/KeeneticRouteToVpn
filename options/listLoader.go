package options

import (
	"bufio"
	"fmt"
	"github.com/go-pkgz/lgr"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func loadHostsList(path string) ([]string, error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("can't open file: %w", err)
	}

	result := make([]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		s, err := parseHostLine(line)
		if err != nil {
			return nil, fmt.Errorf("can't parse '%s': %w", line, err)
		}
		result = append(result, s...)
	}

	return result, nil
}

func parseHostLine(l string) ([]string, error) {
	l = trimSymbol(l, "#")
	l = strings.TrimSpace(l)
	if l == "" {
		return nil, nil
	}

	if ip := net.ParseIP(l); ip != nil {
		if !isIP4(ip) {
			return nil, fmt.Errorf("only IPv4 is supported")
		}
		return []string{l}, nil
	}

	if ip, _, err := net.ParseCIDR(l); err == nil {
		if !isIP4(ip) {
			return nil, fmt.Errorf("only IPv4 is supported")
		}
		return []string{l}, nil
	}

	ips, err := net.LookupIP(l)
	if err != nil {
		return nil, fmt.Errorf("can't resolve host '%s': %w", l, err)
	}

	result := make([]string, 0, len(ips))
	for _, ip := range ips {
		if isIP4(ip) {
			result = append(result, ip.String())
		} else {
			lgr.Printf("[WARN] IPv6 record for host '%s' was ignored", l)
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("host %s doesn't have IPv4 records", l)
	}

	return result, nil
}

func isIP4(ip net.IP) bool {
	return ip.To4() != nil
}

func trimSymbol(str, symbol string) string {
	if idx := strings.Index(str, symbol); idx != -1 {
		return str[:idx]
	}
	return str
}
