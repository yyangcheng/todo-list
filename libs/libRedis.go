package libs

import (
	"github.com/go-redis/redis"
	"github.com/yyangc/todo-list/config"
)

func InitRedis() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Env.Redis.Host + ":" + config.Env.Redis.Port,
		Password: config.Env.Redis.Password,
		DB:       config.Env.Redis.DB,
	})
	if _, err := rdb.Ping().Result(); err != nil {
		return nil, err
	}
	return rdb, nil
}
