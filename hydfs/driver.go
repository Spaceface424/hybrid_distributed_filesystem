package hydfs

import (
	"bufio"
	"cs425/mp3/hydfs/repl"
	"cs425/mp3/hydfs/swim"
	"cs425/mp3/shared"
	"fmt"
	"hash/fnv"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/huandu/skiplist"
)

func StartHydfs(introducer string, verbose bool, cache_size int) {
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
	member_change_chan <- struct{}{}
	hash_func = fnv.New32()
	hydfs_log = log.New(os.Stdout, fmt.Sprintf("hydfs.main Node %d, ", cur_member.ID), log.Ltime|log.Lshortfile)
	enable_cache = cache_size > 0
	if enable_cache {
		cache = createCache(cache_size)
	}

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
	_, target := getMainFileTarget(file_hash)
	hydfs_log.Printf("%s hash to %d", hydfs_filename, file_hash)
	// construct rpc structures and make rpc call
	block := []*repl.FileBlock{{Data: contents}}
	file_rpc := &repl.File{Filename: hydfs_filename, Blocks: block}
	hydfs_log.Printf("[INFO] Sending create request to node %d, hash %d", target.ID, target.Hash)
	return sendCreateRPC(target, file_rpc), nil
}

func hydfsGet(hydfs_filename string, local_filename string) (bool, error) {
	if !enoughMembers() {
		return false, fmt.Errorf("Error: %w with current number of members: %d", ErrNotEnoughMembers, members.Len())
	}
	file_hash := hashFilename(hydfs_filename)
	// cache check
	if enable_cache {
		data := cache.getFile(hydfs_filename)
		if data != nil {
			hydfs_log.Printf("[INFO] GET cache hit on file %s!", hydfs_filename)
			err := os.MkdirAll(filepath.Dir(local_filename), 0777)
			if err != nil {
				fmt.Println("[ERROR] os create error:", err)
			}
			local_file, err := os.Create(local_filename)
			if err != nil {
				fmt.Println("[ERROR] os create error:", err)
			}
			defer local_file.Close()
			local_file.Write(data)
			hydfs_log.Printf("[INFO] GET Wrote %s into %s from cache", hydfs_filename, local_filename)
			return true, nil
		}
		hydfs_log.Printf("[INFO] GET cache miss on file %s!", hydfs_filename)
	}

	// make get rpc call to a replica
	_, target := getReplicaFileTarget(file_hash)
	get_rpc := &repl.RequestGetData{Filename: hydfs_filename}
	hydfs_filedata := sendGetRPC(target, get_rpc)
	if hydfs_filedata == nil {
		return false, fmt.Errorf("Error: GetRPC Call had an error")
	}

	// add to cache
	if enable_cache {
		cache.addFile(hydfs_filename, hydfs_filedata)
	}

	// write get data to local_filename
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

	// make get rpc call to vm_address
	_, target := getVMFileTarget(file_hash, vm_address)
	if target == nil {
		return false, fmt.Errorf("[ERROR] Replica_GET could not find node with address %s", vm_address)
	}
	get_rpc := &repl.RequestGetData{Filename: hydfs_filename}
	hydfs_filedata := sendGetRPC(target, get_rpc)
	if hydfs_filedata == nil {
		return false, fmt.Errorf("Error: GetRPC Call had an error")
	}

	// write get data to local_filename
	local_file, err := os.Create(local_filename)
	if err != nil {
		fmt.Println("[ERROR] os create error:", err)
	}
	defer local_file.Close()
	local_file.Write(hydfs_filedata)
	hydfs_log.Printf("[INFO] Replica_GET Wrote %s into %s", hydfs_filename, local_filename)
	return true, nil
}

func hydfsAppend(local_filename string, hydfs_filename string) (bool, error) {
	if !enoughMembers() {
		return false, fmt.Errorf("Error: %w with current number of members: %d", ErrNotEnoughMembers, members.Len())
	}
	// read local_file into memory and send append rpc
	contents := readFile(local_filename)
	file_hash := hashFilename(hydfs_filename)
	// check cache for remove
	if enable_cache {
		cache.invalidateFile(hydfs_filename)
	}
	_, target := getReplicaFileTarget(file_hash)
	data := &repl.FileBlock{Data: contents}
	append_rpc := &repl.AppendData{Filename: hydfs_filename, Block: data}
	hydfs_log.Printf("[INFO] Sending append request for file %s to node %d, hash %d", hydfs_filename, target.ID, target.Hash)
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
		// check if became primary replica for new range
		// make re-replication calls to secondary replicas
		primary_replica_filehashes := getPrimaryReplicaFiles()
		// only make replication calls if primary files owned > 0 and enough members in group
		if len(primary_replica_filehashes) > 0 && enoughMembers() {
			hydfs_log.Println("[INFO] making replication requests")
			var wg sync.WaitGroup
			for _, repl_hash := range getReplicas() {
				wg.Add(1)
				target_replica := members.Get(repl_hash).Value.(*shared.MemberInfo)
				go func() {
					sendReplicationRPC(target_replica, primary_replica_filehashes)
					wg.Done()
				}()
			}
			mu.Unlock()
			wg.Wait()
			mu.Lock()
		}

		// check if files left replication range
		garbageCollectFiles()

		mu.Unlock()
		swim.Members.Mu.Unlock()

		swim.PrintMembershipList()
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
				fmt.Printf("SUCCESS Create %s in HyDFS from %s\n", hydfs_filename, local_filename)
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
			fmt.Printf("SUCCESS Get %s from HyDFS into %s\n", hydfs_filename, local_filename)
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
			fmt.Printf("SUCCESS Append %s into %s in HyDFS\n", local_filename, hydfs_filename)
		case "merge":
			if len(commandParts) != 2 {
				fmt.Println("usage: merge HyDFSfilename")
				continue
			}
			// TODO: merge
		case "mem":
			printMemberDict()
		case "getfromreplica":
			if len(commandParts) != 4 {
				fmt.Println("usage: get VMAddress HyDFSfilename localfilename")
				continue
			}
			vm_address := VM[commandParts[1]]
			hydfs_filename := commandParts[2]
			local_filename := commandParts[3]
			_, err := hydfsVMGet(hydfs_filename, local_filename, vm_address)
			if err != nil {
				fmt.Println("Get error:", err)
				continue
			}
			fmt.Printf("Get %s from VM Address %s into %s\n", hydfs_filename, vm_address, local_filename)
		case "ls":
			if len(commandParts) != 2 {
				fmt.Println("usage: ls HyDFSfilename")
				continue
			}
			ls(commandParts[1])
		case "multiappend":
			if len(commandParts) < 5 {
				fmt.Println("usage: multiappend hydfs_filename vm_1 vm_2 ... files local_filename_1 local_filename_2 ...")
				continue
			}
			vms := make([]string, 0)
			local_files := make([]string, 0)
			file_idx := 0
			for i, arg := range commandParts[2:] {
				if arg == "files" {
					file_idx = i + 3
					break
				}
				vms = append(vms, VM[arg]+":9000")
			}
			for _, arg := range commandParts[file_idx:] {
				local_files = append(local_files, arg)
			}
			if len(vms) != len(local_files) {
				fmt.Println("usage: multiappend hydfs_filename vm_1 vm_2 ... files local_filename_1 local_filename_2 ...")
				fmt.Println(vms)
				fmt.Println(local_files)
				continue
			}
			multiappend(commandParts[1], vms, local_files)
		case "load_dataset":
			num_files, _ := strconv.Atoi(commandParts[1])
			loadDataset(num_files, 4)
		case "exp3":
			num_files, _ := strconv.Atoi(commandParts[2])
			percent_append, _ := strconv.Atoi(commandParts[3])
			if commandParts[1] == "zipf" {
				getExperiment(num_files, num_files*2, zipfianSample, percent_append)
			} else {
				getExperiment(num_files, num_files*2, uniformSample, percent_append)
			}
		case "clear_cache":
			if enable_cache {
				cache.clear()
			}
			fmt.Println("Cache cleared")
		case "set_cache_size":
			new_limit, _ := strconv.Atoi(commandParts[1])
			if enable_cache {
				cache.setLimit(new_limit)
			}
			fmt.Printf("Cache cleared and size set to %dKB\n", new_limit)
		default:
			fmt.Println("Unknown command...")
		}
	}
}
