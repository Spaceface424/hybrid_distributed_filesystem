package hydfs

import (
	"bufio"
	"cs425/mp3/hydfs/swim"
	"os"
	"strings"
)

func StartHydfs() {
	cleanup()
	cur_member, member_change_chan := swim.StartServer()
	node_hash := cur_member.ID
}

func Create(filepath string) bool {
	// TODO:
	return false
}

func Get(filepath string) []byte {
	// TODO:
	return make([]byte, 0)
}

func Append(filepath string) bool {
	// TODO:
	return false
}

func Merge(filepath string) bool {
	// TODO:
	return false
}

func send_create_request_main() bool {
	// TODO:
	// client = gRPC setup
	// target_node = get_replica()
	//
	return false
}

func send_create_request_replica() bool {
	// TODO:
	return false
}

func handle_create_request(filename string, data []byte) bool {
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
			list_files()
			// TODO: rest of commands
		}
	}
}
