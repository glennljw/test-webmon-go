package store

import (
	"sort"
	"sync"
	"time"

	"mini-webmon/internal/record"
)

type Store interface {
	Add(r record.CheckRecord)
	History(url string) []record.CheckRecord
	PurgeBefore(cutoff time.Time) int
}

type Memory struct {
	mu      sync.RWMutex
	records []record.CheckRecord
}

func NewMemory() *Memory {
	return &Memory{}
}

func (m *Memory) Add(r record.CheckRecord) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.records = append(m.records, r)
}

func (m *Memory) History(url string) []record.CheckRecord {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := []record.CheckRecord{}
	for _, r := range m.records {
		if r.URL == url {
			res = append(res, r)
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].CheckedAt.After(res[j].CheckedAt)
	})
	return res
}

func (m *Memory) PurgeBefore(cutoff time.Time) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	newRecords := m.records[:0]
	for _, r := range m.records {
		if r.CheckedAt.After(cutoff) {
			newRecords = append(newRecords, r)
		}
	}
	m.records = newRecords
	return len(m.records) - len(newRecords)
}
