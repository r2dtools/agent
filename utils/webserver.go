package utils

import "github.com/r2dtools/agentintegration"

// FilterVhosts filter vhosts
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
