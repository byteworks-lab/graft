package store

type StoreInf interface {
	Put(key any, value any)
	Get(key any) (any, bool)
	Delete(key any) bool
}
