package hydfs

import (
	"context"
	"cs425/mp3/hydfs/repl"
	"cs425/mp3/shared"
	"fmt"
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
	// measure re-replication time
	start := time.Now()

	rpc_files := make([]*repl.File, 0)
	for _, file_hash := range primary_replica_filehashes {
		cur_file := files.Get(file_hash).Value.(*File)
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
	if err != nil {
		hydfs_log.Println("[ERROR] Calling RPC Request Send error: ", err)
		return false
	}

	elapsed := time.Since(start)
	hydfs_log.Printf("[INFO] Re-replication took %v", elapsed)
	return response_ack.OK
}

func sendGetRPC(target *shared.MemberInfo, file_rpc *repl.RequestGetData) []byte {
	// get file locally if stored at node
	if files.Get(hashFilename(file_rpc.Filename)) != nil {
		hydfs_log.Printf("[INFO] GET file %s found locally", file_rpc.Filename)
		file_blocks := getBlocks(file_rpc.Filename, true)
		sortBlocks(file_blocks)
		return readAllBlocks(file_rpc.Filename, file_blocks)
	}

	// make rpc get request to replica node
	target_addr := strings.Split(target.Address, ":")[0] + ":" + GRPC_PORT
	hydfs_log.Printf("[INFO] Sending get request to node %d, hash %d for file: %s", target.ID, target.Hash, file_rpc.Filename)
	conn, err := grpc.NewClient(target_addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		hydfs_log.Printf("[WARNING] gRPC did not connect: %v", err)
		return nil
	}
	defer conn.Close()

	client := repl.NewReplicationClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	response_data, err := client.RequestGet(ctx, file_rpc)
	if err != nil {
		hydfs_log.Printf("[WARNING] gRPC call error: %v", err)
		return nil
	}
	return response_data.FileData
}

func sendAppendRPC(target *shared.MemberInfo, file_rpc *repl.AppendData) bool {
	target_addr := strings.Split(target.Address, ":")[0] + ":" + GRPC_PORT
	hydfs_log.Printf("[INFO] RPC Sending append request to node %d, hash %d for file: %s", target.ID, target.Hash, file_rpc.Filename)
	conn, err := grpc.NewClient(target_addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		hydfs_log.Printf("[WARNING] gRPC did not connect: %v", err)
		return false
	}
	defer conn.Close()

	client := repl.NewReplicationClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	response, err := client.RequestAppend(ctx, file_rpc)
	if err != nil {
		hydfs_log.Printf("[WARNING] gRPC call error: %v", err)
		return false
	}
	return response.OK
}

func sendAppendReplicaRPC(target *shared.MemberInfo, file_rpc *repl.AppendData) bool {
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

	response, err := client.RequestReplicaAppend(ctx, file_rpc)
	if err != nil {
		hydfs_log.Printf("[WARNING] gRPC call error: %v", err)
		return false
	}
	return response.OK
}

func sendLsRPC(target *shared.MemberInfo, hydfs_filename string, ch chan *shared.MemberInfo) {
	target_addr := strings.Split(target.Address, ":")[0] + ":" + GRPC_PORT
	hydfs_log.Printf("[INFO] RPC Sending ls request to %s for file: %s", target_addr, hydfs_filename)
	conn, err := grpc.NewClient(target_addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		hydfs_log.Printf("[WARNING] gRPC did not connect: %v", err)
		ch <- nil
		return
	}
	defer conn.Close()

	client := repl.NewReplicationClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	request_data := &repl.File{Filename: hydfs_filename}
	response, err := client.RequestList(ctx, request_data)
	if err != nil {
		hydfs_log.Printf("[WARNING] gRPC call error: %v", err)
		ch <- nil
		return
	}

	if response.OK {
		ch <- target
	} else {
		ch <- nil
	}
}

func sendMultiAppendReqeustRPC(target *shared.MemberInfo, hydfs_filename string, local_filename string) {
	target_addr := strings.Split(target.Address, ":")[0] + ":" + GRPC_PORT
	hydfs_log.Printf("[INFO] RPC Sending multi-append request to node %d, hash %d for hydfs_file: %s to local_file: %s", target.ID, target.Hash, hydfs_filename, local_filename)
	conn, err := grpc.NewClient(target_addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		hydfs_log.Printf("[WARNING] gRPC did not connect: %v", err)
		return
	}
	defer conn.Close()

	client := repl.NewReplicationClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	request_data := &repl.MultiAppendData{HydfsFilename: hydfs_filename, LocalFilename: local_filename}
	response, err := client.RequestMultiAppend(ctx, request_data)
	if err != nil {
		hydfs_log.Printf("[WARNING] gRPC call error: %v", err)
		return
	}
	if response.OK {
		fmt.Printf("SUCCESS MultiAppend to node %d hash %d for hydfs_file: %s to local_file: %s\n", target.ID, target.Hash, hydfs_filename, local_filename)
	} else {
		fmt.Printf("FAILED MultiAppend to node %d hash %d for hydfs_file: %s to local_file: %s\n", target.ID, target.Hash, hydfs_filename, local_filename)
	}
}
