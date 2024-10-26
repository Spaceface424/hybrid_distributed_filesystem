package server

import (
	"bytes"
	"context"
	"cs425/mp2/grep"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

type GrepRPCserver struct {
	grep.UnimplementedGrepServer
}

func StartGrepGRPCServer(host string) {
	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	grep.RegisterGrepServer(s, &GrepRPCserver{})
	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *GrepRPCserver) CallGrep(ctx context.Context, params *grep.GRPCParams) (*grep.GrepOutput, error) {
	// rpc call to grep local log file
	log.Println("Serving RPC.Grep on", params.Logfile)
	buf := new(bytes.Buffer)
	params_struct := grep.ProtoToStruct(params)
	cmd := grep.Command(os.Stdin, buf, os.Stderr, params_struct, []string{params_struct.Logfile})
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return &grep.GrepOutput{Result: buf.String()}, nil
}
