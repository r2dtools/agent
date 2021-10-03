package handler

import (
	"github.com/r2dtools/agent/modules/servermonitor/service"
	"github.com/r2dtools/agentintegration"
)

func LoadNetworkTimeLineData(requestData *agentintegration.ServerMonitorTimeLineRequestData) (*agentintegration.ServerMonitorNetworkResponseData, error) {
	var responseData agentintegration.ServerMonitorNetworkResponseData
	responseData.TimeLineData = make(map[string][]agentintegration.ServerMonitorTimeLinePoint)
	responseData.InterfacesInfo = make([]map[string]string, 0)

	if err := loadNetworkInterfaceInfo(&responseData); err != nil {
		return nil, err
	}

	return &responseData, nil
}

func loadNetworkInterfaceInfo(responseData *agentintegration.ServerMonitorNetworkResponseData) error {
	interfacesInfo, err := service.GetNetworkInterfacesInfo()
	if err != nil {
		return err
	}

	responseData.InterfacesInfo = interfacesInfo
	return nil
}
