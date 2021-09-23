package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

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

	return prepareData(items[0])
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

	return prepareData(items[sc.CoreNum-1])
}

func (sc *CoreCPUStatPrivider) GetCode() string {
	return fmt.Sprintf("%s%d", CORE_CPU_PROVIDER_CODE, sc.CoreNum)
}

func (sc *CoreCPUStatPrivider) CheckData(data []string, filter StatProviderFilter) bool {
	return checkData(data, filter)
}

func prepareData(item cpu.TimesStat) ([]string, error) {
	var data []string
	total := item.Total()
	data = append(data, strconv.FormatInt(time.Now().Unix(), 10))
	data = append(data, getCpuPercertValue(item.System, total))
	data = append(data, getCpuPercertValue(item.User, total))
	data = append(data, getCpuPercertValue(item.Nice, total))
	data = append(data, getCpuPercertValue(item.Idle, total))

	// data: time|system|user|nice|idle
	return data, nil
}

func getCpuPercertValue(value, total float64) string {
	return fmt.Sprintf("%.2f", 100*value/total)
}

func checkData(data []string, filter StatProviderFilter) bool {
	if len(data) != 5 {
		return false
	}

	if filter != nil {
		return filter.Check(data)
	}

	return true
}