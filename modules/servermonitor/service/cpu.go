package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

type lastTimes struct {
	lastCPUTimes    []cpu.TimesStat
	lastPerCPUTimes []cpu.TimesStat
}

var lTimes lastTimes

func init() {
	lTimes.lastCPUTimes, _ = cpu.Times(false)
	lTimes.lastPerCPUTimes, _ = cpu.Times(true)
}

// OverallCPUStatPrivider retrieves overall statistics data for cpu
type OverallCPUStatPrivider struct{}

func (sc *OverallCPUStatPrivider) GetData() ([]string, error) {
	items, err := cpu.Times(false)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, nil
	}

	if len(lTimes.lastCPUTimes) == 0 {
		lTimes.lastCPUTimes = items
		return nil, nil
	}

	data, err := prepareCpuData(lTimes.lastCPUTimes[0], items[0])
	lTimes.lastCPUTimes = items

	return data, err
}

func (sc *OverallCPUStatPrivider) GetCode() string {
	return OVERALL_CPU_PROVIDER_CODE
}

func (sc *OverallCPUStatPrivider) CheckData(data []string, filter StatProviderFilter) bool {
	return checkData(data, filter)
}

// CoreCPUStatPrivider retrieves statistics data for cpu per core
type CoreCPUStatPrivider struct {
	CoreNum int
}

func (sc *CoreCPUStatPrivider) GetData() ([]string, error) {
	items, err := cpu.Times(true)
	if err != nil {
		return nil, err
	}

	if len(items) < sc.CoreNum {
		return nil, nil
	}

	if len(lTimes.lastPerCPUTimes) < sc.CoreNum {
		lTimes.lastPerCPUTimes[sc.CoreNum-1] = items[sc.CoreNum-1]
		return nil, nil
	}

	data, err := prepareCpuData(lTimes.lastPerCPUTimes[sc.CoreNum-1], items[sc.CoreNum-1])
	lTimes.lastPerCPUTimes[sc.CoreNum-1] = items[sc.CoreNum-1]

	return data, err
}

func (sc *CoreCPUStatPrivider) GetCode() string {
	return fmt.Sprintf("%s%d", CORE_CPU_PROVIDER_CODE, sc.CoreNum)
}

func (sc *CoreCPUStatPrivider) CheckData(data []string, filter StatProviderFilter) bool {
	return checkData(data, filter)
}

func prepareCpuData(previous, current cpu.TimesStat) ([]string, error) {
	var data []string

	system := current.System - previous.System
	user := current.User - previous.User
	nice := current.Nice - previous.Nice
	idle := current.Idle - previous.Idle
	total := current.Total() - previous.Total()
	usage := (current.Total() - current.Idle) - (previous.Total() - previous.Idle)

	data = append(data, strconv.FormatInt(time.Now().Unix(), 10))
	data = append(data, getCpuPercertValue(system, total))
	data = append(data, getCpuPercertValue(user, total))
	data = append(data, getCpuPercertValue(nice, total))
	data = append(data, getCpuPercertValue(idle, total))
	data = append(data, getCpuPercertValue(usage, total))

	// data: time|system|user|nice|idle|usage
	return data, nil
}

func getCpuPercertValue(value, total float64) string {
	if value <= 0 {
		return "0"
	}
	if total <= 0 {
		return "100"
	}

	return fmt.Sprintf("%.2f", 100*value/total)
}

func checkData(data []string, filter StatProviderFilter) bool {
	if len(data) < 5 {
		return false
	}

	if filter != nil {
		return filter.Check(data)
	}

	return true
}
