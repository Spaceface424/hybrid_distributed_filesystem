package hydfs

import "fmt"

// Prints the files currently stored at this node
// Iterates through the file hashmap to print the names
func list_files() {
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
func get_replica_idx(main_idx int, length int, cur_node_hash int32) int {
	return (main_idx + (int(cur_node_hash) % REPL_FACTOR)) % length
}

// Handle events from membership change channel
// Make re-replication calls and re-assign sorted_members slice
func handle_membership_change() {
	// TODO: loop forever and read from membership change channel
	for {
	}
}

// Hash a filename
func hash_filename(filename string) int32 {
	// TODO: use FNV hash algorithm
	return 0
}

// Block info to filepath
func block_filepath(filename string, block_node int32, blockID int32) string {
	return fmt.Sprintf("%s.hydfs.%s.%s")
}
