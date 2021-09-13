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
	default:
		responseData, err = nil, fmt.Errorf("invalid category '%s' provided", requestData.Category)
	}

	return responseData, err
}

func loadCpuTimeLineData(requestData *agentintegration.ServerMonitorTimeLineRequestData) (*agentintegration.ServerMonitorTimeLineResponseData, error) {
	var responseData agentintegration.ServerMonitorTimeLineResponseData
	responseData.Data = make(map[string][]agentintegration.ServerMonitorTimeLinePoint)
	filter := &service.StatProviderTimeFilter{
		FromTime: requestData.FromTime,
		ToTime:   requestData.ToTime,
	}
	overallCpuData, err := loadOverallCpuTimeLineData(filter)
	if err != nil {
		return nil, err
	}
	responseData.Data["overall"] = overallCpuData

	return &responseData, nil
}

func loadOverallCpuTimeLineData(filter service.StatProviderFilter) ([]agentintegration.ServerMonitorTimeLinePoint, error) {
	overallCpuStatCollector, err := service.GetStatCollector(&service.OverallCPUStatPrivider{})
	if err != nil {
		return nil, err
	}

	rows, err := overallCpuStatCollector.Load(filter)
	if err != nil {
		return nil, err
	}

	var overallCpuData []agentintegration.ServerMonitorTimeLinePoint
	for _, row := range rows {
		overallCpuData = append(overallCpuData, agentintegration.ServerMonitorTimeLinePoint{
			Time: row[0],
			Value: map[string]string{
				"system": row[1],
				"user":   row[2],
				"nice":   row[3],
				"idle":   row[4],
			},
		})
	}

	return overallCpuData, nil
}

func loadCoreCpuTimeLineData(filter service.StatProviderFilter) {

}
