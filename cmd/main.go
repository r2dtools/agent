package main

import (
	"github.com/r2dtools/sslbot/cmd/server"
)

func main() {
	if err := server.RootCmd.Execute(); err != nil {
		panic(err)
	}
}
