package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"redis-monitor/internal/colloector"
	"time"
)

func main() {
	// redis connection info
	redisAddr := "redis-19186.c278.us-east-1-4.ec2.redns.redis-cloud.com:19186"
	redisPassword := "SF5vF4Rbz5nEmSQJb9W2fwUwRuFJWlUS"

	// initialize collector
	col := colloector.NewRedisCollector(redisAddr, redisPassword)

	// set up context
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		metrics, err := col.CollectAll(ctx)
		if err != nil {
			log.Fatal("Failed to collect metrics: %w", err)
		}
		// marshal json for pretty print
		jsonData, err := json.MarshalIndent(metrics, "", "     ")
		if err != nil {
			log.Fatal("Failed to marshal metrics: %w", err)
		}

		fmt.Printf("\n=== Metrics Collection %d ===\n", i+1)
		fmt.Println(string(jsonData))

		time.Sleep(1 * time.Second)
	}
}
