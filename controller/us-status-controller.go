package controller

import (
	"context"
	"crawler-visa/config"
	"crawler-visa/models"
	"crawler-visa/utils"
	"encoding/json"
	"errors"
	"fmt"
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
	fmt.Println(queryUsStatus)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = redisClient.Set(ctx, keyPrefix+queryUsStatus.ApplicationID, marshal, 0).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// RetrieveApplication 通过Application ID从Redis获取签证申请记录
func RetrieveApplication(w http.ResponseWriter, r *http.Request) {
	appID := r.URL.Query().Get("application_id")
	if appID == "" {
		http.Error(w, "Application ID is required", http.StatusBadRequest)
		return
	}
	result, err := redisClient.Get(ctx, keyPrefix+appID).Result()
	if errors.Is(err, redis.Nil) {
		http.Error(w, "Application not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}
