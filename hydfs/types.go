package hydfs

import (
	"fmt"
	"sync"

	"github.com/huandu/skiplist"
)

const (
	CACHE_DIR   = "cache/"     // cache directory to store cached reads
	HYDFS_DIR   = "hydfs_dir/" // hydfs directory to store blocks
	REPL_FACTOR = 3
)

var (
	files     *skiplist.SkipList // ordered map hash to file
	members   *skiplist.SkipList // ordered map of current nodes sorted by node hashes
	mu        sync.Mutex
	node_hash uint32 // current node hash
)

var ErrNotEnoughMembers error = fmt.Errorf("Not enough members to satisfy replication factor of %d", REPL_FACTOR)

type File struct {
	filename string // user defined filename
	nextID   uint32 // monotonically increasing local fileblock ID
}
