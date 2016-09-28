package utils

import "bufio"
import "strings"
import "os"

// A map from IP to FQDN
type Hosts map[string]string

// Ok, it's a horrible hack... but I don't know yet how to propagated
// this variable from Arguments to Printf!
var HostMap Hosts

// Take a hosts file in the same format than /etc/hosts file.
// Currently, the only two first elements are used.
func LoadHostsFile(filepath string) (Hosts, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	m := make(Hosts)
	for scanner.Scan() {
		line := scanner.Text()
		ips := strings.Split(line, " ")
		m[ips[0]] = ips[1]
	}
	return m, nil
}

func ResolveIp(ip string) string {
	fqdn, ok := HostMap[ip]
	if ok {
		return fqdn
	} else {
		return ip
	}
}
