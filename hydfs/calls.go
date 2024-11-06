package hydfs

import (
	"context"
	"cs425/mp3/hydfs/repl"
	"cs425/mp3/shared"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Send a create rpc to target node
func sendCreateRPC(target *shared.MemberInfo, file_rpc *repl.File) bool {
	conn, err := grpc.NewClient(target.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("[WARNING] gRPC did not connect: %v", err)
		return false
	}
	defer conn.Close()

	client := repl.NewReplicationClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	request_data := &repl.CreateData{NewFile: file_rpc}
	response, err := client.RequestCreate(ctx, request_data)
	if err != nil {
		log.Printf("[WARNING] gRPC call error: %v", err)
		return false
	}
	return response.OK
}
