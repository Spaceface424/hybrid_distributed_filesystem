package main

import (
	server "cs425/mp2/grep_server/server_lib"
	"os"
)

var usageMessage = `usage: go run server/main.go port`

func main() {
	// check correct num args
	args := os.Args
	hostname, _ := os.Hostname()
	if len(args) != 2 {
		server.StartGrepGRPCServer(hostname + ":8000")
	} else {
		server.StartGrepGRPCServer(hostname + ":" + args[1])
	}
}
