package hydfs

import (
	"github.com/huandu/skiplist"
)

// lru cache is a skiplist that holds filehash as key and filedata as value
type LRU_cache struct {
	list  *skiplist.SkipList
	front *node
	back  *node
	size  int // cache size in bytes
	limit int // cache limit in bytes
}

type node struct {
	next *node
	prev *node
	key  uint32
}

func createCache(cache_size int) *LRU_cache {
	return &LRU_cache{
		list:  skiplist.New(skiplist.Uint32),
		front: nil,
		back:  nil,
		size:  0,
		limit: cache_size,
	}
}

func (cache *LRU_cache) addFile(filename string, data []byte) {
	if len(data) > cache.limit {
		hydfs_log.Printf("[WARNING] file %s too big for cache", filename)
		return
	}
	filehash := hashFilename(filename)
	for cache.size+len(data) > cache.limit {
		cache.evict()
	}
	if cache.front == nil && cache.back == nil {
		new_head := &node{key: filehash, next: nil, prev: nil}
		cache.front = new_head
		cache.back = new_head
		cache.list.Set(filehash, data)
		cache.size += len(data)
		hydfs_log.Printf("[INFO] Adding file %s to cache, cache size = %.2f KB", filename, float64(cache.size)/1024)
		return
	}
	if cache.front == cache.back {
		new_head := &node{key: filehash, next: cache.front, prev: nil}
		cache.front = new_head
		cache.back.prev = new_head
		cache.list.Set(filehash, data)
		cache.size += len(data)
		hydfs_log.Printf("[INFO] Adding file %s to cache, cache size = %.2f KB", filename, float64(cache.size)/1024)
		return
	}
	new_head := &node{key: filehash, next: cache.front, prev: nil}
	cache.front.prev = new_head
	cache.front = new_head
	cache.list.Set(filehash, data)
	cache.size += len(data)
	hydfs_log.Printf("[INFO] Adding file %s to cache, cache size = %.2f KB", filename, float64(cache.size)/1024)
}

func (cache *LRU_cache) evict() {
	if cache.front == nil && cache.back == nil {
		return
	}
	hydfs_log.Printf("[INFO] GET exceeded cache capacity, deleting oldest entry")
	evict_key := cache.back.key
	evict_size := len(cache.list.Get(evict_key).Value.([]byte))
	if cache.front == cache.back {
		cache.front = nil
		cache.back = nil
		cache.size -= evict_size
		cache.list.Remove(evict_key)
		return
	}
	cache.back = cache.back.prev
	cache.back.next = nil
	cache.size -= evict_size
	cache.list.Remove(evict_key)
}

func (cache *LRU_cache) getFile(filename string) []byte {
	file_element := cache.list.Get(hashFilename(filename))
	if file_element == nil {
		return nil
	}
	return file_element.Value.([]byte)
}

func (cache *LRU_cache) invalidateFile(filename string) {
	filehash := hashFilename(filename)
	if cache.list.Get(filehash) == nil {
		return
	}
	hydfs_log.Printf("[INFO] Invalidating file %s from cache", filename)
	for cur_node := cache.front; cur_node != nil; cur_node = cur_node.next {
		if cur_node.key == filehash {
			switch cur_node {
			case cache.front:
				cache.front = cur_node.next
				cache.front.prev = nil
			case cache.back:
				cache.back = cache.back.prev
				cache.back.next = nil
			default:
				cur_node.prev.next = cur_node.next
				cur_node.next.prev = cur_node.prev
			}
			break
		}
	}
	evict_element := cache.list.Get(filehash)
	evict_size := len(evict_element.Value.([]byte))
	cache.size -= evict_size
	cache.list.Remove(filehash)
}

func (cache *LRU_cache) clear() {
	cache.list.Init()
	cache.front = nil
	cache.back = nil
	cache.size = 0
}

func (cache *LRU_cache) setLimit(new_limit_kb int) {
	cache.clear()
	cache.limit = new_limit_kb * 1000
}
