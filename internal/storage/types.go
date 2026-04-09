package storage

type Storage interface {
	GetBucket(key, namespace string) *Bucket
}
