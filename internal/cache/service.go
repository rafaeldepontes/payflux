package cache

import "time"

const DefaultDuration = 24

// Cache implements a internal cache system, every method IS CASE SENSITIVE...
// So "A" and "a" gives different results.
type Cache[K comparable, T any] interface {

	// Add adds something to cache for 48 hours.
	Add(key K, value T)

	// Add adds something to cache, with no TTL were specified it will use the
	// default value of 48 hours.
	AddWithTTL(time time.Duration, key K, value ...T)

	// Set updates the cache value and also refresh the TTL, if the TTL were
	// expired then it removes the old value and create a new one...
	Set(key K, value T)

	// Remove removes...
	Remove(key K)

	// Get gets the value if any and also returns a boolean to check if it is
	// expired.
	Get(key K) (T, bool)

	// Clear clears cache...
	Clear()

	// FullClear should only be called in case of lack of memory, ideally this
	// would never be used since it calls the GC directly and can, depending on
	// the workload, slow down A LOT... Be careful.
	FullClear()
}
