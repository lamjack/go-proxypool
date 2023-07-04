package storage

import (
	"context"
	"go-proxypool/pkg/models"
)

type Storage interface {
	// Add adds an ip to storage.
	Add(ip models.Ip, ctx context.Context) error

	// Remove removes an ip from storage.
	Remove(ip models.Ip, ctx context.Context) error

	// GetAll returns all ips from storage.
	GetAll(ctx context.Context) ([]models.Ip, error)
}
