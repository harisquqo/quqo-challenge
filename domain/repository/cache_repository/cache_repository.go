package cache_repository

import "time"

type CacheRepository interface {
	SetKey(key string, value interface{}, expiration time.Duration) error
	GetKey(key string, src interface{}) error
	DelKey(key string) error
}
