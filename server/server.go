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
		err = handle(conn)

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

func handle(conn net.Conn) error {
	//defer conn.Close()
	buffer := make([]byte, 1024)
	len, err := conn.Read(buffer)

	if err != nil {
		return err
	}

	data := buffer[:len]
	logger.Info(fmt.Sprintf("received data: %v", string(data)))
	responseData, err := HandleRequest(data)
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
