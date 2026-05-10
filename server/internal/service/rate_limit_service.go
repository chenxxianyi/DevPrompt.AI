package service

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimitService struct {
	redis  *redis.Client
	enabled bool
}

func NewRateLimitService(rdb *redis.Client, enabled bool) *RateLimitService {
	return &RateLimitService{redis: rdb, enabled: enabled}
}

const (
	dailyGenerateKey = "daily_generate:%d:%s"
	rateLimitKey     = "rate_limit:%d:%s"
)

// CheckDailyLimit 检查用户当日生成次数是否超限
// 返回 (当前次数, 限制, 是否允许)
func (s *RateLimitService) CheckDailyLimit(userID uint64, dailyLimit int) (int, int, bool, error) {
	if !s.enabled {
		return 0, dailyLimit, true, nil
	}

	ctx := context.Background()
	today := time.Now().Format("2006-01-02")
	key := fmt.Sprintf(dailyGenerateKey, userID, today)

	count, err := s.redis.Get(ctx, key).Int()
	if err != nil && err != redis.Nil {
		return 0, dailyLimit, false, err
	}

	if count >= dailyLimit {
		return count, dailyLimit, false, nil
	}

	return count, dailyLimit, true, nil
}

// IncrementDailyCount 增加用户当日生成次数
func (s *RateLimitService) IncrementDailyCount(userID uint64) error {
	if !s.enabled {
		return nil
	}

	ctx := context.Background()
	today := time.Now().Format("2006-01-02")
	key := fmt.Sprintf(dailyGenerateKey, userID, today)

	pipe := s.redis.Pipeline()
	pipe.Incr(ctx, key)
	// TTL 到次日凌晨
	now := time.Now()
	nextDay := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	pipe.ExpireAt(ctx, key, nextDay)
	_, err := pipe.Exec(ctx)
	return err
}

// CheckRateLimit 接口级别限流 (按用户+路由)
func (s *RateLimitService) CheckRateLimit(userID uint64, route string, rpm int) (bool, error) {
	if !s.enabled {
		return true, nil
	}

	ctx := context.Background()
	key := fmt.Sprintf(rateLimitKey, userID, route)

	pipe := s.redis.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, 60*time.Second)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return true, err
	}

	return incr.Val() <= int64(rpm), nil
}

// GetDailyCount 获取用户当日生成次数
func (s *RateLimitService) GetDailyCount(userID uint64) int {
	if !s.enabled {
		return 0
	}

	ctx := context.Background()
	today := time.Now().Format("2006-01-02")
	key := fmt.Sprintf(dailyGenerateKey, userID, today)

	count, _ := s.redis.Get(ctx, key).Int()
	return count
}
