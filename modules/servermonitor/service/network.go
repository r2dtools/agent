package service

import (
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/net"
)

type lastCounters struct {
	lastNetworkCounters []net.IOCountersStat
	lastTime            int64
}

var lCounters lastCounters

func init() {
	lCounters.lastNetworkCounters = make([]net.IOCountersStat, 0)
	lCounters.lastTime = time.Now().Unix()
}

// OverallNetworkStatProvider retrieves statistics data for the network usage
type OverallNetworkStatProvider struct {
	BaseStatProvider
}

func (n *OverallNetworkStatProvider) GetRecord() ([]string, error) {
	iCountersStat, err := net.IOCounters(false)
	if err != nil {
		return nil, err
	}
	if len(iCountersStat) == 0 {
		return nil, nil
	}

	currentTime := time.Now().Unix()

	if len(lCounters.lastNetworkCounters) == 0 {
		lCounters.lastNetworkCounters = iCountersStat
		lCounters.lastTime = currentTime
		return nil, nil
	}

	current := iCountersStat[0]
	previous := lCounters.lastNetworkCounters[0]
	timeDelta := currentTime - lCounters.lastTime

	var data []string
	data = append(data, formatSpeed(previous.BytesRecv, current.BytesRecv, timeDelta))
	data = append(data, formatSpeed(previous.BytesSent, current.BytesSent, timeDelta))
	data = append(data, formatSpeed(previous.PacketsRecv, current.PacketsRecv, timeDelta))
	data = append(data, formatSpeed(previous.PacketsSent, current.PacketsSent, timeDelta))
	data = append(data, formatSpeed(previous.Errin, current.Errin, timeDelta))
	data = append(data, formatSpeed(previous.Errout, current.Errout, timeDelta))
	lCounters.lastNetworkCounters = iCountersStat
	lCounters.lastTime = currentTime

	// bytesrecv|bytessent|packetsrecv|packetssent|errin|errout
	return data, nil
}

func (n *OverallNetworkStatProvider) GetAverageRecord(records [][]string) []string {
	return n.getAverageRecord(records, n.GetFieldsCount(), false, n.GetEmptyRecordValue)
}

func (n *OverallNetworkStatProvider) GetFieldsCount() int {
	return 6
}

func (n *OverallNetworkStatProvider) GetCode() string {
	return OVERALL_NETWORK_PROVIDER_CODE
}

func (n *OverallNetworkStatProvider) CheckRecord(data []string, filter StatProviderFilter) bool {
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
		interfaceInfo["flags"] = strings.Join(iStat.Flags, ",")
		interfaceInfo["hadwareaddr"] = iStat.HardwareAddr
	}

	return data, nil
}

func formatSpeed(previous, current uint64, time int64) string {
	if previous > current || time <= 0 {
		return "0"
	}

	return strconv.FormatUint((current-previous)/uint64(time), 10)
}
