package bootstrap

import (
	"context"

	"smart.gitlab.biomind.com.cn/infrastructure/file-server-gateway/model"

	"smart.gitlab.biomind.com.cn/infrastructure/biogo/config"
	"smart.gitlab.biomind.com.cn/infrastructure/biogo/config/center"
	"smart.gitlab.biomind.com.cn/infrastructure/biogo/redis"
)

var redisNames = []string{
	model.RedisEngine,
}

func InitRedis(ctx context.Context) error {

	configCenterClient := config.GetClient()
	for _, redisName := range redisNames {
		conf, err := configCenterClient.Get(ctx, redisName)
		if err != nil {
			panic(err)
		}
		err = redis.Init(redisName, conf.GetValue())
		if err != nil {
			panic(err)
		}
		// 监听修改
		configCenterClient.Watch(ctx, conf, func(item *center.Item) {
			err = redis.Init(conf.Key, item.GetValue())
			if err != nil {
				panic(err)
			}
		})
	}

	return nil
}
