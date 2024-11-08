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
	file, err := os.Open(local_filename)
	if err != nil {
		hydfs_log.Println("[ERROR] Open file error:", err)
		return false, nil
	}
	defer file.Close()

	// read contents of file into memory
	contents, err := io.ReadAll(file)
	if err != nil {
		hydfs_log.Println("[ERROR] Open file error:", err)
		return false, nil
	}

	file_hash := hashFilename(hydfs_filename)
	target_hash, target := getFileTarget(file_hash)
	hydfs_log.Printf("%s hash to %d", hydfs_filename, file_hash)
	block := []*repl.FileBlock{{BlockNode: target_hash, BlockID: 0, Data: contents}}
	file_rpc := &repl.File{Filename: hydfs_filename, Blocks: block}
	hydfs_log.Printf("[INFO] Sending create request to node %d", target.ID)
	return sendCreateRPC(target, file_rpc), nil
}

func hydfsGet(hydfs_filename string, local_filename string) (bool, error) {
	// TODO:
    if !enoughMembers() {
        return false, fmt.Errorf("Error: %w with current number of members: %d", ErrNotEnoughMembers, members.Len())
    }
    file_hash := hashFilename(hydfs_filename)
    _, target := getFileTarget(file_hash)
    get_rpc := &repl.GetData{Filename: hydfs_filename}
    hydfs_log.Printf("[INFO] Sending get request for file %s to node %d", hydfs_filename, target.ID)
    file := sendGetRPC(target, get_rpc)
    if file == nil {
        return false, fmt.Errorf("Error: GetRPC Call had an error")
    }
	for _, block := range file.Blocks {
		createBlock(local_filename, block.BlockNode, block.BlockID, block.Data)
	}
    // block := file.Blocks[0]
    // createBlock(local_filename, block.BlockNode, block.BlockID, block.Data) //IDK if this works
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
    _, target := getFileTarget(file_hash)
    data := &repl.FileBlock{BlockNode: file_hash, BlockID: 9999, Data: contents} //BlockID has temp val, will be set in rpc call
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
			members.Set(member.Hash, member)
		}
		members.Set(node_hash, this_member)

		// TODO: handle re-replication
		// check if became primary replica for new range
		// make re-replication calls to secondary replicas
		// for _, repl_hash := range getReplicas() {
		//
		// }

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
		default:
			fmt.Println("Unknown command...")
		}
	}
}
