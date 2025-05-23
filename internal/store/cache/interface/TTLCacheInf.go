package _interface

type TTLCacheInf interface {
	Cache
	Evict(key any)
}
