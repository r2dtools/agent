package server

import (
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

		handle(conn)
	}
}

func handle(conn net.Conn) {
	logger.Info("Connection successfully handled")
}
