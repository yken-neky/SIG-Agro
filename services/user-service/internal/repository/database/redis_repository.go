package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sig-agro/services/user-service/internal/entity"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{client: client}
}

func (r *RedisRepository) GetUser(id int64) (bool, *entity.User) {
	key := fmt.Sprintf("user:%d", id)
	val, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		return false, nil
	}

	var user entity.User
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return false, nil
	}

	return true, &user
}

func (r *RedisRepository) SetUser(id int64, u *entity.User) {
	key := fmt.Sprintf("user:%d", id)
	data, _ := json.Marshal(u)
	r.client.Set(context.Background(), key, data, time.Hour)
}
