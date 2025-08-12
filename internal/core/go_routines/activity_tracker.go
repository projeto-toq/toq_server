package goroutines

import (
	"context"
	"log/slog"
	"strconv"
	"sync"
	"time"

	userservices "github.com/giulio-alfieri/toq_server/internal/core/service/user_service"
	"github.com/redis/go-redis/v9"
)

// ActivityTracker manages user activity tracking with Redis + batch updates
type ActivityTracker struct {
	redisClient   *redis.Client
	userService   userservices.UserServiceInterface
	batchSize     int
	flushInterval time.Duration
}

func NewActivityTracker(redisClient *redis.Client, userService userservices.UserServiceInterface) *ActivityTracker {
	return &ActivityTracker{
		redisClient:   redisClient,
		userService:   userService,
		batchSize:     100,              // Process 100 users per batch
		flushInterval: 30 * time.Second, // Flush every 30 seconds
	}
}

// SetUserService sets the user service after initialization
func (at *ActivityTracker) SetUserService(userService userservices.UserServiceInterface) {
	at.userService = userService
}

// TrackActivity records user activity in Redis (very fast, non-blocking)
func (at *ActivityTracker) TrackActivity(ctx context.Context, userID int64) {
	key := "user_activity:" + strconv.FormatInt(userID, 10)

	// Set with TTL - will expire if user becomes inactive
	err := at.redisClient.Set(ctx, key, time.Now().Unix(), 5*time.Minute).Err()
	if err != nil {
		slog.Error("Failed to track activity in Redis", "userID", userID, "error", err)
		// Fallback: could add to channel for immediate DB update
	}
}

// StartBatchWorker starts the worker that periodically flushes Redis data to MySQL
func (at *ActivityTracker) StartBatchWorker(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	ticker := time.NewTicker(at.flushInterval)
	defer ticker.Stop()

	slog.Info("Activity batch worker started", "interval", at.flushInterval)

	for {
		select {
		case <-ctx.Done():
			slog.Info("Activity batch worker stopped")
			return
		case <-ticker.C:
			at.flushActivitiesToDB(ctx)
		}
	}
}

// flushActivitiesToDB processes Redis data and updates MySQL in batches
func (at *ActivityTracker) flushActivitiesToDB(ctx context.Context) {
	// Get all active users from Redis
	keys, err := at.redisClient.Keys(ctx, "user_activity:*").Result()
	if err != nil {
		slog.Error("Failed to get activity keys from Redis", "error", err)
		return
	}

	if len(keys) == 0 {
		slog.Debug("No active users to process")
		return
	}

	slog.Debug("Processing active users", "count", len(keys))

	// Process in batches
	for i := 0; i < len(keys); i += at.batchSize {
		end := i + at.batchSize
		if end > len(keys) {
			end = len(keys)
		}

		batch := keys[i:end]
		at.processBatch(ctx, batch)
	}
}

// processBatch handles a batch of user activities
func (at *ActivityTracker) processBatch(ctx context.Context, keys []string) {
	userIDs := make([]int64, 0, len(keys))
	timestamps := make([]int64, 0, len(keys))

	// Get timestamps from Redis
	pipe := at.redisClient.Pipeline()
	for _, key := range keys {
		pipe.Get(ctx, key)
	}

	results, err := pipe.Exec(ctx)
	if err != nil {
		slog.Error("Failed to get batch data from Redis", "error", err)
		return
	}

	// Parse results
	for i, result := range results {
		if result.Err() != nil {
			continue // Skip failed keys
		}

		timestampStr := result.(*redis.StringCmd).Val()
		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			continue
		}

		// Extract userID from key: "user_activity:123" -> 123
		key := keys[i]
		userIDStr := key[len("user_activity:"):]
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			continue
		}

		userIDs = append(userIDs, userID)
		timestamps = append(timestamps, timestamp)
	}

	if len(userIDs) == 0 {
		return
	}

	// Check if userService is available
	if at.userService == nil {
		slog.Warn("User service not available, skipping batch update", "count", len(userIDs))
		return
	}

	// Batch update MySQL
	err = at.userService.BatchUpdateLastActivity(ctx, userIDs, timestamps)
	if err != nil {
		slog.Error("Failed to batch update activities", "count", len(userIDs), "error", err)
		return
	}

	slog.Debug("Successfully updated user activities", "count", len(userIDs))
}

// GetActiveUsers returns list of currently active users (from Redis)
func (at *ActivityTracker) GetActiveUsers(ctx context.Context) ([]int64, error) {
	keys, err := at.redisClient.Keys(ctx, "user_activity:*").Result()
	if err != nil {
		return nil, err
	}

	userIDs := make([]int64, 0, len(keys))
	for _, key := range keys {
		userIDStr := key[len("user_activity:"):]
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			continue
		}
		userIDs = append(userIDs, userID)
	}

	return userIDs, nil
}

// GetActiveUserCount returns count of active users (very fast)
func (at *ActivityTracker) GetActiveUserCount(ctx context.Context) (int64, error) {
	keys, err := at.redisClient.Keys(ctx, "user_activity:*").Result()
	if err != nil {
		return 0, err
	}
	return int64(len(keys)), nil
}
