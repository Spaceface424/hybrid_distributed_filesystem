package main

import (
	client "cs425/mp2/client/client_lib"
	"cs425/mp2/grep"
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

var usageMessage = `usage: go run client/main.go -o hosts_file [-clFivnhqre] regex`

func main() {
	// parse grep flags
	params := grep.ParseParams()

	// check correct num args
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println(usageMessage)
		os.Exit(2)
	}
	params.Expr = args[0]
	if params.Swim {
		client.SwimCMD(args, params)
		return
	}
	res := client.Grep(params)
	if params.Count {
		fmt.Println("Total lines matched:", res)
	}
}
