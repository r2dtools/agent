package handler

import (
	"github.com/r2dtools/agent/modules/servermonitor/service"
	"github.com/r2dtools/agentintegration"
)

func LoadProcessStatisticsData(requestData *agentintegration.ServerMonitorStatisticsRequestData) (*agentintegration.ServerMonitorProcessResponseData, error) {
	var responseData agentintegration.ServerMonitorProcessResponseData
	processesData, err := service.GetProcessesData()
	if err != nil {
		return nil, err
	}
	responseData.ProcessesData = processesData

	return &responseData, nil
}
