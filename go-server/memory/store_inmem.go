package memory

import (
	"errors"
	"strings"
	"sync"
	"time"
)

type InMemoryStore struct {
	mu       sync.RWMutex
	memories []*Memory
}

func NewInMemStore() *InMemoryStore {
	return &InMemoryStore{memories: []*Memory{}}
}

func (s *InMemoryStore) Insert(m *Memory) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.memories = append(s.memories, m)
	return nil
}

func (s *InMemoryStore) QueryRelevant(userID, query string, k int, kinds []Kind) ([]*Memory, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	q := strings.ToLower(query)
	kindSet := map[Kind]bool{}
	for _, kd := range kinds {
		kindSet[kd] = true
	}

	var out []*Memory
	for _, m := range s.memories {
		if m.UserID != userID {
			continue
		}
		if len(kindSet) > 0 && !kindSet[m.Kind] {
			continue
		}
		if q == "" || strings.Contains(strings.ToLower(m.Content), q) {
			out = append(out, m)
		}
		if len(out) >= k {
			break
		}
	}
	now := time.Now().UTC()
	for _, m := range out {
		m.LastUsedAt = now
	}

	return out, nil
}

func (s *InMemoryStore) Update(updated *Memory) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, m := range s.memories {
		if m.ID == updated.ID && m.UserID == updated.UserID {
			s.memories[i] = updated
			return nil
		}
	}
	return errors.New("memory not found")
}

func (s *InMemoryStore) Delete(userID, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, m := range s.memories {
		if m.ID == id && m.UserID == userID {
			s.memories = append(s.memories[:i], s.memories[i+1:]...)
			return nil
		}
	}
	return errors.New("memory not found")
}
