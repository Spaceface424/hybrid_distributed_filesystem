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

// Send re-replication request to target replica node
func sendReplicationRPC(target *shared.MemberInfo, primary_replica_filehashes []uint32) bool {
	rpc_files := make([]*repl.File, 0)
	for _, file_hash := range primary_replica_filehashes {
		cur_file := files.Get(file_hash).Value.(File)
		rpc_files = append(rpc_files, &repl.File{Filename: cur_file.filename, Blocks: getBlocks(cur_file.filename, false)})
	}
	rpc_request_files := &repl.RequestFiles{Files: rpc_files}

	target_addr := strings.Split(target.Address, ":")[0] + ":" + GRPC_PORT
	conn, err := grpc.NewClient(target_addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		hydfs_log.Printf("[WARNING] gRPC did not connect: %v", err)
		return false
	}
	defer conn.Close()

	client := repl.NewReplicationClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	response_missing, err := client.RequestAsk(ctx, rpc_request_files)
	if err != nil {
		hydfs_log.Printf("[WARNING] gRPC call error: %v", err)
		return false
	}
	num_blocks_missing := getNumBlocks(response_missing)
	hydfs_log.Printf("[INFO] Received request ask response with %d blocks missing", num_blocks_missing)
	// skip sending data if no data to be replicated
	if num_blocks_missing == 0 {
		return true
	}

	rpc_request_data := fillData(response_missing)
	response_ack, err := client.RequestSend(ctx, rpc_request_data)

	return response_ack.OK
}
