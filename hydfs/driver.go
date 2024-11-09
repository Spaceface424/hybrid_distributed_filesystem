package hydfs

import (
	"bufio"
	"cs425/mp3/hydfs/repl"
	"cs425/mp3/hydfs/swim"
	"cs425/mp3/shared"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"os"
	"strings"

	"github.com/huandu/skiplist"
)

func StartHydfs(introducer string, verbose bool) {
	hostname, _ := os.Hostname()
	cur_member, member_change_chan := swim.StartServer(hostname, introducer, verbose)
	this_member = &shared.MemberInfo{
		Address: cur_member.Address,
		ID:      cur_member.ID,
		Hash:    cur_member.Hash,
		State:   cur_member.State,
		IncNum:  cur_member.IncNum,
	}
	node_hash = this_member.Hash
	files = skiplist.New(skiplist.Uint32)
	members = skiplist.New(skiplist.Uint32)
	cache = make(map[uint32]*CachedFile)
	enable_cache = true
	member_change_chan <- struct{}{}
	hash_func = fnv.New32()
	hydfs_log = log.New(os.Stdout, fmt.Sprintf("hydfs.main Node %d, ", cur_member.ID), log.Ltime|log.Lshortfile)

	cleanupDir()

	go handleMembershipChange(member_change_chan)
	go StartGRPCServer(hostname + ":" + GRPC_PORT)

	commandLoop()
}

func hydfsCreate(local_filename string, hydfs_filename string) (bool, error) {
	if !enoughMembers() {
		return false, fmt.Errorf("Error: %w with current number of members: %d", ErrNotEnoughMembers, members.Len())
	}
	// read in file and hash filename
	contents := readFile(local_filename)
	file_hash := hashFilename(hydfs_filename)
	// get target node based on hashed filename
	target_hash, target := getMainFileTarget(file_hash)
	hydfs_log.Printf("%s hash to %d", hydfs_filename, file_hash)
	// construct rpc structures and make rpc call
	block := []*repl.FileBlock{{BlockNode: target_hash, BlockID: 0, Data: contents}}
	file_rpc := &repl.File{Filename: hydfs_filename, Blocks: block}
	hydfs_log.Printf("[INFO] Sending create request to node %d, hash %d", target.ID, target.Hash)
	return sendCreateRPC(target, file_rpc), nil
}

func hydfsGet(hydfs_filename string, local_filename string) (bool, error) {
	if !enoughMembers() {
		return false, fmt.Errorf("Error: %w with current number of members: %d", ErrNotEnoughMembers, members.Len())
	}
	file_hash := hashFilename(hydfs_filename)
	//cache check
	if enable_cache && cache[file_hash] != nil {
		cached_file := cache[file_hash]
		local_file, err := os.Create(local_filename)
		if err != nil {
			fmt.Println("[ERROR] os create error:", err)
		}
		defer local_file.Close()
		local_file.Write(cached_file.file)
		hydfs_log.Printf("[INFO] GET Wrote %s into %s, pulled from cache", hydfs_filename, local_filename)
		return true, nil
	}

	_, target := getReplicaFileTarget(file_hash)
	get_rpc := &repl.RequestGetData{Filename: hydfs_filename}
	hydfs_filedata := sendGetRPC(target, get_rpc)
	if hydfs_filedata == nil {
		return false, fmt.Errorf("Error: GetRPC Call had an error")
	}
	//add to cache
	if enable_cache {
		if len(cache) >= cache_cap {
			var min_key uint32 = 0
			var min_timestamp uint32 = 9999999
			for key, value := range cache {
				if min_timestamp > value.timestamp {
					min_key = key
				}
			}
			hydfs_log.Printf("[INFO] GET caused exceeded capacity in cache, deleted oldest entry")
			delete(cache, min_key)
		}
		hydfs_log.Printf("[INFO] GET Wrote %s into the cache", local_filename)
		new_cache := &CachedFile{timestamp: cache_ts, file: hydfs_filedata}
		cache_ts += 1
		cache[file_hash] = new_cache
	}
	local_file, err := os.Create(local_filename)
	if err != nil {
		fmt.Println("[ERROR] os create error:", err)
	}
	defer local_file.Close()
	local_file.Write(hydfs_filedata)
	hydfs_log.Printf("[INFO] GET Wrote %s into %s", hydfs_filename, local_filename)
	return true, nil
}

func hydfsVMGet(hydfs_filename string, local_filename string, vm_address string) (bool, error) {
	if !enoughMembers() {
		return false, fmt.Errorf("Error: %w with current number of members: %d", ErrNotEnoughMembers, members.Len())
	}
	file_hash := hashFilename(hydfs_filename)
	//cache check
	if enable_cache && cache[file_hash] != nil {
		cached_file := cache[file_hash]
		local_file, err := os.Create(local_filename)
		if err != nil {
			fmt.Println("[ERROR] os create error:", err)
		}
		defer local_file.Close()
		local_file.Write(cached_file.file)
		hydfs_log.Printf("[INFO] GET Wrote %s into %s, pulled from cache", hydfs_filename, local_filename)
		return true, nil
	}

	_, target := getVMFileTarget(file_hash, vm_address)
	if target == nil {
		return false, fmt.Errorf("[ERROR] VM_GET could not find node with address %s", vm_address)
	}
	get_rpc := &repl.RequestGetData{Filename: hydfs_filename}
	hydfs_filedata := sendGetRPC(target, get_rpc)
	if hydfs_filedata == nil {
		return false, fmt.Errorf("Error: GetRPC Call had an error")
	}
	//add to cache
	if enable_cache {
		if len(cache) >= cache_cap {
			var min_key uint32 = 0
			var min_timestamp uint32 = 9999999
			for key, value := range cache {
				if min_timestamp > value.timestamp {
					min_key = key
				}
			}
			hydfs_log.Printf("[INFO] GET caused exceeded capacity in cache, deleted oldest entry")
			delete(cache, min_key)
		}
		hydfs_log.Printf("[INFO] GET Wrote %s into the cache", local_filename)
		new_cache := &CachedFile{timestamp: cache_ts, file: hydfs_filedata}
		cache_ts += 1
		cache[file_hash] = new_cache
	}
	local_file, err := os.Create(local_filename)
	if err != nil {
		fmt.Println("[ERROR] os create error:", err)
	}
	defer local_file.Close()
	local_file.Write(hydfs_filedata)
	hydfs_log.Printf("[INFO] GET Wrote %s into %s", hydfs_filename, local_filename)
	return true, nil
}

func hydfsAppend(local_filename string, hydfs_filename string) (bool, error) {
	// TODO:
	if !enoughMembers() {
		return false, fmt.Errorf("Error: %w with current number of members: %d", ErrNotEnoughMembers, members.Len())
	}
	file, err := os.Open(local_filename)
	if err != nil {
		hydfs_log.Println("[ERROR] Open file error:", err)
		return false, nil
	}
	defer file.Close()
	contents, err := io.ReadAll(file)
	if err != nil {
		hydfs_log.Println("[ERROR] Open file error:", err)
		return false, nil
	}

	file_hash := hashFilename(hydfs_filename)
	//check cache for remove
	if enable_cache && cache[file_hash] != nil {
		hydfs_log.Printf("[INFO] APPEND deleted entry for %s from the cache", local_filename)
		delete(cache, file_hash)
	}
	_, target := getReplicaFileTarget(file_hash)
	data := &repl.FileBlock{BlockNode: file_hash, BlockID: 9999, Data: contents} // BlockID has temp val, will be set in rpc call
	append_rpc := &repl.AppendData{Filename: hydfs_filename, Block: data}
	hydfs_log.Printf("[INFO] Sending append request for file %s to node %d", hydfs_filename, target.ID)
	return sendAppendRPC(target, append_rpc), nil
}

func hydfsMerge(filepath string) (bool, error) {
	// TODO:
	if !enoughMembers() {
		return false, fmt.Errorf("Error: %w with current number of members: %d", ErrNotEnoughMembers, members.Len())
	}
	return false, nil
}

// Handle events from membership change channel
// Make re-replication calls and re-assign sorted_members ordered map
func handleMembershipChange(member_change_chan chan struct{}) {
	for {
		<-member_change_chan
		swim.Members.Mu.Lock()
		mu.Lock()

		hydfs_log.Println("[INFO] Handling membership change")
		members.Init()
		for _, member := range swim.Members.MemberMap {
			if member.State == shared.NodeState_ALIVE {
				members.Set(member.Hash, member)
			}
		}
		members.Set(node_hash, this_member)
		hydfs_log.Println("[INFO] Reset members")
		// check if became primary replica for new range
		// make re-replication calls to secondary replicas
		primary_replica_filehashes := getPrimaryReplicaFiles()
		// only make replication calls if primary files owned > 0 and enough members in group
		if len(primary_replica_filehashes) > 0 && enoughMembers() {
			hydfs_log.Println("[INFO] making replication requests")
			for _, repl_hash := range getReplicas() {
				target_replica := members.Get(repl_hash).Value.(*shared.MemberInfo)
				go sendReplicationRPC(target_replica, primary_replica_filehashes)
			}
		}

		mu.Unlock()
		swim.Members.Mu.Unlock()
	}
}

// user command loop
func commandLoop() {
	scanner := bufio.NewReader(os.Stdin)
	for {
		text, _ := scanner.ReadBytes('\n')
		command := string(text[:len(text)-1])

		switch commandParts := strings.Split(command, " "); commandParts[0] {
		case "list_mem":
			swim.PrintMembershipList()
		case "store":
			listFiles()
		case "create":
			if len(commandParts) != 3 {
				fmt.Println("usage: create localfilename HyDFSfilename")
				continue
			}
			local_filename := commandParts[1]
			hydfs_filename := commandParts[2]
			res, err := hydfsCreate(local_filename, hydfs_filename)
			if err != nil {
				fmt.Println("Create error:", err)
				continue
			}
			if res {
				fmt.Printf("Create %s in HyDFS from %s\n", hydfs_filename, local_filename)
			} else {
				fmt.Println("Create failed")
			}
		case "get":
			if len(commandParts) != 3 {
				fmt.Println("usage: get HyDFSfilename localfilename")
				continue
			}
			hydfs_filename := commandParts[1]
			local_filename := commandParts[2]
			_, err := hydfsGet(hydfs_filename, local_filename)
			if err != nil {
				fmt.Println("Get error:", err)
				continue
			}
			fmt.Printf("Get %s from HyDFS into %s\n", hydfs_filename, local_filename)
		case "append":
			if len(commandParts) != 3 {
				fmt.Println("usage: append localfilename HyDFSfilename")
				continue
			}
			local_filename := commandParts[1]
			hydfs_filename := commandParts[2]
			_, err := hydfsAppend(local_filename, hydfs_filename)
			if err != nil {
				fmt.Println("Append error:", err)
				continue
			}
			fmt.Printf("Append %s into %s in HyDFS\n", local_filename, hydfs_filename)
		case "merge":
			if len(commandParts) != 2 {
				fmt.Println("usage: merge HyDFSfilename")
				continue
			}
			// TODO: merge
		case "mem":
			printMemberDict()
		case "getfromreplica":
			
		default:
			fmt.Println("Unknown command...")
		}
	}
}
