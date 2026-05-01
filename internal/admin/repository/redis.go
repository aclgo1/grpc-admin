package repository

import (
	"github.com/aclgo/grpc-admin/internal/admin"
	"github.com/redis/go-redis/v9"
)

func NewRedisRepo(rds *redis.Client) admin.RedisRepo {
	return rds
}
