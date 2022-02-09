package keeneticmanager

const ipRulesComment = "!KeeneticManager ip Rule"

const ipRouteCommand = "ip route "

const autoKeyword = "auto"

const ipRouteSuffix = " " + autoKeyword + " " + ipRulesComment

func makeIPRouteCommand(host, interfaceName string) string {
	// ip route 10.1.1.1 Wireguard0 auto !KeeneticManager ip Rule
	return ipRouteCommand + host + " " + interfaceName + ipRouteSuffix
}
