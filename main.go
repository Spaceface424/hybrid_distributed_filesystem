package main

import (
	"cs425/mp3/hydfs"
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

var usageMessage = `Usage: go run main.go [OPTIONS]

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

	var grpc_port string = "8000"
	var udp_port string = "9000"

	// if len(args) != 2 {
	// 	grpc_port = "8000"
	// 	udp_port = "9000"
	// } else {
	// 	grpc_port = args[]
	// 	udp_port = "9000"
	//  }

	hostname, _ := os.Hostname()
	if introducer == "" {
		fmt.Printf("Introducer Started at %s\n", hostname+":"+args[0])
	}

	hydfs.StartHydfs()

	wait := make(chan struct{})
	<-wait
}
