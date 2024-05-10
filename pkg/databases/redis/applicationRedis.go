package database

import (
	"context"
	"fmt"
	"main/logs"
	"main/pkg/constants"
	"time"

	// "wh_dequeuer/pkg/constants"

	"github.com/redis/go-redis/v9"
)

var RedisQueueClient *redis.Client
var RedisCacheClient *redis.Client

var Ctx = context.Background()

func EstablishRedisQueueConnecion() {
	logs.InfoLog("Establishing Redis Queue Connection")
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v", constants.ApplicationConfig.RedisQueue.Host+":"+constants.ApplicationConfig.RedisQueue.Port),
		Username: constants.ApplicationConfig.RedisQueue.Username,
		Password: constants.ApplicationConfig.RedisQueue.Password, // no password set
		DB:       0,                                               // use default DB
	})
	RedisQueueClient = rdb
	logs.InfoLog("Redis Queue Connection Estblished....")
}

func EstablishRedisCacheConnecion() {
	logs.InfoLog("Establishing Redis Cache Connection")
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v", constants.ApplicationConfig.RedisCache.Host+":"+constants.ApplicationConfig.RedisCache.Port),
		Username: constants.ApplicationConfig.RedisCache.Username,
		Password: constants.ApplicationConfig.RedisCache.Password, // no password set
		DB:       0,                                               // use default DB

	})
	RedisCacheClient = rdb
	logs.InfoLog("Redis Cache Estblished....")
}

func CustomBLpop(listName string) ([]string, error) {
	// logs.InfoLog("In custom pop function")
	if RedisQueueClient == nil {
		EstablishRedisQueueConnecion()
	}
	poppedString, err := RedisQueueClient.BLPop(Ctx, time.Second*1, listName).Result()
	// logs.InfoLog("poppedString::%v", poppedString)
	// logs.ErrorLog("poppedString::%v", poppedString)
	if err != nil {
		// fmt.Println("err", err)
		// fmt.Println("errString", err.Error())
		return nil, err
	}

	return poppedString, nil
}

func CustomRpush(listName string, stringParam string) {
	logs.InfoLog("Performing rpush, QUEUE[%v], String[%v]", listName, stringParam)
	if RedisQueueClient == nil {
		EstablishRedisQueueConnecion()
	}
	_, err := RedisQueueClient.RPush(Ctx, listName, stringParam).Result()
	if err != nil {
		logs.ErrorLog("Error to perform rpush, QUEUE[%v], String[%v]", listName, stringParam)
		logs.ErrorLog("Error in CustomRpush function, %v", err)
	}
}

func CustomSetKey(key string, value string, expiration time.Duration) {
	if RedisCacheClient == nil {
		EstablishRedisCacheConnecion()
	}
	err := RedisCacheClient.Set(Ctx, key, value, expiration).Err()
	if err != nil {
		logs.ErrorLog("Error in CustomSetKey function, %v", err)

	}

}

func CustomHgetAll(hashKey string) map[string]string {
	logs.InfoLog("HashKey is :: %v", hashKey)
	if RedisCacheClient == nil {
		EstablishRedisCacheConnecion()
	}
	hashMap, err := RedisCacheClient.HGetAll(Ctx, hashKey).Result()
	if err != nil {
		logs.ErrorLog("Error in CustomHGetAll function, %v", err)
		return nil
	}
	logs.InfoLog("Map for client formed :: %v", hashMap)
	return hashMap
}

func CustomHGet(hashKey string, key string) string {
	logs.InfoLog("CustomHGet :: HashKey::%v :: key:: %v", hashKey, key)
	if RedisCacheClient == nil {
		EstablishRedisCacheConnecion()
	}
	value, err := RedisCacheClient.HGet(Ctx, hashKey, key).Result()

	if err != nil {
		logs.ErrorLog("Error in CustomHGet function, %v", err)
		return ""
	}
	return value
}

func CustomHSet(hashKey string, key string, field string) error {
	logs.InfoLog("CustomHSet :: HashKey::%v :: key:: %v :: field :: %v", hashKey, key, field)
	if RedisCacheClient == nil {
		EstablishRedisCacheConnecion()
	}
	_, err := RedisCacheClient.HSet(Ctx, hashKey, key, field).Result()

	if err != nil {
		logs.ErrorLog("Error in CustomHSet function, %v", err)
		return err
	}
	return nil
}
