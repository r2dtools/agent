package server

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"

	"github.com/r2dtools/agent/config"
	"github.com/r2dtools/agent/internal/pkg/logger"
	"github.com/r2dtools/agent/internal/pkg/router"
	"github.com/r2dtools/agent/internal/pkg/service"
)

const HEADER_DATA_LENGTH = 4 // bytes

// Server structure
type Server struct {
	Port           int
	ServiceManager service.ServiceManager
	Router         router.Router
	Logger         logger.Logger
	Config         *config.Config
	listener       net.Listener
}

// Serve starts TCP server
func (s *Server) Serve() error {
	port := strconv.Itoa(s.Port)
	s.Logger.Info("starting TCP server on port %s ...", port)
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(s.Port))

	if err != nil {
		s.Logger.Error("error starting TCP server: %v", err)
		return err
	}

	s.listener = listener
	s.Logger.Info("TCP server successfully started")
	defer listener.Close()

	s.ServiceManager.RunServices()

	for {
		s.Logger.Info("listening to a remote conection ...")
		conn, err := listener.Accept()

		if err != nil {
			s.Logger.Error("error accepting remote connection: %v", err)
			continue
		}

		s.Logger.Info("accepted connection from the remote address: %v", conn.RemoteAddr())
		s.handleConn(conn)
	}
}

func prepareResponse(data interface{}, err error) router.Response {
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

func (s *Server) getResponse(reader io.Reader) router.Response {
	dataLen, err := s.readDataLen(reader)

	if err != nil {
		return prepareResponse(nil, err)
	}

	var data []byte
	buffer := make([]byte, 256)
	rLen := 0

	for {
		len, err := reader.Read(buffer)

		if err != nil {
			if err == io.EOF {
				break
			}

			return prepareResponse(nil, err)
		}

		data = append(data, buffer[:len]...)
		rLen += len

		if rLen >= dataLen {
			break
		}
	}

	s.Logger.Info("received data: %v", string(data))
	responseData, err := s.handleRequest(data)

	return prepareResponse(responseData, err)
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	var response router.Response

	response = s.getResponse(conn)

	if response.Error != "" {
		s.Logger.Error(response.Error)
	}

	responseByte, err := json.Marshal(response)
	if err != nil {
		response = prepareResponse(nil, fmt.Errorf("could not encode response data: %v", err))
		responseByte, _ = json.Marshal(response)
		s.Logger.Error(response.Error)
	}

	if err = s.writeData(conn, responseByte); err != nil {
		s.Logger.Error(err.Error())
		return
	}

	s.Logger.Info("Connection successfully handled")
}

// handleRequest handles requests from the main server
func (s *Server) handleRequest(data []byte) (interface{}, error) {
	var request router.Request
	err := json.Unmarshal(data, &request)

	if err != nil {
		return nil, fmt.Errorf("could not decode request data: %v", err)
	}

	if request.Token == "" || request.Token != s.Config.Token {
		return nil, fmt.Errorf("invalid request token is specified: %s", request.Token)
	}

	return s.Router.HandleRequest(request)
}

func (s *Server) writeData(writer io.Writer, data []byte) error {
	// First, write sending data length
	header := make([]byte, HEADER_DATA_LENGTH)
	dataLen := len(data)
	binary.BigEndian.PutUint32(header, uint32(dataLen))

	if _, err := writer.Write(header); err != nil {
		return fmt.Errorf("could not write response header: %v", err)
	}

	if _, err := writer.Write(data); err != nil {
		return fmt.Errorf("could not write response data: %v", err)
	}

	return nil
}

// readDataLen reads first bytes where data length is stored
func (s *Server) readDataLen(reader io.Reader) (int, error) {
	header := make([]byte, HEADER_DATA_LENGTH)
	if _, err := reader.Read(header); err != nil {
		return 0, fmt.Errorf("could not read data length: %v", err)
	}

	return int(binary.BigEndian.Uint32(header)), nil
}
