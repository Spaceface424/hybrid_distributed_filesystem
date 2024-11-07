package hydfs

import (
	"context"
	"cs425/mp3/hydfs/repl"
	"cs425/mp3/shared"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Send a create rpc to target node
func sendCreateRPC(target *shared.MemberInfo, file_rpc *repl.File) bool {
	target_addr := strings.Split(target.Address, ":")[0] + ":" + GRPC_PORT
	hydfs_log.Printf("[INFO] RPC Sending create request to %s for file: %s", target_addr, file_rpc.Filename)
	conn, err := grpc.NewClient(target_addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		hydfs_log.Printf("[WARNING] gRPC did not connect: %v", err)
		return false
	}
	defer conn.Close()

	client := repl.NewReplicationClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	request_data := &repl.CreateData{NewFile: file_rpc}
	response, err := client.RequestCreate(ctx, request_data)
	if err != nil {
		hydfs_log.Printf("[WARNING] gRPC call error: %v", err)
		return false
	}
	return response.OK
}

// Send a create rpc to target replica node
func sendCreateReplicaRPC(target *shared.MemberInfo, file_rpc *repl.File) bool {
	target_addr := strings.Split(target.Address, ":")[0] + ":" + GRPC_PORT
	conn, err := grpc.NewClient(target_addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		hydfs_log.Printf("[WARNING] gRPC did not connect: %v", err)
		return false
	}
	defer conn.Close()

	client := repl.NewReplicationClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	request_data := &repl.CreateData{NewFile: file_rpc}
	response, err := client.RequestReplicaCreate(ctx, request_data)
	if err != nil {
		hydfs_log.Printf("[WARNING] gRPC call error: %v", err)
		return false
	}
	return response.OK
}
