package storage

import (
	"context"
	"go-proxypool/pkg/models"
	"sync"
)

type MemoryStorage struct {
	ips *[]models.Ip
	mu  sync.RWMutex
}

// NewMemoryStorage returns a new MemoryStorage.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		ips: &[]models.Ip{},
	}
}

func (s *MemoryStorage) Add(ip models.Ip, ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	*s.ips = append(*s.ips, ip)
	return nil
}

func (s *MemoryStorage) Remove(ip models.Ip, ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, v := range *s.ips {
		if v.Ip == ip.Ip && v.Port == ip.Port {
			*s.ips = append((*s.ips)[:i], (*s.ips)[i+1:]...)
			return nil
		}
	}
	return nil
}

func (s *MemoryStorage) GetAll(ctx context.Context) ([]models.Ip, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return *s.ips, nil
}
