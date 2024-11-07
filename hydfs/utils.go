package hydfs

import (
	"cs425/mp3/hydfs/swim"
	"cs425/mp3/shared"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Prints the files currently stored at this node
// Iterates through the file hashmap to print the names
func listFiles() {
	node_id := members.Get(node_hash).Value.(*shared.MemberInfo).ID
	res := fmt.Sprintf("---------------------------- NODE %d, HASH %d FILES -----------------------\n", node_id, node_hash)
	cur_file := files.Front()
	for range files.Len() {
		res += fmt.Sprintf("Hash: %d\tFilename: %s\n", cur_file.Key().(uint32), cur_file.Element().Value.(File).filename)
	}
	res += "--------------------------------------------------------------------------\n"
	fmt.Print(res)
}

// Make sure we start with no files stored
// Delete all files in HYDFS_DIR or create it if it doesn't exist
func cleanupDir() {
	files, err := os.ReadDir(HYDFS_DIR)
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(HYDFS_DIR, 0644)
			return
		} else {
			hydfs_log.Fatal("[ERROR] Failed to read directory:", err)
		}
	}

	for _, file := range files {
		filePath := filepath.Join(HYDFS_DIR, file.Name())
		if err := os.RemoveAll(filePath); err != nil {
			hydfs_log.Fatalf("[ERROR] Failed to delete file %s: %v", filePath, err)
		}
	}
}

// Get random replica index
// Consistent for same node
func getReplicaIdx(main_idx int, length int, cur_node_hash int32) int {
	return (main_idx + (int(cur_node_hash) % REPL_FACTOR)) % length
}

// Hash a filename
func hashFilename(filename string) uint32 {
	io.WriteString(hash_func, filename)
	hash_val := hash_func.Sum32()
	hash_func.Reset()
	return hash_val % (2 << swim.M)
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
		hydfs_log.Fatal("[ERROR] stat error:", err)
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
		hydfs_log.Fatal("[ERROR] stat error:", err)
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
		hydfs_log.Fatal("[ERROR] stat error:", err)
	}
	return !info.IsDir()
}

// create a directory for a new file in Hydfs
func createDir(filename string) {
	err := os.MkdirAll(HYDFS_DIR+"/"+filename, 0777)
	if err != nil {
		hydfs_log.Fatal("[ERROR] mkdir error:", err)
	}
}

// create a block file and write data into it
func createBlock(filename string, block_node uint32, blockID uint32, data []byte) {
	block_filepath := blockFilepath(filename, block_node, blockID)
	// create the directory for the file if it doesn't exist
	if !dirExists(filename) {
		createDir(filename)
	}
	if fileExists(block_filepath) {
		hydfs_log.Printf("[WARNING] Trying to create block that already exists: %s", block_filepath)
	}
	file, err := os.Create(block_filepath)
	if err != nil {
		hydfs_log.Fatal("[ERROR] Create file error:", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		hydfs_log.Fatal("[ERROR] Write file error:", err)
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
	hydfs_log.Printf("[INFO] Replicating create request to replicas: %v", res)
	return res
}

// return target node for a file
func getFileTarget(file_hash uint32) (uint32, *shared.MemberInfo) {
	mu.Lock()
	defer mu.Unlock()

	main_replica := members.Find(file_hash)
	if main_replica == nil {
		main_replica = members.Front()
	}
	offset := node_hash % REPL_FACTOR
	for range offset {
		main_replica = main_replica.Next()
		if main_replica == nil {
			main_replica = members.Front()
		}
	}

	return main_replica.Key().(uint32), main_replica.Value.(*shared.MemberInfo)
}

// get all file hashes that current node is the primary replica for
func getPrimaryReplicaFiles() {
	file_hashes := make([]uint32, 0)
	prev_node := members.Get(this_member.Hash).Prev()
	if prev_node == nil {
		prev_node = members.Back()
	}

	start := (prev_node.Key().(uint32) + 1) % (2 << M)
	end := this_member.Hash
	cur_file := files.Front()
	for cur_file != nil {
		cur_file_hash := cur_file.Key().(uint32)
		// if cur_node and prev_node between 0
		if start > end {
			if cur_file_hash < start || cur_file_hash >= end {
				file_hashes = append(file_hashes, cur_file_hash)
			}
		} else {
			if cur_file_hash < start && cur_file_hash >= end {
				file_hashes = append(file_hashes, cur_file_hash)
			}
		}
	}
}

func printMemberDict() {
	mu.Lock()
	defer mu.Unlock()

	res := "--------- Mem Dict ----------\n"
	cur := members.Front()
	for cur != nil {
		res += fmt.Sprintf("Hash: %d\tNode: %d\n", cur.Key(), cur.Value.(*shared.MemberInfo).ID)
		// res += fmt.Sprintf("Hash: %d\n", cur.Key())
		cur = cur.Next()
	}
	res += "-----------------------------\n"
	fmt.Print(res)
}
