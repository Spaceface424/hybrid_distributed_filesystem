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
	hydfs_log.Println("[INFO] RPC Serving request ask")
	res := &repl.RequestMissing{MissingFiles: make([]*repl.File, 0)}

	mu.Lock()
	defer mu.Unlock()

	// iterate over files from request and check if exists on current node
	// reply with missing file blocks
	num_blocks_missing := 0
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
		if len(cur_file.Blocks) != 0 {
			res.MissingFiles = append(res.MissingFiles, cur_file)
			num_blocks_missing += len(cur_file.Blocks)
		}
	}
	hydfs_log.Printf("[INFO] RPC requesting %v missing blocks", num_blocks_missing)
	return res, nil
}

// Receive a list of files with data to be replicated
// Create the new blocks on current node
// Return OK
func (s *HydfsRPCserver) RequestSend(ctx context.Context, request *repl.RequestData) (*repl.RequestAck, error) {
	hydfs_log.Println("[INFO] RPC serving request send")
	for _, file := range request.DataFiles {
		filehash := hashFilename(file.Filename)
		if files.Get(filehash) == nil {
			new_file := &File{filename: file.Filename, nextID: 0}
			files.Set(filehash, new_file)
		}
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
	file_struct := &File{filename: file_rpc.Filename, nextID: 1}
	file_hash := hashFilename(file_rpc.Filename)

	// fail request if file already exists
	if files.Get(file_hash) != nil {
		hydfs_log.Printf("[WARNING] RPC Serving create request file %s already exists", request.NewFile.Filename)
		return &repl.RequestAck{OK: false}, nil
	}

	files.Set(file_hash, file_struct)
	block := request.NewFile.Blocks[0]
	block.BlockNode = this_member.Hash
	createBlock(file_rpc.Filename, block.BlockNode, 0, block.Data)
	replicas := getReplicas()
	hydfs_log.Printf("[INFO] RPC Replicating create request to replicas: %v", replicas)
	for _, replica_hash := range replicas {
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
	file_struct := &File{filename: file_rpc.Filename, nextID: 0}
	file_hash := hashFilename(file_rpc.Filename)
	files.Set(file_hash, file_struct)
	block := request.NewFile.Blocks[0]

	createBlock(file_rpc.Filename, block.BlockNode, block.BlockID, block.Data)

	return &repl.RequestAck{OK: true}, nil
}

func (s *HydfsRPCserver) RequestGet(ctx context.Context, request *repl.RequestGetData) (*repl.ResponseGetData, error) {
	mu.Lock()
	defer mu.Unlock()

	hydfs_log.Printf("[INFO] RPC Serving get request for file: %s", request.Filename)
	filename := request.Filename
	// check if file exists
	if files.Get(hashFilename(filename)) == nil {
		return nil, nil
	}
	// get blocks and sort and read into memory
	file_blocks := getBlocks(filename, true)
	sortBlocks(file_blocks)
	return &repl.ResponseGetData{FileData: readAllBlocks(filename, file_blocks)}, nil
}

func (s *HydfsRPCserver) RequestAppend(ctx context.Context, request *repl.AppendData) (*repl.RequestAck, error) {
	mu.Lock()
	defer mu.Unlock()

	file_hash := hashFilename(request.Filename)
	block := request.Block
	// fails if file does not exist in files
	if files.Get(file_hash) == nil {
		hydfs_log.Printf("[WARNING] RPC Serving append request file %s does not exist", request.Filename)
		return &repl.RequestAck{OK: false}, nil
	}
	local_file := files.Get(file_hash).Value.(*File)
	block.BlockNode = this_member.Hash
	block.BlockID = local_file.nextID
	local_file.nextID += 1
	hydfs_log.Printf("[INFO] RPC Serving append request for file %s", request.Filename)
	createBlock(request.Filename, block.BlockNode, block.BlockID, block.Data)
	hydfs_log.Printf("[INFO] RPC append created new block for file %s", request.Filename)
	for _, replica_hash := range getPrimaryReplicas(file_hash) {
		// send append request to neighbor nodes
		hydfs_log.Printf("[INFO] RPC sending append replica created new block for file %s", request.Filename)
		replica := members.Get(replica_hash)
		if replica == nil {
			hydfs_log.Println("[WARNING] replica == nil")
			continue
		}
		go sendAppendReplicaRPC(replica.Value.(*shared.MemberInfo), request)
	}

	return &repl.RequestAck{OK: true}, nil
}

func (s *HydfsRPCserver) RequestReplicaAppend(ctx context.Context, request *repl.AppendData) (*repl.RequestAck, error) {
	mu.Lock()
	defer mu.Unlock()

	file_hash := hashFilename(request.Filename)
	block := request.Block
	// fails if file does not exist in files
	if files.Get(file_hash) == nil {
		hydfs_log.Printf("[WARNING] RPC Serving replica append request file %s does not exist", request.Filename)
		return &repl.RequestAck{OK: false}, nil
	}
	hydfs_log.Printf("[INFO] RPC Serving replica append request for file: %s", request.Filename)
	createBlock(request.Filename, block.BlockNode, block.BlockID, block.Data)
	hydfs_log.Printf("[INFO] RPC append created new block for file %s", request.Filename)
	return &repl.RequestAck{OK: true}, nil
}

// responds whether or not the current node is storing this file
func (s *HydfsRPCserver) RequestList(ctx context.Context, request *repl.File) (*repl.RequestAck, error) {
	mu.Lock()
	defer mu.Unlock()

	filehash := hashFilename(request.Filename)
	storing_file := files.Get(filehash) != nil
	hydfs_log.Printf("[INFO] RPC Serving ls request for file %s responding %v", request.Filename, storing_file)
	return &repl.RequestAck{OK: storing_file}, nil
}

func (s *HydfsRPCserver) RequestMultiAppend(ctx context.Context, request *repl.MultiAppendData) (*repl.RequestAck, error) {
	hydfs_log.Printf("[INFO] RPC MultiAppend received request for hydfs file %s to local file %s", request.HydfsFilename, request.LocalFilename)
	ok, _ := hydfsAppend(request.LocalFilename, request.HydfsFilename)
	return &repl.RequestAck{OK: ok}, nil
}
