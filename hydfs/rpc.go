package hydfs

import (
	"context"
	"cs425/mp3/hydfs/repl"
	"cs425/mp3/shared"
	"net"

	"google.golang.org/grpc"
)

type HydfsRPCserver struct {
	repl.UnimplementedReplicationServer
}

func StartGRPCServer(host string) {
	lis, err := net.Listen("tcp", host)
	if err != nil {
		hydfs_log.Fatalf("[ERROR] failed to listen: %v", err)
	}
	s := grpc.NewServer()
	repl.RegisterReplicationServer(s, &HydfsRPCserver{})
	hydfs_log.Printf("[INFO] HyDFS gRPC server started at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		hydfs_log.Fatalf("[ERROR] Failed to serve: %v", err)
	}
}

// Receive a list of files that need to be replicated
// Return list of files missing from local copies
func (s *HydfsRPCserver) RequestAsk(ctx context.Context, request *repl.RequestFiles) (*repl.RequestMissing, error) {
	res := &repl.RequestMissing{MissingFiles: make([]*repl.File, 0)}

	mu.Lock()
	defer mu.Unlock()

	// iterate over files from request and check if exists on current node
	// reply with missing file blocks
	for _, file := range request.Files {
		cur_file := &repl.File{Filename: file.Filename, Blocks: make([]*repl.FileBlock, 0)}
		// if file doesn't exist then add all blocks to missing list
		if files.Get(hashFilename(file.Filename)) == nil {
			cur_file.Blocks = file.Blocks
		} else {
			for _, block := range file.Blocks {
				if !blockExists(file.Filename, block.BlockNode, block.BlockID) {
					cur_file.Blocks = append(cur_file.Blocks, block)
				}
			}
		}
		res.MissingFiles = append(res.MissingFiles, cur_file)
	}

	return res, nil
}

// Receive a list of files with data to be replicated
// Create the new blocks on current node
// Return OK
func (s *HydfsRPCserver) RequestSend(ctx context.Context, request *repl.RequestData) (*repl.RequestAck, error) {
	for _, file := range request.DataFiles {
		for _, block := range file.Blocks {
			createBlock(file.Filename, block.BlockNode, block.BlockID, block.Data)
		}
	}

	return &repl.RequestAck{OK: true}, nil
}

// Process a file create request
// Create the file directory and the initial block
func (s *HydfsRPCserver) RequestCreate(ctx context.Context, request *repl.CreateData) (*repl.RequestAck, error) {
	mu.Lock()
	defer mu.Unlock()

	hydfs_log.Printf("[INFO] RPC Serving create request for file: %s", request.NewFile.Filename)

	file_rpc := request.NewFile
	file_struct := File{filename: file_rpc.Filename, nextID: 1}
	file_hash := hashFilename(file_rpc.Filename)
	files.Set(file_hash, file_struct)
	block := request.NewFile.Blocks[0]

	// fail request if file already exists
	if files.Get(file_hash) != nil {
		hydfs_log.Printf("[WARNING] RPC Serving create request file %s already exists", request.NewFile.Filename)
		return &repl.RequestAck{OK: false}, nil
	}

	createBlock(file_rpc.Filename, block.BlockNode, block.BlockID, block.Data)

	for _, replica_hash := range getReplicas() {
		// send create request to neighbor nodes
		replica := members.Get(replica_hash)
		if replica == nil {
			hydfs_log.Println("[WARNING] replica == nil")
			continue
		}
		go sendCreateReplicaRPC(replica.Value.(*shared.MemberInfo), file_rpc)
	}

	return &repl.RequestAck{OK: true}, nil
}

// Process a file create request for a replica
// Create the file directory and the initial block
func (s *HydfsRPCserver) RequestReplicaCreate(ctx context.Context, request *repl.CreateData) (*repl.RequestAck, error) {
	mu.Lock()
	defer mu.Unlock()

	hydfs_log.Printf("[INFO] RPC Serving replica create request for file: %s", request.NewFile.Filename)

	file_rpc := request.NewFile
	file_struct := File{filename: file_rpc.Filename, nextID: 1}
	file_hash := hashFilename(file_rpc.Filename)
	files.Set(file_hash, file_struct)
	block := request.NewFile.Blocks[0]

	createBlock(file_rpc.Filename, block.BlockNode, block.BlockID, block.Data)

	return &repl.RequestAck{OK: true}, nil
}
