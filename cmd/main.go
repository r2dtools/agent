package main

import (
	"github.com/r2dtools/agent/cmd/server"
)

func main() {

	if err := server.RootCmd.Execute(); err != nil {
		panic(err)
	}
}
