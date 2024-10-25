package models

import "time"

// MetricType represents the type of Redis metric
type MetricType string

const (
	MemoryMetric     MetricType = "memory"
	ConnectionMetric MetricType = "connection"
	CacheMetric      MetricType = "cache"
)

// BaseMetric contains common fields for all metrics
type BaseMetric struct {
	Timestamp time.Time  `json:"timestamp"`
	Type      MetricType `json:"type"`
	Instance  string     `json:"instance"`
}

// MemoryMetrics represents Redis memory usage statistics
type MemoryMetrics struct {
	BaseMetric
	UsedMemoryBytes int64   `json:"used_memory_bytes"`
	MaxMemoryBytes  int64   `json:"max_memory_bytes"`
	MemoryFragRatio float64 `json:"memory_frag_ratio"`
	MemoryUsagePerc float64 `json:"memory_usage_percentage"`
}

// ConnectionMetrics represents Redis connection statistics
type ConnectionMetrics struct {
	BaseMetric
	ConnectedClients int64 `json:"connected_clients"`
	BlockedClients   int64 `json:"blocked_clients"`
	ConnectedSlaves  int64 `json:"connected_slaves"`
	MaxClients       int64 `json:"max_clients"`
}

// CacheMetrics represents Redis cache performance statistics
type CacheMetrics struct {
	BaseMetric
	KeyspaceHits   int64   `json:"keyspace_hits"`
	KeyspaceMisses int64   `json:"keyspace_misses"`
	HitRate        float64 `json:"hit_rate"`
	Evictions      int64   `json:"evictions"`
	ExpiredKeys    int64   `json:"expired_keys"`
}

// MetricsCollection represents a collection of all metrics at a point in time
type MetricsCollection struct {
	Memory      *MemoryMetrics     `json:"memory,omitempty"`
	Connection  *ConnectionMetrics `json:"connection,omitempty"`
	Cache       *CacheMetrics      `json:"cache,omitempty"`
	CollectedAt time.Time          `json:"collected_at"`
}
