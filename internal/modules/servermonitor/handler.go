package servermonitor

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/r2dtools/agent/internal/modules/servermonitor/handler"
	"github.com/r2dtools/agent/pkg/logger"
	"github.com/r2dtools/agent/pkg/router"
	"github.com/r2dtools/agentintegration"
)

type Handler struct {
	logger logger.LoggerInterface
}

func (h *Handler) Handle(request router.Request) (interface{}, error) {
	var response interface{}
	var err error

	switch action := request.GetAction(); action {
	case "loadStatisticsData":
		response, err = h.loadStatisticsData(request.Data)
	default:
		response, err = nil, fmt.Errorf("invalid action '%s' for module '%s'", action, request.GetModule())
	}

	return response, err
}

func (h *Handler) loadStatisticsData(data interface{}) (interface{}, error) {
	var requestData agentintegration.ServerMonitorStatisticsRequestData
	err := mapstructure.Decode(data, &requestData)

	if err != nil {
		return nil, fmt.Errorf("servermonitor: invalid request data: %v", err)
	}

	var responseData interface{}

	switch requestData.Category {
	case "cpu":
		responseData, err = handler.LoadCpuTimeLineData(&requestData, h.logger)
	case "memory":
		responseData, err = handler.LoadMemoryTimeLineData(&requestData, h.logger)
	case "disk":
		responseData, err = handler.LoadDiskUsageTimeLineData(&requestData, h.logger)
	case "network":
		responseData, err = handler.LoadNetworkTimeLineData(&requestData, h.logger)
	case "process":
		responseData, err = handler.LoadProcessStatisticsData(&requestData)
	default:
		responseData, err = nil, fmt.Errorf("invalid category '%s' provided", requestData.Category)
	}

	return responseData, err
}

func GetHandler(logger logger.LoggerInterface) router.HandlerInterface {
	return &Handler{logger}
}
