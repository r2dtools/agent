package servermonitor

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/r2dtools/agent/modules/servermonitor/service"
	"github.com/r2dtools/agent/router"
	"github.com/r2dtools/agentintegration"
)

// Handler handles requests to the module
type Handler struct{}

// Handle handles request to the module
func (h *Handler) Handle(request router.Request) (interface{}, error) {
	var response interface{}
	var err error

	switch action := request.GetAction(); action {
	case "loadTimeLineData":
		response, err = loadTimeLineData(request.Data)
	default:
		response, err = nil, fmt.Errorf("invalid action '%s' for module '%s'", action, request.GetModule())
	}

	return response, err
}

func loadTimeLineData(data interface{}) (*agentintegration.ServerMonitorTimeLineResponseData, error) {
	var requestData agentintegration.ServerMonitorTimeLineRequestData
	err := mapstructure.Decode(data, &requestData)

	if err != nil {
		return nil, fmt.Errorf("servermonitor: invalid request data: %v", err)
	}

	var responseData *agentintegration.ServerMonitorTimeLineResponseData

	switch requestData.Category {
	case "cpu":
		responseData, err = loadCpuTimeLineData(&requestData)
	case "memory":
		responseData, err = loadMemoryTimeLineData(&requestData)
	case "disk":
		responseData, err = loadDiskUsageTimeLineData(&requestData)
	default:
		responseData, err = nil, fmt.Errorf("invalid category '%s' provided", requestData.Category)
	}

	return responseData, err
}

func loadDiskUsageTimeLineData(requestData *agentintegration.ServerMonitorTimeLineRequestData) (*agentintegration.ServerMonitorTimeLineResponseData, error) {
	return nil, nil
}

func loadCpuTimeLineData(requestData *agentintegration.ServerMonitorTimeLineRequestData) (*agentintegration.ServerMonitorTimeLineResponseData, error) {
	var responseData agentintegration.ServerMonitorTimeLineResponseData
	responseData.Data = make(map[string][]agentintegration.ServerMonitorTimeLinePoint)
	filter := &service.StatProviderTimeFilter{
		FromTime: requestData.FromTime,
		ToTime:   requestData.ToTime,
	}
	if err := loadOverallCpuTimeLineData(&responseData, filter); err != nil {
		return nil, err
	}
	if err := loadCoreCpuTimeLineData(&responseData, filter); err != nil {
		return nil, err
	}

	return &responseData, nil
}

func loadMemoryTimeLineData(requestData *agentintegration.ServerMonitorTimeLineRequestData) (*agentintegration.ServerMonitorTimeLineResponseData, error) {
	var responseData agentintegration.ServerMonitorTimeLineResponseData
	responseData.Data = make(map[string][]agentintegration.ServerMonitorTimeLinePoint)
	filter := &service.StatProviderTimeFilter{
		FromTime: requestData.FromTime,
		ToTime:   requestData.ToTime,
	}

	if err := loadVirtualMemoryTimeLineData(&responseData, filter); err != nil {
		return nil, err
	}
	if err := loadSwapMemoryTimeLineData(&responseData, filter); err != nil {
		return nil, err
	}

	return &responseData, nil
}

func loadVirtualMemoryTimeLineData(responseData *agentintegration.ServerMonitorTimeLineResponseData, filter service.StatProviderFilter) error {
	virtualMemoryStatCollector, err := service.GetStatCollector(&service.VirtualMemoryStatPrivider{})
	if err != nil {
		return nil
	}

	rows, err := virtualMemoryStatCollector.Load(filter)
	if err != nil {
		return nil
	}

	var memoryData []agentintegration.ServerMonitorTimeLinePoint
	for _, row := range rows {
		memoryData = append(memoryData, getVirtualMemoryTimeLinePoint(row))
	}
	responseData.Data["virtual"] = memoryData

	return nil
}

func loadSwapMemoryTimeLineData(responseData *agentintegration.ServerMonitorTimeLineResponseData, filter service.StatProviderFilter) error {
	swapMemoryStatCollector, err := service.GetStatCollector(&service.SwapMemoryStatPrivider{})
	if err != nil {
		return nil
	}

	rows, err := swapMemoryStatCollector.Load(filter)
	if err != nil {
		return nil
	}

	var memoryData []agentintegration.ServerMonitorTimeLinePoint
	for _, row := range rows {
		memoryData = append(memoryData, getSwapMemoryTimeLinePoint(row))
	}
	responseData.Data["swap"] = memoryData

	return nil
}

func loadOverallCpuTimeLineData(responseData *agentintegration.ServerMonitorTimeLineResponseData, filter service.StatProviderFilter) error {
	overallCpuStatCollector, err := service.GetStatCollector(&service.OverallCPUStatPrivider{})
	if err != nil {
		return err
	}

	rows, err := overallCpuStatCollector.Load(filter)
	if err != nil {
		return err
	}

	var overallCpuData []agentintegration.ServerMonitorTimeLinePoint
	for _, row := range rows {
		overallCpuData = append(overallCpuData, getCpuTimeLinePoint(row))
	}
	responseData.Data["overall"] = overallCpuData

	return nil
}

func loadCoreCpuTimeLineData(responseData *agentintegration.ServerMonitorTimeLineResponseData, filter service.StatProviderFilter) error {
	coreCpuStatCollectors, err := service.GetCoreCpuStatCollectors()
	if err != nil {
		return err
	}

	for index, collector := range coreCpuStatCollectors {
		var coreCpuData []agentintegration.ServerMonitorTimeLinePoint
		rows, err := collector.Load(filter)
		if err != nil {
			return err
		}

		for _, row := range rows {
			coreCpuData = append(coreCpuData, getCpuTimeLinePoint(row))
		}
		responseData.Data[fmt.Sprintf("core%d", index+1)] = coreCpuData
	}

	return nil
}

func getCpuTimeLinePoint(row []string) agentintegration.ServerMonitorTimeLinePoint {
	return agentintegration.ServerMonitorTimeLinePoint{
		Time: row[0],
		Value: map[string]string{
			"system": row[1],
			"user":   row[2],
			"nice":   row[3],
			"idle":   row[4],
		},
	}
}

// time|total|available|free|used|active|inactive|cached|buffered
func getVirtualMemoryTimeLinePoint(row []string) agentintegration.ServerMonitorTimeLinePoint {
	return agentintegration.ServerMonitorTimeLinePoint{
		Time: row[0],
		Value: map[string]string{
			"total":     row[1],
			"available": row[2],
			"free":      row[3],
			"used":      row[4],
			"active":    row[5],
			"inactive":  row[6],
			"cached":    row[7],
			"buffered":  row[8],
		},
	}
}

// time|total|used|free
func getSwapMemoryTimeLinePoint(row []string) agentintegration.ServerMonitorTimeLinePoint {
	return agentintegration.ServerMonitorTimeLinePoint{
		Time: row[0],
		Value: map[string]string{
			"total": row[1],
			"used":  row[2],
			"free":  row[3],
		},
	}
}
