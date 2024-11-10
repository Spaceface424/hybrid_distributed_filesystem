package hydfs

import (
	"cs425/mp3/shared"
	"fmt"
	"hash"
	"log"
	"sync"

	"github.com/huandu/skiplist"
)

const (
	CACHE_DIR   = "cache"     // cache directory to store cached reads
	HYDFS_DIR   = "hydfs_dir" // hydfs directory to store blocks
	REPL_FACTOR = 3
	GRPC_PORT   = "7000"
	M           = 15
)

var (
	files        *skiplist.SkipList // ordered map hash to file
	members      *skiplist.SkipList // ordered map of current nodes sorted by node hashes
	cache        *LRU_cache
	enable_cache bool
	mu           sync.Mutex
	node_hash    uint32 // current node hash
	this_member  *shared.MemberInfo
	hash_func    hash.Hash32
	hydfs_log    *log.Logger
)

var ErrNotEnoughMembers error = fmt.Errorf("Not enough members to satisfy replication factor of %d", REPL_FACTOR)

var VM = map[string]string{
	"m1":  "fa24-cs425-1401.cs.illinois.edu",
	"m2":  "fa24-cs425-1402.cs.illinois.edu",
	"m3":  "fa24-cs425-1403.cs.illinois.edu",
	"m4":  "fa24-cs425-1404.cs.illinois.edu",
	"m5":  "fa24-cs425-1405.cs.illinois.edu",
	"m6":  "fa24-cs425-1406.cs.illinois.edu",
	"m7":  "fa24-cs425-1407.cs.illinois.edu",
	"m8":  "fa24-cs425-1408.cs.illinois.edu",
	"m9":  "fa24-cs425-1409.cs.illinois.edu",
	"m10": "fa24-cs425-1410.cs.illinois.edu",
}

type File struct {
	filename string // user defined filename
	nextID   uint32 // monotonically increasing local fileblock ID
}

type CachedFile struct {
	timestamp uint32 // strictly increasing id for cache removal
	file      []byte
}
