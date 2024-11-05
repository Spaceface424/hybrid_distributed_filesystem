package hydfs

import "github.com/huandu/skiplist"

const (
	CACHE_DIR   = "cache/" // cache directory to store cached reads
	HYDFS_DIR   = "tmp/"   // hydfs directory to store blocks
	REPL_FACTOR = 3
)

var (
	hashes map[int32]file // map filename to hash
	// sorted_members []*shared.MemberInfo // slice of current nodes sorted by node hashes
	sorted_members *skiplist.SkipList // ordered map of current nodes sorted by node hashes
	node_hash      int32              // current node hash
)

type file struct {
	filename string // user defined filename
	nextID   int32  // monotonically increasing local fileblock ID
}
