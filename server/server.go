package server

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/r2dtools/agent/logger"
)

// Server structure
type Server struct {
	Port     int
	listener net.Listener
}

// Response that will be sent to the mail server
type Response struct {
	Status,
	Error string
	Data interface{}
}

// Handler is an interface that must be implemented by any action handler
type Handler interface {
	Handle(request Request) (interface{}, error)
}

// Serve starts TCP server
func (s *Server) Serve() error {
	port := strconv.Itoa(s.Port)
	logger.Info(fmt.Sprintf("starting TCP server on port %s ...", port))
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(s.Port))

	if err != nil {
		logger.Error(fmt.Sprintf("error starting TCP server: %v", err))
		return err
	}

	s.listener = listener
	logger.Info("TCP server successfully started")
	defer listener.Close()

	for {
		logger.Info(fmt.Sprintf("listening to a remote conection ..."))
		conn, err := listener.Accept()

		if err != nil {
			logger.Error(fmt.Sprintf("error accepting remote connection: %v", err))
			continue
		}

		logger.Info(fmt.Sprintf("accepted connection from the remote address: %v", conn.RemoteAddr()))
		err = handleConn(conn)

		if err != nil {
			logger.Error(err.Error())
		}
	}
}

func getResponse(data interface{}, err error) Response {
	var response Response

	if err != nil {
		response.Status = "error"
		response.Error = err.Error()
	} else {
		response.Status = "ok"
		response.Data = data
	}

	return response
}

func handleConn(conn net.Conn) error {
	//defer conn.Close()
	buffer := make([]byte, 1024)
	len, err := conn.Read(buffer)

	if err != nil {
		return err
	}

	data := buffer[:len]
	logger.Info(fmt.Sprintf("received data: %v", string(data)))
	responseData, err := handleRequest(data)
	response := getResponse(responseData, err)
	responseByte, err := json.Marshal(response)

	if err != nil {
		return fmt.Errorf("could not encode response data: %v", err)
	}

	_, err = conn.Write(responseByte)

	if err != nil {
		return fmt.Errorf("could not send response: %v", err)
	}

	logger.Info("Connection successfully handled")

	return nil
}

// HandleRequest handles requests from the main server
func handleRequest(data []byte) (interface{}, error) {
	var request Request
	err := json.Unmarshal(data, &request)

	if err != nil {
		return nil, fmt.Errorf("could not decode request data: %v", err)
	}

	//if request.Token == "" || request.Token != config.GetConfig().Token {
	//	return nil, fmt.Errorf("invalid token specified: %s", request.Token)
	//}

	handler, err := getHadler(request)

	if err != nil {
		return nil, err
	}

	return handler.Handle(request)
}

func getHadler(request Request) (Handler, error) {
	var handler Handler
	var err error

	switch module := request.GetModule(); module {
	case "":
		handler = &MainHandler{}
	default:
		err = fmt.Errorf("could not find handler for command '%s'", request.Command)
	}

	return handler, err
}
