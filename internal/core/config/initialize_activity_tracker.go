package config

import (
	"fmt"
	"log/slog"

	goroutines "github.com/giulio-alfieri/toq_server/internal/core/go_routines"
	"github.com/redis/go-redis/v9"
)

func (c *config) InitializeActivityTracker() error {
	// Initialize Activity Tracker early (before gRPC middleware needs it)
	redisOpts, err := redis.ParseURL(c.env.REDIS.URL)
	if err != nil {
		slog.Error("failed to parse Redis URL for activity tracker", "error", err)
		return fmt.Errorf("failed to parse Redis URL: %w", err)
	}
	redisClient := redis.NewClient(redisOpts)

	// Test Redis connection
	_, err = redisClient.Ping(c.context).Result()
	if err != nil {
		slog.Error("failed to connect to Redis for activity tracker", "error", err)
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	// Create activity tracker with a placeholder user service (will be set later)
	c.activityTracker = goroutines.NewActivityTracker(redisClient, nil)

	slog.Info("Activity tracker initialized successfully")
	return nil
}

func (c *config) SetActivityTrackerUserService() {
	// Set the user service after it's been initialized
	if c.activityTracker != nil {
		c.activityTracker.SetUserService(c.userService)
	}
}
