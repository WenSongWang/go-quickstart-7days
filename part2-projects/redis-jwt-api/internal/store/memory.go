package store

import (
	"sync"
	"sync/atomic"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// MemoryStore：先用内存跑通业务，再替换为 DB；同时可叠加 Redis 缓存层
type MemoryStore struct {
	mu    sync.RWMutex
	next  int64
	users map[int]User
}

func NewMemoryStore() *MemoryStore {
	s := &MemoryStore{users: map[int]User{}}
	// 预置两条数据，方便直接 GET /api/users 看效果
	s.Create("张三")
	s.Create("李四")
	return s
}

func (s *MemoryStore) List() []User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]User, 0, len(s.users))
	for _, u := range s.users {
		out = append(out, u)
	}
	return out
}

func (s *MemoryStore) Create(name string) User {
	id := int(atomic.AddInt64(&s.next, 1))
	s.mu.Lock()
	defer s.mu.Unlock()
	u := User{ID: id, Name: name}
	s.users[id] = u
	return u
}

