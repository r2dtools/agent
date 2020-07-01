package server

import (
	"encoding/json"
	"fmt"
	// "github.com/r2dtools/agent/config"
)

// Request is a request object from the main server
type Request struct {
	Command,
	Token string
	Data interface{}
}

func register(data interface{}) (interface{}, error) {
	var respData struct {
		OsCode,
		OsVersion,
		AgentVersion string
	}
	respData.AgentVersion = "1.0.0"
	respData.OsCode = "ubuntu"
	respData.OsVersion = "18.04"

	return respData, nil
}

// HandleRequest handles requests from the main server
func HandleRequest(data []byte) (interface{}, error) {
	var request Request
	err := json.Unmarshal(data, &request)

	if err != nil {
		return nil, fmt.Errorf("could not decode request data: %v", err)
	}

	//if request.Token == "" || request.Token != config.GetConfig().Token {
	//	return nil, fmt.Errorf("invalid token specified: %s", request.Token)
	//}

	var response interface{}

	switch command := request.Command; command {
	case "register":
		response, err = register(request.Data)
	default:
		response, err = nil, fmt.Errorf("unsupported command: %s", request.Command)
	}

	return response, nil
}
