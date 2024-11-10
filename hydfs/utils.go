package hydfs

import (
	"bytes"
	"cs425/mp3/hydfs/repl"
	"cs425/mp3/hydfs/swim"
	"cs425/mp3/shared"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Prints the files currently stored at this node
// Iterates through the file hashmap to print the names
func listFiles() {
	node_id := members.Get(node_hash).Value.(*shared.MemberInfo).ID
	res := fmt.Sprintf("---------------------------- NODE %d, HASH %d FILES -----------------------\n", node_id, node_hash)
	cur_file := files.Front()
	for range files.Len() {
		filehash := cur_file.Key().(uint32)
		file_struct := cur_file.Value.(*File)
		res += fmt.Sprintf("Hash: %d\tFilename: %s\n", filehash, file_struct.filename)
		res += "Blocks:\n"
		blocks := getBlocks(file_struct.filename, false)
		for _, block := range blocks {
			res += fmt.Sprintf("\tNode: %d\t ID:%d\n", block.BlockNode, block.BlockID)
		}
		res += "\n"
		cur_file = cur_file.Next()
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
func getReplicaIdx(cur_node_hash uint32) uint32 {
	return cur_node_hash % REPL_FACTOR
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
	return members.Len() >= REPL_FACTOR
}

// check if a block exists for a file
func blockExists(filename string, block_node uint32, blockID uint32) bool {
	if !dirExists(HYDFS_DIR + "/" + filename) {
		return false
	}
	return fileExists(fmt.Sprintf("%s/%s/%s", HYDFS_DIR, filename, blockName(filename, block_node, blockID)))
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
	hydfs_log.Println("[INFO] Created block file", block_filepath)
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

// return the replicas minus current node for a given filehash
func getPrimaryReplicas(filehash uint32) []uint32 {
	res := make([]uint32, REPL_FACTOR-1)
	cur_elem := members.Find(filehash)
	i := 0
	for range REPL_FACTOR {
		if cur_elem == nil {
			cur_elem = members.Front()
		}
		cur_elem_hash := cur_elem.Key().(uint32)
		if cur_elem_hash != this_member.Hash {
			res[i] = cur_elem.Key().(uint32)
			i++
		}
		cur_elem = cur_elem.Next()
	}
	return res
}

// return main target node for a file
func getMainFileTarget(file_hash uint32) (uint32, *shared.MemberInfo) {
	mu.Lock()
	defer mu.Unlock()

	main_replica := members.Find(file_hash)
	if main_replica == nil {
		main_replica = members.Front()
	}
	return main_replica.Key().(uint32), main_replica.Value.(*shared.MemberInfo)
}

// return target node for a file
func getReplicaFileTarget(file_hash uint32) (uint32, *shared.MemberInfo) {
	mu.Lock()
	defer mu.Unlock()

	main_replica := members.Find(file_hash)
	if main_replica == nil {
		main_replica = members.Front()
	}
	offset := this_member.Hash % REPL_FACTOR
	for range offset {
		main_replica = main_replica.Next()
		if main_replica == nil {
			main_replica = members.Front()
		}
	}

	return main_replica.Key().(uint32), main_replica.Value.(*shared.MemberInfo)
}

func getVMFileTarget(file_hash uint32, vm_address string) (uint32, *shared.MemberInfo) {
	mu.Lock()
	defer mu.Unlock()

	main_replica := members.Find(file_hash)
	if main_replica == nil {
		main_replica = members.Front()
	}
	for range REPL_FACTOR - 1 {
		curr_member := main_replica.Value.(*shared.MemberInfo)
		if curr_member.Address == vm_address+":9000" {
			break
		}
		main_replica = main_replica.Next()
		if main_replica == nil {
			main_replica = members.Front()
		}
	}
	if main_replica.Value.(*shared.MemberInfo).Address != vm_address+":9000" {
		return 0, nil
	}
	return main_replica.Key().(uint32), main_replica.Value.(*shared.MemberInfo)
}

// get all file hashes that current node is the primary replica for
func getPrimaryReplicaFiles() []uint32 {
	file_hashes := make([]uint32, 0)
	prev_node := members.Get(this_member.Hash).Prev()
	if prev_node == nil {
		prev_node = members.Back()
	}

	start := prev_node.Key().(uint32)
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
			if start < cur_file_hash && cur_file_hash <= end {
				file_hashes = append(file_hashes, cur_file_hash)
			}
		}
		cur_file = cur_file.Next()
	}

	return file_hashes
}

// get all blocks stored at the current node for a given filename
// if data false, dont include data in block structs
func getBlocks(filename string, data bool) []*repl.FileBlock {
	files, err := os.ReadDir(HYDFS_DIR + "/" + filename)
	if err != nil {
		hydfs_log.Fatal("[ERROR] Failed to read directory:", err)
	}

	res_blocks := make([]*repl.FileBlock, 0)
	for _, file := range files {
		res_blocks = append(res_blocks, constructBlock(filename, file.Name(), data))
	}
	return res_blocks
}

// create block struct from block filename
// if data false, dont include data in block struct
func constructBlock(filename string, blockname string, data bool) *repl.FileBlock {
	res_block := &repl.FileBlock{}

	split_string := ".hydfs."
	i := strings.LastIndex(blockname, split_string)
	block_info := strings.Split(blockname[i+len(split_string):], ".")

	block_node, err := strconv.Atoi(block_info[0])
	if err != nil {
		hydfs_log.Fatal("[ERROR] parse int error:", err)
	}
	res_block.BlockNode = uint32(block_node)

	block_id, err := strconv.Atoi(block_info[1])
	if err != nil {
		hydfs_log.Fatal("[ERROR] parse int error:", err)
	}
	res_block.BlockID = uint32(block_id)

	if data {
		res_block.Data = readFile(fmt.Sprintf("%s/%s/%s", HYDFS_DIR, filename, blockname))
	}

	return res_block
}

// read in local_filename into byte slice
func readFile(local_filename string) []byte {
	file, err := os.Open(local_filename)
	if err != nil {
		hydfs_log.Fatal("[ERROR] Open file error:", err)
	}
	defer file.Close()

	// read contents of file into memory
	contents, err := io.ReadAll(file)
	if err != nil {
		hydfs_log.Fatal("[ERROR] Open file error:", err)
	}
	return contents
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

// fill in data to missing files
func fillData(response_missing *repl.RequestMissing) *repl.RequestData {
	request_data := &repl.RequestData{DataFiles: response_missing.MissingFiles}
	for _, rpc_file := range request_data.DataFiles {
		filename := rpc_file.Filename
		for _, rpc_block := range rpc_file.Blocks {
			rpc_block.Data = readFile(blockFilepath(filename, rpc_block.BlockNode, rpc_block.BlockID))
		}
	}
	return request_data
}

func getNumBlocks(response_missing *repl.RequestMissing) int {
	res := 0
	for _, file := range response_missing.MissingFiles {
		res += len(file.Blocks)
	}
	return res
}

// sort blocks by blockNode then blockID
func sortBlocks(file_blocks []*repl.FileBlock) {
	sort.Slice(file_blocks, func(i, j int) bool {
		if file_blocks[i].BlockNode != file_blocks[j].BlockNode {
			return file_blocks[i].BlockNode < file_blocks[j].BlockNode
		}
		return file_blocks[i].BlockID < file_blocks[j].BlockID
	})
}

// allocate []byte and read in all blocks to it
func readAllBlocks(filename string, file_blocks []*repl.FileBlock) []byte {
	var buffer bytes.Buffer
	for _, block := range file_blocks {
		block_data := readFile(blockFilepath(filename, block.BlockNode, block.BlockID))
		buffer.Grow(len(block_data))
		buffer.Write(block_data)
	}
	return buffer.Bytes()
}

// remove files from files list if they are outside range
func garbageCollectFiles() {
	end := node_hash
	start_node := members.Get(node_hash)
	for range REPL_FACTOR {
		start_node = start_node.Prev()
		if start_node == nil {
			start_node = members.Back()
		}
	}
	start := start_node.Key().(uint32)
	for cur_file := files.Front(); cur_file != nil; cur_file = cur_file.Next() {
		cur_file_hash := cur_file.Key().(uint32)
		if start < end && (cur_file_hash <= start || cur_file_hash > end) {
			files.Remove(cur_file_hash)
			os.RemoveAll(HYDFS_DIR + "/" + cur_file.Value.(*File).filename)
		} else if start > end && (end < cur_file_hash && cur_file_hash <= start) {
			files.Remove(cur_file_hash)
			os.RemoveAll(HYDFS_DIR + "/" + cur_file.Value.(*File).filename)
		}
	}
}
