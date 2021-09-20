package server

import (
	"compress/flate"
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
	"github.com/r2dtools/agent/service"
)

const HEADER_DATA_LENGTH = 4 // bytes

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

	serviceManager := service.ServiceManager{}
	if err = modules.RegisterSercices(&serviceManager); err != nil {
		return err
	}

	if err = serviceManager.RunServices(); err != nil {
		logger.Error(err.Error())
		return err
	}

	for {
		logger.Info("listening to a remote conection ...")
		conn, err := listener.Accept()

		if err != nil {
			logger.Error(fmt.Sprintf("error accepting remote connection: %v", err))
			continue
		}

		logger.Info(fmt.Sprintf("accepted connection from the remote address: %v", conn.RemoteAddr()))
		go handleConn(conn)
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

func getResponse(reader io.Reader) router.Response {
	dataLen, err := readDataLen(reader)
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

	logger.Info(fmt.Sprintf("received data: %v", string(data)))
	responseData, err := handleRequest(data)

	return prepareResponse(responseData, err)
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	var response router.Response

	cReader := flate.NewReader(conn)
	defer cReader.Close()
	response = getResponse(cReader)

	if response.Error != "" {
		logger.Error(response.Error)
	}

	cWriter, err := flate.NewWriter(conn, flate.BestCompression)
	if err != nil {
		logger.Error(fmt.Sprintf("could not create compressed data writer: %v", err))
		return
	}
	defer cWriter.Close()

	jEncoder := json.NewEncoder(cWriter)
	if err = jEncoder.Encode(response); err != nil {
		logger.Error(err.Error())
		return
	}
	cWriter.Flush()
	logger.Info("Connection successfully handled")
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

// readDataLen reads first bytes where data length is stored
func readDataLen(reader io.Reader) (int, error) {
	header := make([]byte, HEADER_DATA_LENGTH)
	if _, err := reader.Read(header); err != nil {
		return 0, fmt.Errorf("could not read data length: %v", err)
	}

	return int(binary.BigEndian.Uint32(header)), nil
}
