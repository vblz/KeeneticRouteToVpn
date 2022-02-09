# KeeneticRouteToVpn [![build](https://github.com/vblz/KeeneticRouteToVpn/actions/workflows/tests.yml/badge.svg)](https://github.com/vblz/KeeneticRouteToVpn/actions/workflows/tests.yml) [![Coverage Status](https://coveralls.io/repos/github/vblz/KeeneticRouteToVpn/badge.svg?branch=master)](https://coveralls.io/github/vblz/KeeneticRouteToVpn?branch=master)

KeeneticRouteToVpn is simple app updating Keenetic Router rules for some hosts to go through VPN interface.

It has defaults values and just need to be provided with a text files containing lists of hosts needed access through VPN.

Every execution the app clears all rules generated by it and add new ones from the list provided.

Password will either prompted or will read from standard input with `--password-stdin` flag.

## Usage

```shell
$ cat list.txt                                                                                                                                                                                                                                                                                                            386ms  Wed Feb  9 01:03:37 2022
# IPv4
10.12.1.5
192.168.10.44

# CIDR with IPv4
172.16.0.0/16

# HostNames
internal.company.com # only v4 IPs records will be resolved

$ ./krv list.txt
```

```
Application Options:
  -u, --username=                      username (default: admin) [$USERNAME]
  -H, --host=                          host to connect to (default: 192.168.1.1) [$HOST]
  -p, --port=                          port to connect to (default: 22) [$PORT]
  -i, --interface=                     interface name (default: Wireguard0) [$INTERFACE]
      --password-stdin                 take password from stdin [$PASSWORD_STDIN]
      --insecure-ignore-host-checking  ignore known_hosts checking [$INSECURE_IGNORE_HOST_CHECKING]

Help Options:
  -h, --help                           Show this help message

```


## Limitations
- user must have access to command line
- only ipv4 supported (Keenetic limitation)
- to provide security, you should connect via SSH to the server before and accept its public key. You can also provide `--insecure-ignore-host-checking` to ignore checking server's public key, but it is insecure. 