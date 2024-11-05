package hydfs

import (
	"cs425/mp3/hydfs/repl"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type HydfsRPCserver struct {
	repl.UnimplementedReplicationServer
}

func StartGRPCServer(host string) {
	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	repl.RegisterReplicationServer(s, &HydfsRPCserver{})
	fmt.Printf("gRPC server started at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
