package redis

import (
	"errors"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	goredis "github.com/go-redis/redis"
	"github.com/gomodule/redigo/redis"
)

var (
	// defaultRedisTomlConfigPath 默认redis配置
	defaultRedisTomlConfigPath = "./config/redis.toml"
)

func getRedisOptions(index string) (options goredis.Options, err error) {
	if _, err := os.Stat(defaultRedisTomlConfigPath); os.IsNotExist(err) {
		defaultRedisTomlConfigPath = "../config/redis.toml"
	}

	var config map[string]interface{}
	_, err = toml.DecodeFile(defaultRedisTomlConfigPath, &config)
	if err != nil {
		return
	}
	var env string
	instanceSlice := strings.Split(index, ".")
	if len(instanceSlice) == 2 {
		index = instanceSlice[0]
		env = instanceSlice[1]
	} else {
		env = config["env"].(string)
	}
	if instanceConf, ok := config[index].(map[string]interface{}); ok {
		if envConf, ok := instanceConf[env].(map[string]interface{}); ok {
			options = goredis.Options{
				Host:           envConf["host"].(string),
				Port:           envConf["port"].(string),
				Password:       envConf["password"].(string),
				MaxIdle:        int(envConf["maxIdle"].(int64)),
				MaxOpen:        int(envConf["maxOpen"].(int64)),
				ConnectTimeout: int(envConf["connect_timeout"].(int64)),
				ReadTimeout:    int(envConf["read_timeout"].(int64)),
				WriteTimeout:   int(envConf["write_timeout"].(int64)),
				IdleTimeout:    int(envConf["idle_timeout"].(int64)),
			}
			return
		}
	}
	err = errors.New("redis options read failed, " + index)
	return
}

// NewRedis 获取redis
func NewRedis(instance string) *redis.Pool {
	if poolLoad, ok := goredis.Pool.Load(instance); ok {
		redisPool := poolLoad.(*redis.Pool)
		return redisPool
	}

	options, err := getRedisOptions(instance)
	if err != nil {
		panic(err)
	}

	if err := goredis.NewClient(instance, options); err != nil {
		panic(err)
	}

	poolLoad, _ := goredis.Pool.Load(instance)
	redisPool := poolLoad.(*redis.Pool)
	return redisPool
}
