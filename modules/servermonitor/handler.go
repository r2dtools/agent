package servermonitor

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
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

	var responseData agentintegration.ServerMonitorTimeLineResponseData
	responseData.Data = make(map[string][]agentintegration.ServerMonitorTimeLinePoint)
	responseData.Data["overall"] = []agentintegration.ServerMonitorTimeLinePoint{
		{Time: 1, Value: map[string]float32{"user": 1, "system": 2}},
		{Time: 2, Value: map[string]float32{"user": 2, "system": 4}},
		{Time: 3, Value: map[string]float32{"user": 3, "system": 9}},
		{Time: 5, Value: map[string]float32{"user": 7, "system": 1}},
	}

	return &responseData, nil
}
