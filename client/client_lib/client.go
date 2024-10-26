package client

import (
	"bufio"
	"context"
	"cs425/mp2/grep"
	"cs425/mp2/shared"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type serverInfo struct {
	m           map[string]string
	sorted_keys []string
}

func getServers(filename string) serverInfo {
	// open hosts config file and create map of hostnames, logfile
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	m := make(map[string]string)
	hosts := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(strings.TrimSpace(scanner.Text()), ",")
		if len(line) != 2 {
			fmt.Printf("Skipping invalid hostfile line: %s\n", line)
			continue
		}
		cur_host := strings.TrimSpace(line[0])
		hosts = append(hosts, cur_host)
		m[cur_host] = strings.TrimSpace(line[1])
	}
	return serverInfo{m: m, sorted_keys: hosts}
}

func writeOutput(id int, output []byte) {
	// write grep output to file
	filename := fmt.Sprintf("out/%v.txt", id)
	_ = os.MkdirAll("out", 0755)
	err := os.WriteFile(filename, output, 0644)
	if err != nil {
		log.Fatal("Failed to write to file", err)
	}
	log.Printf("output written to out/%v.txt", id)
}

func Grep(params grep.Params) int {
	// create channel to signal when goroutines are done
	doneChan := make(chan int)
	servers := getServers(params.Hosts)
	// create sorted slice of server names to make sure they are given correct id for output
	i := 1
	for _, host := range servers.sorted_keys {
		params.Logfile = servers.m[host]
		// call goroutine to make RPC call to server
		go callRPC(host, params, doneChan, i)
		i += 1
	}

	// wait for all go routine calls to finish and collect number of line matches
	grep_count := 0
	for i := 0; i < len(servers.m); i++ {
		grep_count += <-doneChan
	}
	return grep_count
}

func SwimCMD(args []string, params grep.Params) {
	cmd := args[0]
	SWIMIn := &shared.SWIMIn{
		Cmd: cmd,
	}
	if cmd == "drop_rate" {
		drop_rate, err := strconv.ParseFloat(args[2], 64)
		if err == nil {
			SWIMIn.DropRate = drop_rate
		}
	}
	doneChan := make(chan string)

	if cmd == "enable_sus" || cmd == "disable_sus" || cmd == "drop_rate" || cmd == "list_sus" || cmd == "list_mem" {
		servers := getServers(params.Hosts)
		// create sorted slice of server names to make sure they are given correct id for output
		for _, host := range servers.sorted_keys {
			// call goroutine to make RPC call to server
			go callSWIMRPC(host, SWIMIn, doneChan)
		}
		log.Printf("%s command sent to all machines", cmd)
		for i := 0; i < len(servers.m); i++ {
			fmt.Printf("%s", <-doneChan)
		}
	} else {
		server := args[1]
		go callSWIMRPC(server, SWIMIn, doneChan)
		fmt.Printf(<-doneChan)
	}
}

func callSWIMRPC(host string, command *shared.SWIMIn, doneChan chan<- string) {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect: %v", err)
		doneChan <- ""
		return
	}
	defer conn.Close()
	client := shared.NewIntroducerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	response, err := client.SWIMcmd(ctx, command)
	if err != nil {
		// log.Printf("gRPC call error: %v", err)
		doneChan <- ""
		return
	}
	doneChan <- response.Output
}

func callRPC(host string, args grep.Params, doneChan chan<- int, id int) {
	log.Printf("Sending gRPC call to node %s", host)
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := grep.NewGrepClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	proto_args := grep.StructToProto(args)
	response, err := client.CallGrep(ctx, proto_args)
	if err != nil {
		log.Fatalf("gRPC call error: %v", err)
	}

	cur_count := 0
	writeOutput(id, []byte(response.Result))
	fmt.Println(response.Result)
	if args.Count {
		cur_count, _ = strconv.Atoi(strings.TrimSpace(strings.Split(response.Result, ":")[1]))
		fmt.Print(response.Result)
	}
	doneChan <- cur_count
}
