package hydfs

import "hash/fnv"

type HashRing struct {
	hash_func hash.Hash32

}

func CreateHashRing() HashRing {
  hr := HashRing{fnv.New32()}
  return hr
}

func (hr HashRing) AddNode(data byte[]) {

}

func (hr HashRing) DeleteNode()

