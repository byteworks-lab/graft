package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestLRUCacheSetValueExists verifies that a value inserted into the cache
// can be retrieved using its corresponding key.
func TestLRUCacheSetValueExists(t *testing.T) {
	cache := NewCache(4)
	cache.Put("12", "Twelve")
	cache.Put(14, 14)
	val, _ := cache.Get(14)
	assert.Equal(t, 14, val.(int), "Value fetched from the cache should be 'Twelve'")
}

// TestLRUCacheSetValueNotExists ensures that requesting a non-existent key returns nil.
func TestLRUCacheSetValueNotExists(t *testing.T) {
	cache := NewCache(4)
	cache.Put("12", "Twelve")
	val, _ := cache.Get("11")
	assert.Nil(t, val, "Value should be nil as key '11' does not exist")
}

// TestLRUCacheKeyEvicted checks that the least recently used item is evicted
// when the cache exceeds its capacity.
func TestLRUCacheKeyEvicted(t *testing.T) {
	cache := NewCache(4)
	cache.Put("12", "Twelve")
	cache.Put("11", "Eleven")
	cache.Put("10", "Ten")
	cache.Put("9", "Nine")
	cache.Put("8", "Eight") // Should evict key "12"
	val, _ := cache.Get("12")
	assert.Nil(t, val, "Key '12' should be evicted from the cache due to LRU policy")
}
