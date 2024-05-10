package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/logs"
	"main/pkg/constants"
	databaseMySQL "main/pkg/databases/mysql"
	database "main/pkg/databases/redis"
	"time"

	"github.com/allegro/bigcache"
)

var ApplicationCache *ApplicationCacheModel

type ApplicationCacheModel struct {
	ClientDetail *bigcache.BigCache
}

func NewBigCache() error {
	bCache, err := bigcache.NewBigCache(bigcache.Config{
		// number of shards (must be a power of 2)
		Shards: 1024,

		// time after which entry can be evicted
		LifeWindow: 10 * time.Minute,

		// Interval between removing expired entries (clean up).
		// If set to <= 0 then no action is performed.
		// Setting to < 1 second is counterproductive — bigcache has a one second resolution.
		CleanWindow: 5 * time.Minute,

		// rps * lifeWindow, used only in initial memory allocation
		MaxEntriesInWindow: 1000 * 10 * 60,

		// max entry size in bytes, used only in initial memory allocation
		MaxEntrySize: 500,

		// prints information about additional memory allocation
		Verbose: false,

		// cache will not allocate more memory than this limit, value in MB
		// if value is reached then the oldest entries can be overridden for the new ones
		// 0 value means no size limit
		HardMaxCacheSize: 256,

		// callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A bitmask representing the reason will be returned.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		OnRemove: nil,

		// OnRemoveWithReason is a callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A constant representing the reason will be passed through.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		// Ignored if OnRemove is specified.
		OnRemoveWithReason: nil,
	})
	if err != nil {
		return fmt.Errorf("new big cache: %w", err)
	}

	var applicationCache ApplicationCacheModel
	applicationCache.ClientDetail = bCache
	ApplicationCache = &applicationCache
	return nil

}

func (applicationCache *ApplicationCacheModel) UpdateCache(key string, clientDetail map[string]string) error {
	// fmt.Println("UPdating Cache.....................")
	bs, err := json.Marshal(&clientDetail)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	return applicationCache.ClientDetail.Set(key, bs)
}

// func userKey(id int64) string {
// 	return strconv.FormatInt(id, 10)
// }

func (applicationCache *ApplicationCacheModel) ReadFromCache(phoneNumber string) (map[string]string, error) {
	bs, err := applicationCache.ClientDetail.Get(phoneNumber)
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			// /*
			// 	first identify the clientId from PhoneNumber
			// 	and   hgetall WA_CLT_DET_PHASE2_
			// */
			logs.InfoLog("Getting Client Details from Redis")
			username := database.CustomHGet(constants.ApplicationConfig.Application.WaPhNoMap, phoneNumber)
			logs.InfoLog("username found %v", username)
			clientId := database.CustomHGet(constants.ApplicationConfig.Application.WaUsrMap, username)
			logs.InfoLog("clientId found %v", clientId)
			clientDetail := make(map[string]string)
			clientDetail = database.CustomHgetAll(constants.ApplicationConfig.Application.MasterCacheKey + clientId + "_" + phoneNumber)
			logs.InfoLog("ClientDetails found from Redis for phoneNumber %v :: %v", phoneNumber, clientDetail)

			if len(clientDetail) == 0 {
				//slack Alert
				logs.InfoLog("Details not found in Redis checking in MySQL for ::%s", phoneNumber)
				clientDetailSql := databaseMySQL.GetClientDetail(phoneNumber)
				if len(clientDetailSql) > 0 {
					clientDetail["wa_url"] = fmt.Sprintf("%v", clientDetailSql["wa_url"])
					clientDetail["wa_token"] = fmt.Sprintf("%v", clientDetailSql["wa_token"])
					clientDetail["username"] = fmt.Sprintf("%v", clientDetailSql["username"])
					clientDetail["phone_number"] = fmt.Sprintf("%v", clientDetailSql["phone_number"])
					logs.InfoLog("Details Found from MySQL for ::%s :: %v", phoneNumber, clientDetail)
				} else {
					logs.InfoLog("No data found in MSQL fallback for phoneNumber::%s", phoneNumber)
				}
			}
			applicationCache.UpdateCache(phoneNumber, clientDetail)
			return clientDetail, nil
		}

		return map[string]string{}, fmt.Errorf("get: %w", err)
	}

	clientDetail := make(map[string]string)
	err = json.Unmarshal(bs, &clientDetail)
	if err != nil {
		return map[string]string{}, fmt.Errorf("unmarshal: %w", err)
	}

	return clientDetail, nil
}

func (applicationCache *ApplicationCacheModel) DeleteCacheEntry(phoneNumber string) {
	applicationCache.ClientDetail.Delete(phoneNumber)
}
