package shardedmap

import (
	"crypto/sha1"
	"sync"
)

type Shard[T any] struct {
	mu sync.RWMutex
	m  map[string]T
}

type ShardedMap[T any] []*Shard[T]

func New[T any](shardCount int) ShardedMap[T] {
	shards := make([]*Shard[T], shardCount)

	for i := range shardCount {
		shardMap := make(map[string]T)
		shards[i] = &Shard[T]{m: shardMap}
	}

	return shards
}

func (sm *ShardedMap[T]) shardIndex(key string) int {
	checksum := sha1.Sum([]byte(key))
	index := int(checksum[10])
	return index % len(*sm)
}

func (sm *ShardedMap[T]) getShard(key string) *Shard[T] {
	index := sm.shardIndex(key)
	return (*sm)[index]
}

func (sm *ShardedMap[T]) Get(key string) (T, bool) {
	shard := sm.getShard(key)
	shard.mu.RLock()
	defer shard.mu.RUnlock()

	value, ok := shard.m[key]
	return value, ok
}

func (sm *ShardedMap[T]) Put(key string, value T) {
	shard := sm.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	shard.m[key] = value
}
