package controller

import (
	"context"
	"crawler-visa/config"
	"crawler-visa/models"
	"crawler-visa/utils"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"net/http"
)

var ctx = context.Background()
var redisClient = config.ConfigureRedis()

const keyPrefix = "application:status:"

func CreateApplication(w http.ResponseWriter, r *http.Request) {
	queryUsStatus := &models.QueryUsStatus{}
	utils.ParseBody(r, queryUsStatus)
	marshal, err := json.Marshal(queryUsStatus)
	if err != nil {
		utils.ResultError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = redisClient.Set(ctx, keyPrefix+queryUsStatus.ApplicationID, marshal, 0).Err()
	if err != nil {
		utils.ResultError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.ResultJSON(w, queryUsStatus, "保存成功")
}

// RetrieveApplication 通过Application ID从Redis获取签证申请记录
func RetrieveApplication(w http.ResponseWriter, r *http.Request) {
	appID := r.URL.Query().Get("application_id")
	if appID == "" {
		utils.ResultError(w, "Application ID is required", http.StatusInternalServerError)
		return
	}
	result, err := redisClient.Get(ctx, keyPrefix+appID).Result()
	if errors.Is(err, redis.Nil) {
		utils.ResultError(w, "Application not found", http.StatusInternalServerError)
		return
	} else if err != nil {
		utils.ResultError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.ResultJSON(w, result, "检索成功")
}
