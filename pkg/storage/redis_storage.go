package storage

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go-proxypool/pkg/models"
	"strconv"
	"strings"
)

const (
	key = "proxypool_proxies"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(options *redis.Options) *RedisStorage {
	client := redis.NewClient(options)
	return &RedisStorage{client: client}
}

func (s *RedisStorage) Add(ip models.Ip, ctx context.Context) error {
	proxyStr := fmt.Sprintf("%s:%d", ip.Ip, ip.Port)
	_, err := s.client.SAdd(ctx, key, proxyStr).Result()
	if err != nil {
		return err
	}
	return nil
}

func (s *RedisStorage) Remove(ip models.Ip, ctx context.Context) error {
	proxyStr := fmt.Sprintf("%s:%d", ip.Ip, ip.Port)
	_, err := s.client.SRem(ctx, key, proxyStr).Result()
	if err != nil {
		return err
	}
	return nil
}

func (s *RedisStorage) GetAll(ctx context.Context) ([]models.Ip, error) {
	proxyStrs, err := s.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("could not get IPs from set: %v", err)
	}

	var ips []models.Ip
	for _, proxyStr := range proxyStrs {
		parts := strings.Split(proxyStr, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid proxy string format: %s", proxyStr)
		}
		port, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid port number: %s", parts[1])
		}
		ips = append(ips, models.Ip{Ip: parts[0], Port: port})
	}

	return ips, nil
}
