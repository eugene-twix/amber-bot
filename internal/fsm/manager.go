// internal/fsm/manager.go
package fsm

import (
	"context"
	"fmt"
	"time"

	"github.com/eugene-twix/amber-bot/internal/cache"
)

const fsmTTL = 10 * time.Minute

type Manager struct {
	cache *cache.Cache
}

func NewManager(cache *cache.Cache) *Manager {
	return &Manager{cache: cache}
}

type UserState struct {
	State State `json:"state"`
	Data  Data  `json:"data"`
}

func (m *Manager) key(userID int64) string {
	return fmt.Sprintf("fsm:%d", userID)
}

func (m *Manager) Get(ctx context.Context, userID int64) (*UserState, error) {
	var state UserState
	err := m.cache.Get(ctx, m.key(userID), &state)
	if err != nil {
		return &UserState{State: StateNone, Data: make(Data)}, nil
	}
	if state.Data == nil {
		state.Data = make(Data)
	}
	return &state, nil
}

func (m *Manager) Set(ctx context.Context, userID int64, state State, data Data) error {
	return m.cache.Set(ctx, m.key(userID), &UserState{State: state, Data: data}, fsmTTL)
}

func (m *Manager) Clear(ctx context.Context, userID int64) error {
	return m.cache.Delete(ctx, m.key(userID))
}

func (m *Manager) Update(ctx context.Context, userID int64, state State, key string, value any) error {
	current, _ := m.Get(ctx, userID)
	current.Data[key] = value
	current.State = state
	return m.cache.Set(ctx, m.key(userID), current, fsmTTL)
}
