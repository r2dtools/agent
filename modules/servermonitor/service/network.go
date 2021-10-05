package service

import (
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/net"
)

// OverallNetworkStatProvider retrieves statistics data for the network usage
type OverallNetworkStatProvider struct{}

func (n *OverallNetworkStatProvider) GetData() ([]string, error) {
	iCountersStat, err := net.IOCounters(false)
	if err != nil {
		return nil, err
	}
	if len(iCountersStat) == 0 {
		return nil, nil
	}
	stat := iCountersStat[0]
	var data []string
	data = append(data, strconv.FormatInt(time.Now().Unix(), 10))
	data = append(data, formatInterfaceValue(stat.BytesRecv))
	data = append(data, formatInterfaceValue(stat.BytesSent))
	data = append(data, formatInterfaceValue(stat.PacketsRecv))
	data = append(data, formatInterfaceValue(stat.PacketsSent))

	// time|bytesrecv|bytessent|packetsrecv|packetssent
	return data, nil
}

func (n *OverallNetworkStatProvider) GetCode() string {
	return OVERALL_NETWORK_PROVIDER_CODE
}

func (n *OverallNetworkStatProvider) CheckData(data []string, filter StatProviderFilter) bool {
	if len(data) != 5 {
		return false
	}
	if filter != nil {
		return filter.Check(data)
	}

	return true
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
		name := iCounterStat.Name
		// skip docker, lo network interfsces
		if strings.HasPrefix(name, "docker") || name == "lo" {
			continue
		}
		interfaceInfo := make(map[string]string)
		interfaceInfo["name"] = name
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
