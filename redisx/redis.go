package redisx

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tanenking/svrframe/constants"
	"github.com/tanenking/svrframe/logx"

	"github.com/go-redis/redis"
)

func GetRedis() redis.Cmdable {
	var cmdable redis.Cmdable = nil
	if rdb_cluster != nil {
		cmdable = redis.Cmdable(rdb_cluster)
	} else if rdb_client != nil {
		cmdable = redis.Cmdable(rdb_client)
	}
	return cmdable
}

func Pipelined(fn func(redis.Pipeliner) error) error {
	cmder := GetRedis()
	if cmder == nil {
		logx.ErrorF("no redis instance")
		return nil
	}
	_, err := cmder.Pipelined(fn)
	return err
}

func Exists(key string) bool {
	key = constants.ProjectName + ":" + key
	cmder := GetRedis()
	if cmder == nil {
		logx.ErrorF("no redis instance")
		return false
	}
	i, err := cmder.Exists(key).Result()
	if err != nil {
		logx.ErrorF("Exists %s, err = %v", key, err)
		return false
	}

	if i > 0 {
		return true
	}

	return false
}

func Expire(key string, expiration time.Duration) bool {
	key = constants.ProjectName + ":" + key
	cmder := GetRedis()
	if cmder == nil {
		logx.ErrorF("no redis instance")
		return false
	}
	b, err := cmder.Expire(key, expiration).Result()
	if err != nil {
		logx.ErrorF("Expire %s, err = %v", key, err)
		return false
	}
	return b
}

func Del(keys ...string) {
	cmder := GetRedis()
	if cmder == nil {
		logx.ErrorF("no redis instance")
		return
	}
	real_keys := []string{}
	for _, s := range keys {
		real_keys = append(real_keys, constants.ProjectName+s)
	}
	cmder.Del(real_keys...)
}

func SetNX(key string, val interface{}) bool {
	key = constants.ProjectName + ":" + key
	cmder := GetRedis()
	if cmder == nil {
		err := fmt.Errorf("no redis instance")
		logx.ErrorF("%v", err)
		return false
	}
	status, err := cmder.SetNX(key, val, time.Second*-1).Result()
	if err != nil {
		logx.ErrorF("Set %s, err = %v", key, err)
		return false
	}
	return status
}

func SetAndExp(key string, val interface{}, expiration time.Duration) (status string, err error) {
	key = constants.ProjectName + ":" + key
	cmder := GetRedis()
	if cmder == nil {
		err = fmt.Errorf("no redis instance")
		logx.ErrorF("%v", err)
		return
	}
	status, err = cmder.Set(key, val, expiration).Result()
	if err != nil {
		logx.ErrorF("SetAndExp %s, err = %v", key, err)
	}
	return
}
func Set(key string, val interface{}) (status string, err error) {
	key = constants.ProjectName + ":" + key
	cmder := GetRedis()
	if cmder == nil {
		err = fmt.Errorf("no redis instance")
		logx.ErrorF("%v", err)
		return
	}
	status, err = cmder.Set(key, val, time.Second*-1).Result()
	if err != nil {
		logx.ErrorF("Set %s, err = %v", key, err)
	}
	return
}
func Get(key string) (val string, err error) {
	key = constants.ProjectName + ":" + key
	cmder := GetRedis()
	if cmder == nil {
		err = fmt.Errorf("no redis instance")
		logx.ErrorF("%v", err)
		return
	}
	val, err = cmder.Get(key).Result()
	if err != nil {
		logx.ErrorF("Get %s, err = %v", key, err)
	}
	return
}
func SetStruct(key string, stc interface{}) (status string, err error) {
	j, err := json.Marshal(stc)
	if err != nil {
		logx.ErrorF("%v", err)
		return
	}
	return Set(key, j)
}
func HSet(key, field string, value interface{}) (status bool, err error) {
	key = constants.ProjectName + ":" + key
	cmder := GetRedis()
	if cmder == nil {
		err = fmt.Errorf("no redis instance")
		logx.ErrorF("%v", err)
		return
	}
	status, err = cmder.HSet(key, field, value).Result()
	if err != nil {
		logx.ErrorF("HSet %s, field %s, err = %v", key, field, err)
	}
	return
}
func HGet(key, field string) (val string, err error) {
	key = constants.ProjectName + ":" + key
	cmder := GetRedis()
	if cmder == nil {
		err = fmt.Errorf("no redis instance")
		logx.ErrorF("%v", err)
		return
	}
	val, err = cmder.HGet(key, field).Result()
	if err != nil {
		logx.ErrorF("HGet %s, field %s, err = %v", key, field, err)
	}
	return
}
func HDel(key string, field ...string) (err error) {
	key = constants.ProjectName + ":" + key
	cmder := GetRedis()
	if cmder == nil {
		err = fmt.Errorf("no redis instance")
		logx.ErrorF("%v", err)
		return
	}
	_, err = cmder.HDel(key, field...).Result()
	if err != nil {
		logx.ErrorF("HDel %s, field %v, err = %v", key, field, err)
	}
	return
}

// -1代表出错
func HLen(key string) int64 {
	key = constants.ProjectName + ":" + key
	cmder := GetRedis()
	if cmder == nil {
		err := fmt.Errorf("no redis instance")
		logx.ErrorF("%v", err)
		return -1
	}
	len, err := cmder.HLen(key).Result()
	if err != nil {
		logx.ErrorF("HLen %s, err = %v", key, err)
		return -1
	}
	return len
}

func HMSet(key string, fields map[string]interface{}) (status string, err error) {
	key = constants.ProjectName + ":" + key
	cmder := GetRedis()
	if cmder == nil {
		err = fmt.Errorf("no redis instance")
		logx.ErrorF("%v", err)
		return
	}
	status, err = cmder.HMSet(key, fields).Result()
	if err != nil {
		logx.ErrorF("HMSet %s, err = %v", key, err)
	}
	return
}
func HMGet(key string, fields []string) (result map[string]interface{}, err error) {
	key = constants.ProjectName + ":" + key

	result = map[string]interface{}{}

	cmder := GetRedis()
	if cmder == nil {
		err = fmt.Errorf("no redis instance")
		logx.ErrorF("%v", err)
		return
	}
	if len(fields) <= 0 {
		return
	}
	var res []interface{}
	res, err = cmder.HMGet(key, fields...).Result()
	if err != nil {
		logx.ErrorF("HMGet %s, err = %v", key, err)
		return
	}
	if err == nil && len(res) > 0 {
		for idx, field := range fields {
			if res[idx] != nil {
				result[field] = res[idx]
			}
		}
	}
	return
}

func HGetAll(key string) (result map[string]string, err error) {
	key = constants.ProjectName + ":" + key
	cmder := GetRedis()
	if cmder == nil {
		err = fmt.Errorf("no redis instance")
		logx.ErrorF("%v", err)
		return
	}
	result, err = cmder.HGetAll(key).Result()
	if err != nil {
		logx.ErrorF("HGetAll %s, err = %v", key, err)
	}
	return
}
