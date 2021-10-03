package service

import (
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/net"
)

// NetworkStatProvider retrieves statistics data for the network usage
type NetworkStatProvider struct{}

func (n *NetworkStatProvider) GetData() ([]string, error) {
	return nil, nil
}

func GetNetworkInterfacesInfo() ([]map[string]string, error) {
	iCountersStat, err := net.IOCounters(true)
	if err != nil {
		return nil, err
	}

	iStats, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	iMap := make(map[string]net.InterfaceStat)
	for _, iStat := range iStats {
		iMap[iStat.Name] = iStat
	}

	var data []map[string]string
	for _, iCounterStat := range iCountersStat {
		interfaceInfo := make(map[string]string)
		interfaceInfo["name"] = iCounterStat.Name
		interfaceInfo["bytesrecv"] = formatInterfaceValue(iCounterStat.BytesRecv)
		interfaceInfo["bytessent"] = formatInterfaceValue(iCounterStat.BytesSent)
		interfaceInfo["packetsrecv"] = formatInterfaceValue(iCounterStat.PacketsRecv)
		interfaceInfo["packetssent"] = formatInterfaceValue(iCounterStat.PacketsSent)
		interfaceInfo["errin"] = formatInterfaceValue(iCounterStat.Errin)
		interfaceInfo["errout"] = formatInterfaceValue(iCounterStat.Errout)

		iStat, ok := iMap[iCounterStat.Name]
		if !ok || len(iStat.Addrs) == 0 {
			continue
		}

		interfaceInfo["mtu"] = strconv.Itoa(iStat.MTU)
		var addresses []string
		for _, addr := range iStat.Addrs {
			addresses = append(addresses, addr.Addr)
		}

		interfaceInfo["addresses"] = strings.Join(addresses, ",")
		data = append(data, interfaceInfo)
	}

	return data, nil
}

func formatInterfaceValue(value uint64) string {
	return strconv.FormatUint(value, 10)
}
