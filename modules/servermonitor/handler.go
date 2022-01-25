package servermonitor

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/r2dtools/agent/modules/servermonitor/handler"
	"github.com/r2dtools/agent/router"
	"github.com/r2dtools/agentintegration"
)

const MODULE_ID = "servermonitor"

// Handler handles requests to the module
type Handler struct{}

// Handle handles request to the module
func (h *Handler) Handle(request router.Request) (interface{}, error) {
	var response interface{}
	var err error

	switch action := request.GetAction(); action {
	case "loadStatisticsData":
		response, err = loadStatisticsData(request.Data)
	default:
		response, err = nil, fmt.Errorf("invalid action '%s' for module '%s'", action, request.GetModule())
	}

	return response, err
}

func loadStatisticsData(data interface{}) (interface{}, error) {
	var requestData agentintegration.ServerMonitorStatisticsRequestData
	err := mapstructure.Decode(data, &requestData)

	if err != nil {
		return nil, fmt.Errorf("servermonitor: invalid request data: %v", err)
	}

	var responseData interface{}

	switch requestData.Category {
	case "cpu":
		responseData, err = handler.LoadCpuTimeLineData(&requestData)
	case "memory":
		responseData, err = handler.LoadMemoryTimeLineData(&requestData)
	case "disk":
		responseData, err = handler.LoadDiskUsageTimeLineData(&requestData)
	case "network":
		responseData, err = handler.LoadNetworkTimeLineData(&requestData)
	case "process":
		responseData, err = handler.LoadProcessStatisticsData(&requestData)
	default:
		responseData, err = nil, fmt.Errorf("invalid category '%s' provided", requestData.Category)
	}

	return responseData, err
}
