package cache

import (
	"ChainClientManageSystem/config"
	"ChainClientManageSystem/internal/model"
	"ChainClientManageSystem/pkg/constant"
	"ChainClientManageSystem/utils"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"time"
)

func GetUserInfoFromCache(username string) (*model.User, error) {
	redisKey := constant.UserInfoPrefix + username
	val, err := utils.GetRedisCli().Get(context.Background(), redisKey).Result()
	if err != nil {
		return nil, err
	}
	user := &model.User{}
	err = json.Unmarshal([]byte(val), user)
	return user, err
}

func SetUserCacheInfo(user *model.User) error {
	redisKey := constant.UserInfoPrefix + user.Name
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}
	expired := time.Second * time.Duration(config.GetGlobalConf().Cache.UserExpired)
	_, err = utils.GetRedisCli().Set(context.Background(), redisKey, val, expired*time.Second).Result()
	return err
}

func UpdateCachedUserInfo(user *model.User) error {
	err := SetUserCacheInfo(user)
	if err != nil {
		redisKey := constant.UserInfoPrefix + user.Name
		utils.GetRedisCli().Del(context.Background(), redisKey).Result()
	}
	log.Info("update cached user info success")
	return err
}

func GetSessionInfo(session string) (*model.User, error) {
	redisKey := constant.SessionKeyPrefix + session
	val, err := utils.GetRedisCli().Get(context.Background(), redisKey).Result()
	if err != nil {
		return nil, err
	}
	user := &model.User{}
	err = json.Unmarshal([]byte(val), &user)
	return user, err
}

func SetSessionInfo(user *model.User, session string) error {
	redisKey := constant.SessionKeyPrefix + session
	val, err := json.Marshal(&user)
	if err != nil {
		return err
	}
	expired := time.Second * time.Duration(config.GetGlobalConf().Cache.SessionExpired)
	_, err = utils.GetRedisCli().Set(context.Background(), redisKey, val, expired*time.Second).Result()
	return err
}

func DelSessionInfo(session string) error {
	redisKey := constant.SessionKeyPrefix + session
	_, err := utils.GetRedisCli().Del(context.Background(), redisKey).Result()
	return err
}

func DelUserCacheInfo(username string) error {
	redisKey := constant.UserInfoPrefix + username
	_, err := utils.GetRedisCli().Del(context.Background(), redisKey).Result()
	return err
}
