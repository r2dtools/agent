package server

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"

	"github.com/r2dtools/agent/config"
	"github.com/r2dtools/agent/logger"
	"github.com/r2dtools/agent/modules"
	"github.com/r2dtools/agent/router"
)

// Server structure
type Server struct {
	Port     int
	listener net.Listener
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

		if err = handleConn(conn); err != nil {
			logger.Error(err.Error())
		}

		if err = conn.Close(); err != nil {
			logger.Error("could not close connection: " + err.Error())
		}
	}
}

func getResponse(data interface{}, err error) router.Response {
	var response router.Response

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
	dataLen, err := readDataLen(conn)

	if err != nil {
		return err
	}

	var data []byte
	buffer := make([]byte, 256)
	rLen := 0

	for {
		len, err := conn.Read(buffer)

		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		data = append(data, buffer[:len]...)
		rLen += len

		if rLen >= dataLen {
			break
		}
	}

	logger.Info(fmt.Sprintf("received data: %v", string(data)))
	responseData, err := handleRequest(data)
	response := getResponse(responseData, err)
	responseByte, err := json.Marshal(response)

	if err != nil {
		return fmt.Errorf("could not encode response data: %v", err)
	}

	if err = writeData(conn, responseByte); err != nil {
		return err
	}

	logger.Info("Connection successfully handled")

	return nil
}

// handleRequest handles requests from the main server
func handleRequest(data []byte) (interface{}, error) {
	var request router.Request
	err := json.Unmarshal(data, &request)

	if err != nil {
		return nil, fmt.Errorf("could not decode request data: %v", err)
	}

	if request.Token == "" || request.Token != config.GetConfig().Token {
		return nil, fmt.Errorf("invalid request token is specified: %s", request.Token)
	}

	router := &router.Router{}
	router.RegisterHandler("main", &MainHandler{})
	modules.RegisterHandlers(router)

	return router.HandleRequest(request)
}

func writeData(conn net.Conn, data []byte) error {
	// First, write sending data length to the two bytes
	header := make([]byte, 2)
	dataLen := len(data)
	binary.BigEndian.PutUint16(header, uint16(dataLen))

	if _, err := conn.Write(header); err != nil {
		return fmt.Errorf("could not write response header: %v", err)
	}

	if _, err := conn.Write(data); err != nil {
		return fmt.Errorf("could not write response data: %v", err)
	}

	return nil
}

// readDataLen reads first 2 bytes where data length is stored
func readDataLen(conn net.Conn) (int, error) {
	header := make([]byte, 2)
	_, err := conn.Read(header)

	if err != nil {
		return 0, fmt.Errorf("could not read data length: %v", err)
	}

	return int(binary.BigEndian.Uint16(header)), nil
}
