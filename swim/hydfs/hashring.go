package hydfs

import "hash"

type HashRing struct {
	hash_func hash.Hash32
}

func (HashRing) AddNode() {
}
