package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestTTLCache verifies that a key can be inserted and immediately retrieved
// before it expires from the cache.
func TestTTLCache(t *testing.T) {
	cache := NewCache(4, time.Duration(200), time.Duration(100)) // TTL = 200ms, cleanup interval = 100ms
	cache.Put(10, 12)
	val, _ := cache.Get(10)
	assert.Equal(t, 12, val, "Value should be retrievable before TTL expiration")
}

// TestTTLCacheExpiredEntries checks that a value is still retrievable before TTL expiry,
// and confirms that subsequent entries do not prematurely remove earlier ones.
func TestTTLCacheExpiredEntries(t *testing.T) {
	cache := NewCache(4, time.Duration(1000), time.Duration(200)) // TTL = 1000ms, cleanup interval = 200ms

	cache.Put(10, 12)              // Add first entry
	time.Sleep(time.Duration(200)) // Wait before adding the next entry

	cache.Put(12, 26)       // Add second entry after delay
	val, _ := cache.Get(10) // First value should still be present
	assert.Equal(t, 12, val, "Key 10 should not be expired yet")

	val, _ = cache.Get(12) // Second value should also be present
	assert.Equal(t, 26, val, "Key 12 should be retrievable immediately after insertion")
}
