package main

import (
	swim "cs425/mp2/swim/swim_lib"
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

var usageMessage = `Usage: go run swim/main.go [OPTIONS] GRPC_PORT UDP_PORT

Options:
  -i, hostname of introducer node - dont set for introducer
  -v, enable verbose output to stdout

Examples:
  go run swim/main.go -v 8000 9000
  go run swim/main.go -v -i fa24-cs425-1401.cs.illinois.edu:8000 8000 9000
  `

func main() {
	var introducer string
	var verbose bool
	flag.StringVarP(&introducer, "Introducer Hostname", "i", "", "introducer hostname or nil to run as introducer")
	flag.BoolVarP(&verbose, "Verbose", "v", false, "set verbose logging to stdout")
	flag.Parse()
	args := flag.Args()
	if len(args) == 2 {
		hostname, _ := os.Hostname()
		if introducer == "" {
			fmt.Printf("Introducer Started at %s\n", hostname+":"+args[0])
		}
		go swim.StartGRPCServer(hostname + ":" + args[0])
		go swim.UDPServer(hostname+":"+args[1], introducer, verbose)
		wait := make(chan struct{})
		<-wait
	} else {
		fmt.Println(usageMessage)
		os.Exit(1)
	}
}
