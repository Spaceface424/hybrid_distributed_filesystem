package hydfs

import (
	"cs425/mp3/hydfs/swim"
	"fmt"
	"log"
	"os"
)

// Prints the files currently stored at this node
// Iterates through the file hashmap to print the names
func listFiles() {
	// TODO:
	res := fmt.Sprintf("---------------------------- NODE %d FILES -----------------------------\n", node_hash)
	// for hash, file := range hashes {
	// }
	fmt.Print(res)
}

// Make sure we start with no files stored
// Delete all files in HYDFS_DIR
func cleanup() {
	// TODO:
}

// Get random replica index
// Consistent for same node
func getReplicaIdx(main_idx int, length int, cur_node_hash int32) int {
	return (main_idx + (int(cur_node_hash) % REPL_FACTOR)) % length
}

// Handle events from membership change channel
// Make re-replication calls and re-assign sorted_members ordered map
func handleMembershipChange(member_change_chan chan struct{}) {
	for {
		<-member_change_chan
		swim.Members.Mu.Lock()
		mu.Lock()

		members.Init()
		for node_hash, member := range swim.Members.MemberMap {
			members.Set(node_hash, member)
		}

		mu.Unlock()
		swim.Members.Mu.Unlock()
	}
}

// Hash a filename
func hashFilename(filename string) int32 {
	// TODO: use FNV hash algorithm
	return 0
}

// Block info to filename
func blockName(filename string, block_node uint32, blockID uint32) string {
	return fmt.Sprintf("%s.hydfs.%d.%d", filename, block_node, blockID)
}

// Block info to filepath
func blockFilepath(filename string, block_node uint32, blockID uint32) string {
	return fmt.Sprintf("%s/%s/%s.hydfs.%d.%d", HYDFS_DIR, filename, filename, block_node, blockID)
}

// check if we have at least REPL_FACTOR members
func enoughMembers() bool {
	mu.Lock()
	defer mu.Unlock()

	return members.Len() >= REPL_FACTOR
}

// check if a block exists for a file
func blockExists(filename string, block_node uint32, blockID uint32) bool {
	if !dirExists(HYDFS_DIR + "/" + filename) {
		return false
	}
	return fileExists(blockName(filename, block_node, blockID))
}

// check if a directory exists
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Fatal("[ERROR] stat error:", err)
	}
	return info.IsDir()
}

// check if a file directory exists
func fileDirExists(filename string) bool {
	info, err := os.Stat(HYDFS_DIR + "/" + filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Fatal("[ERROR] stat error:", err)
	}
	return info.IsDir()
}

// check if a file exists
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Fatal("[ERROR] stat error:", err)
	}
	return !info.IsDir()
}

// create a directory for a new file in Hydfs
func createDir(filename string) {
	err := os.MkdirAll(HYDFS_DIR+"/"+filename, 0644)
	if err != nil {
		log.Fatal("[ERROR] mkdir error:", err)
	}
	return
}

// create a block file and write data into it
func createBlock(filename string, block_node uint32, blockID uint32, data []byte) {
	block_filepath := blockFilepath(filename, block_node, blockID)
	// create the directory for the file if it doesn't exist
	if !dirExists(filename) {
		createDir(filename)
	}
	if fileExists(block_filepath) {
		log.Printf("[WARNING] Trying to create block that already exists: %s", block_filepath)
	}
	file, err := os.Create(block_filepath)
	if err != nil {
		log.Fatal("[ERROR] Create file error:", err)
	}
	_, err = file.Write(data)
	if err != nil {
		log.Fatal("[ERROR] Write file error:", err)
	}
}

// return next REPL_FACTOR - 1 nodes in the hashring
func getReplicas() []uint32 {
	res := make([]uint32, REPL_FACTOR-1)
	cur_elem := members.Get(node_hash).Next()
	for i := range REPL_FACTOR - 1 {
		if cur_elem == nil {
			cur_elem = members.Front()
		}
		res[i] = cur_elem.Key().(uint32)
		cur_elem = cur_elem.Next()
	}
	return res
}
