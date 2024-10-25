package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"redis-monitor/internal/colloector"
	"redis-monitor/pkg/config"
	"time"
)

func main() {
	// redis connection info
	cfg, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatal(err)
	}

	// initialize collector
	col := colloector.NewRedisCollector(cfg.Server.URL, cfg.Server.Password)

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

		fmt.Printf("\n====== Metrics Collection %d ======\n", i+1)
		fmt.Println(string(jsonData))

		time.Sleep(1 * time.Second)
	}
}
