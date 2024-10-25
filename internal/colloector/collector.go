package colloector

import (
	"context"
	"fmt"
	"redis-monitor/internal/models"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type MetricCollector interface {
	CollectMemoryMetrics(ctx context.Context) (*models.MemoryMetrics, error)
	CollectConnectionMetrics(ctx context.Context) (*models.ConnectionMetrics, error)
	CollectCacheMetrics(ctx context.Context) (*models.CacheMetrics, error)
	CollectAll(ctx context.Context) (*models.MetricsCollection, error)
}

type RedisCollector struct {
	client   *redis.Client
	instance string
}

func NewRedisCollector(addr, password string) *RedisCollector {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	return &RedisCollector{
		client:   client,
		instance: addr,
	}
}

func (rc *RedisCollector) CollectMemoryMetrics(ctx context.Context) (*models.MemoryMetrics, error) {
	// INFO MEMORY 명령어 실행
	result, err := rc.client.Info(ctx, "memory").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to collect memory metrics: %v", err)
	}

	// 결과 파싱을 위한 메트릭 객체 생성
	metrics := &models.MemoryMetrics{
		BaseMetric: models.BaseMetric{
			Timestamp: time.Now(),
			Type:      models.MemoryMetric,
			Instance:  rc.instance,
		},
	}

	// INFO 결과 파싱
	for _, line := range strings.Split(result, "\n") {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}

		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])

		switch key {
		case "used_memory":
			metrics.UsedMemoryBytes, _ = strconv.ParseInt(value, 10, 64)
		case "maxmemory":
			metrics.MaxMemoryBytes, _ = strconv.ParseInt(value, 10, 64)
		}
	}

	// 메모리 사용률 계산
	if metrics.MaxMemoryBytes > 0 {
		metrics.MemoryUsagePerc = float64(metrics.UsedMemoryBytes) / float64(metrics.MaxMemoryBytes) * 100
	}

	return metrics, nil
}

func (rc *RedisCollector) CollectConnectionMetrics(ctx context.Context) (*models.ConnectionMetrics, error) {
	// INFO CLIENTS 명령어 실행
	clientsInfo, err := rc.client.Info(ctx, "clients").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to collect clients info: %v", err)
	}

	// INFO REPLICATION 명령어로 슬레이브 정보 수집
	replInfo, err := rc.client.Info(ctx, "replication").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to collect replication info: %v", err)
	}

	metrics := &models.ConnectionMetrics{
		BaseMetric: models.BaseMetric{
			Timestamp: time.Now(),
			Type:      models.ConnectionMetric,
			Instance:  rc.instance,
		},
	}

	// Clients 정보 파싱
	for _, line := range strings.Split(clientsInfo, "\n") {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}

		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		switch key {
		case "connected_clients":
			metrics.ConnectedClients, _ = strconv.ParseInt(value, 10, 64)
		case "blocked_clients":
			metrics.BlockedClients, _ = strconv.ParseInt(value, 10, 64)
		case "maxclients":
			metrics.MaxClients, _ = strconv.ParseInt(value, 10, 64)
		}
	}

	// Replication 정보에서 슬레이브 수 파싱
	for _, line := range strings.Split(replInfo, "\n") {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}

		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		if key == "connected_slaves" {
			metrics.ConnectedSlaves, _ = strconv.ParseInt(value, 10, 64)
		}
	}

	return metrics, nil
}

func (rc *RedisCollector) CollectCacheMetrics(ctx context.Context) (*models.CacheMetrics, error) {
	// INFO STATS 명령어 실행
	result, err := rc.client.Info(ctx, "stats").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to collect cache metrics: %v", err)
	}

	metrics := &models.CacheMetrics{
		BaseMetric: models.BaseMetric{
			Timestamp: time.Now(),
			Type:      models.CacheMetric,
			Instance:  rc.instance,
		},
	}

	// Stats 정보 파싱
	for _, line := range strings.Split(result, "\n") {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}

		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		switch key {
		case "keyspace_hits":
			metrics.KeyspaceHits, _ = strconv.ParseInt(value, 10, 64)
		case "keyspace_misses":
			metrics.KeyspaceMisses, _ = strconv.ParseInt(value, 10, 64)
		case "evicted_keys":
			metrics.Evictions, _ = strconv.ParseInt(value, 10, 64)
		case "expired_keys":
			metrics.ExpiredKeys, _ = strconv.ParseInt(value, 10, 64)
		}
	}

	// Hit rate 계산
	totalOperations := metrics.KeyspaceHits + metrics.KeyspaceMisses
	if totalOperations > 0 {
		metrics.HitRate = float64(metrics.KeyspaceHits) / float64(totalOperations) * 100
	}

	return metrics, nil
}

func (rc *RedisCollector) CollectAll(ctx context.Context) (*models.MetricsCollection, error) {
	memory, err := rc.CollectMemoryMetrics(ctx)
	if err != nil {
		return nil, fmt.Errorf("memory metrics collection failed: %v", err)
	}

	conn, err := rc.CollectConnectionMetrics(ctx)
	if err != nil {
		return nil, fmt.Errorf("connection metrics collection failed: %v", err)
	}

	cache, err := rc.CollectCacheMetrics(ctx)
	if err != nil {
		return nil, fmt.Errorf("cache metrics collection failed: %v", err)
	}

	return &models.MetricsCollection{
		Memory:      memory,
		Connection:  conn,
		Cache:       cache,
		CollectedAt: time.Now(),
	}, nil
}
