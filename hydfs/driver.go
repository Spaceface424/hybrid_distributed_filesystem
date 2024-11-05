package hydfs

import (
	"bufio"
	"cs425/mp3/hydfs/swim"
	"fmt"
	"os"
	"strings"

	"github.com/huandu/skiplist"
)

func StartHydfs() {
	cleanup()
	cur_member, member_change_chan := swim.StartServer() // idk wtd if false positive and swim shutsdown probably not an issue
	node_hash := cur_member.ID
	files = skiplist.New(skiplist.Uint32)
	members = skiplist.New(skiplist.Uint32)
}

func Create(filepath string) (bool, error) {
	// TODO:
	if !enoughMembers() {
		return false, fmt.Errorf("Error: %w with current number of members: %d, ", ErrNotEnoughMembers, members.Len())
	}
	return false, nil
}

func Get(filepath string) ([]byte, error) {
	// TODO:
	if !enoughMembers() {
		return nil, fmt.Errorf("Error: %w with current number of members: %d, ", ErrNotEnoughMembers, members.Len())
	}
	return make([]byte, 0), nil
}

func Append(filepath string) (bool, error) {
	// TODO:
	if !enoughMembers() {
		return false, fmt.Errorf("Error: %w with current number of members: %d, ", ErrNotEnoughMembers, members.Len())
	}
	return false, nil
}

func Merge(filepath string) (bool, error) {
	// TODO:
	if !enoughMembers() {
		return false, fmt.Errorf("Error: %w with current number of members: %d, ", ErrNotEnoughMembers, members.Len())
	}
	return false, nil
}

func sendCreateRequestMain() bool {
	// TODO:
	// client = gRPC setup
	// target_node = get_replica()
	//
	return false
}

func sendCreateRequestReplica() bool {
	// TODO:
	return false
}

func handleCreateRequest(filename string, data []byte) bool {
	// TODO:
	return false
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
			// TODO: rest of commands
		}
	}
}
