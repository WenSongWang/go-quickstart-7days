package store

import (
	"context"
	"encoding/json"
	"time"
)

// Cache 定义一个最小的缓存接口，便于测试与替换实现（Redis / 内存 / 空实现）
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

// CachedUserStore：内存 store + 缓存层（Redis）
// - 读：先读缓存，miss 再读底层 store 并回填缓存（Cache-Aside）
// - 写：写底层 store 后**删缓存**（写后删），下次读时再回填，避免写穿不一致
type CachedUserStore struct {
	Base  *MemoryStore
	Cache Cache

	// ListTTL 列表缓存 TTL
	ListTTL time.Duration
}

const usersListCacheKey = "users:all"

func (s *CachedUserStore) List(ctx context.Context) ([]User, bool, error) {
	if s.Cache != nil {
		if raw, err := s.Cache.Get(ctx, usersListCacheKey); err == nil && raw != "" {
			var users []User
			if jsonErr := json.Unmarshal([]byte(raw), &users); jsonErr == nil {
				return users, true, nil
			}
		}
	}

	users := s.Base.List()
	if s.Cache != nil {
		if b, err := json.Marshal(users); err == nil {
			_ = s.Cache.Set(ctx, usersListCacheKey, string(b), s.ListTTL)
		}
	}
	return users, false, nil
}

func (s *CachedUserStore) Create(ctx context.Context, name string) (User, error) {
	u := s.Base.Create(name)
	// 写后删：写 DB/store 后删除列表缓存 key，下次 List 时 cache miss 再回填（避免长期脏读）
	if s.Cache != nil {
		_ = s.Cache.Delete(ctx, usersListCacheKey)
	}
	return u, nil
}

