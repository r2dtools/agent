package utils

import (
	"github.com/google/go-cmp/cmp"
	"github.com/unknwon/com"

	"github.com/r2dtools/agentintegration"
)

func FilterVhosts(vhosts []agentintegration.VirtualHost) []agentintegration.VirtualHost {
	var fVhosts []agentintegration.VirtualHost

	for _, vhost := range vhosts {
		if vhost.ServerName == "" || !checkVhostPorts(vhost.Addresses, []string{"80", "443"}) {
			continue
		}

		fVhosts = append(fVhosts, vhost)
	}

	return fVhosts
}

// MergeVhosts merge similar vhosts. For example, vhost:443 will be merged with vhost:80
func MergeVhosts(vhosts []agentintegration.VirtualHost) []agentintegration.VirtualHost {
	var fVhosts []agentintegration.VirtualHost
	vhostsMap := make(map[string]agentintegration.VirtualHost)

	for _, vhost := range vhosts {
		if existedVhost, ok := vhostsMap[vhost.ServerName]; ok {
			existedVhost.Ssl = existedVhost.Ssl || vhost.Ssl

			if existedVhost.DocRoot == "" {
				existedVhost.DocRoot = vhost.DocRoot
			}

			// merge addresses (for example ipv4 + ipv6)
			for _, address := range vhost.Addresses {
				var addressExists bool
				for _, eAddress := range existedVhost.Addresses {
					if cmp.Equal(address, eAddress) {
						addressExists = true
						break
					}
				}

				if !addressExists {
					existedVhost.Addresses = append(existedVhost.Addresses, address)
				}
			}

			// merge aliases
			for _, alias := range vhost.Aliases {
				existedVhost.Aliases = com.AppendStr(existedVhost.Aliases, alias)
			}

			vhostsMap[vhost.ServerName] = existedVhost
		} else {
			vhostsMap[vhost.ServerName] = vhost
		}
	}

	for _, vhost := range vhostsMap {
		fVhosts = append(fVhosts, vhost)
	}

	return fVhosts
}

func checkVhostPorts(addresses []agentintegration.VirtualHostAddress, ports []string) bool {
	for _, address := range addresses {
		for _, port := range ports {
			if port == address.Port {
				return true
			}
		}
	}

	return false
}
