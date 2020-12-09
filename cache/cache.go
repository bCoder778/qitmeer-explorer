package cache

import (
	"fmt"
	"github.com/muesli/cache2go"
	"sync"
	"time"
)

type MemCache struct {
	cacheTables map[string]*cache2go.CacheTable
	mutex       sync.RWMutex
}

func NewMemCache() *MemCache {
	return &MemCache{
		cacheTables: make(map[string]*cache2go.CacheTable),
	}
}

func (m *MemCache) Add(table string, key string, sec time.Duration, iTerm interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	t, ok := m.cacheTables[table]
	if !ok {
		cacheTable := cache2go.Cache(table)
		cacheTable.Add(key, sec, iTerm)
		m.cacheTables[table] = cacheTable
	} else {
		t.Add(key, sec, iTerm)
	}
}

func (m *MemCache) Value(table, key string) (interface{}, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	t, ok := m.cacheTables[table]
	if !ok {
		return nil, fmt.Errorf("not table %s", table)
	} else {
		iTerm, err := t.Value(key)
		if err != nil {
			return nil, err
		}
		return iTerm.Data(), nil
	}
}
