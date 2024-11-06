package hydfs

import (
	"bufio"
	"cs425/mp3/hydfs/repl"
	"cs425/mp3/hydfs/swim"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/huandu/skiplist"
)

func StartHydfs(introducer string, verbose bool) {
	cleanup()
	hostname, _ := os.Hostname()
	cur_member, member_change_chan := swim.StartServer(hostname, introducer, verbose) // idk wtd if false positive and swim shutsdown probably not an issue
	node_hash = cur_member.Hash
	files = skiplist.New(skiplist.Uint32)
	members = skiplist.New(skiplist.Uint32)

	go commandLoop()
	go handleMembershipChange(member_change_chan)

	wait_chan := make(chan struct{})
	<-wait_chan
}

func hydfsCreate(local_filename string, hydfs_filename string) (bool, error) {
	if !enoughMembers() {
		return false, fmt.Errorf("Error: %w with current number of members: %d, ", ErrNotEnoughMembers, members.Len())
	}
	file, err := os.Open(local_filename)
	if err != nil {
		log.Println("[ERROR] Open file error:", err)
		return false, nil
	}
	defer file.Close()

	// read contents of file into memory
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("[ERROR] Open file error:", err)
		return false, nil
	}

	target_hash, target := getFileTarget(uint32(hashFilename(hydfs_filename)))
	block := []*repl.FileBlock{{BlockNode: target_hash, BlockID: 0, Data: contents}}
	file_rpc := &repl.File{Filename: hydfs_filename, Blocks: block}
	return sendCreateRPC(target, file_rpc), nil
}

func hydfsGet(filepath string) ([]byte, error) {
	// TODO:
	if !enoughMembers() {
		return nil, fmt.Errorf("Error: %w with current number of members: %d, ", ErrNotEnoughMembers, members.Len())
	}

	return make([]byte, 0), nil
}

func hydfsAppend(filepath string) (bool, error) {
	// TODO:
	if !enoughMembers() {
		return false, fmt.Errorf("Error: %w with current number of members: %d, ", ErrNotEnoughMembers, members.Len())
	}
	return false, nil
}

func hydfsMerge(filepath string) (bool, error) {
	// TODO:
	if !enoughMembers() {
		return false, fmt.Errorf("Error: %w with current number of members: %d, ", ErrNotEnoughMembers, members.Len())
	}
	return false, nil
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
			fmt.Println("Create %v", res)
		case "get":
			if len(commandParts) != 3 {
				fmt.Println("usage: get HyDFSfilename localfilename")
				continue
			}
			hydfs_filename := commandParts[1]
			local_filename := commandParts[2]
			// TODO: get
		case "append":
			if len(commandParts) != 3 {
				fmt.Println("usage: append localfilename HyDFSfilename")
				continue
			}
			local_filename := commandParts[1]
			hydfs_filename := commandParts[2]
			// TODO: append
		case "merge":
			if len(commandParts) != 2 {
				fmt.Println("usage: merge HyDFSfilename")
				continue
			}
			// TODO: merge
		}
	}
}
