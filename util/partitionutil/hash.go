package partitionutil

import (
	"hash/fnv"
)

// See https://github.com/gholt/ring/blob/master/BASIC_HASH_RING.md

func hash(data []byte) uint64 {
	hasher := fnv.New64a()
	_, _ = hasher.Write(data)
	return hasher.Sum64()
}

func HashPartition(data []byte, numPartitions int) int {
	h := hash(data)
	return int(h % uint64(numPartitions))
}
