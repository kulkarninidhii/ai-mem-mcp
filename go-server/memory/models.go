package memory

import "time"

type Kind string

const (
	KindProfile  Kind = "profile"
	KindEpisodic Kind = "episodic"
	KindTask     Kind = "task"
)

type Memory struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Kind       Kind      `json:"kind"`
	Content    string    `json:"content"`
	Importance float64   `json:"importance"`
	Tags       []string  `json:"tags"`
	CreatedAt  time.Time `json:"created_at"`
	LastUsedAt time.Time `json:"last_used_at"`
}

type Store interface {
	Insert(m *Memory) error
	QueryRelevant(userID, query string, k int, kinds []Kind) ([]*Memory, error)
	Update(m *Memory) error
	Delete(userID, id string) error
}
