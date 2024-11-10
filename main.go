package main

import (
	"cs425/mp3/hydfs"

	flag "github.com/spf13/pflag"
)

var usageMessage = `Usage: go run main.go [OPTIONS]

Options:
  -i, hostname of introducer node - dont set for introducer
  -v, enable verbose output to stdout
  `

func main() {
	var introducer string
	var verbose bool
	var cache_size int

	flag.StringVarP(&introducer, "Introducer Hostname", "i", "", "introducer hostname or nil to run as introducer")
	flag.BoolVarP(&verbose, "Verbose", "v", false, "set verbose logging to stdout")
	flag.IntVarP(&cache_size, "Cache Limit", "c", 0, "set cache limit in KB, set 0 to disable caching")

	flag.Parse()

	hydfs.StartHydfs(introducer, verbose, cache_size*1000)
}
