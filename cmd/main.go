package main

import (
	"github.com/r2dtools/sslbot/cmd/server"
	"github.com/r2dtools/sslbot/config"
)

var Version string

func main() {
	config.Version = Version

	if err := server.CreateCli().Execute(); err != nil {
		panic(err)
	}
}
