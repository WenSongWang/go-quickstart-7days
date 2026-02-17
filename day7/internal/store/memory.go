package store

import (
	"sync"
	"sync/atomic"
)

// User 用户实体，和 handler 返回的 JSON 一致
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// MemoryStore 用 map 在内存里存用户，Day7 演示用；实际项目可换成数据库
type MemoryStore struct {
	mu    sync.RWMutex  // 读多写少用 RLock/RUnlock，写用 Lock
	next  int64         // 自增 ID，用 atomic 保证并发安全
	users map[int]User
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{users: map[int]User{}}
}

// List 返回所有用户（复制一份，避免外部改 map）
func (s *MemoryStore) List() []User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]User, 0, len(s.users))
	for _, u := range s.users {
		out = append(out, u)
	}
	return out
}

// Get 根据 id 查一个用户，不存在返回 false
func (s *MemoryStore) Get(id int) (User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	return u, ok
}

// Create 创建用户，自增 ID，返回创建好的用户
func (s *MemoryStore) Create(name string) User {
	id := int(atomic.AddInt64(&s.next, 1))
	s.mu.Lock()
	defer s.mu.Unlock()
	u := User{ID: id, Name: name}
	s.users[id] = u
	return u
}
