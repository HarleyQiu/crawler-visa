package controller

import (
	"context"
	"crawler-visa/config"
	"crawler-visa/models"
	"crawler-visa/utils"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
)

// 全局上下文和Redis客户端初始化
var ctx = context.Background()
var redisClient = config.ConfigureRedis()

// Redis中用于存储应用状态的键的前缀
const keyPrefix = "application:status:"

// CreateApplication 根据提供的请求数据在Redis中创建应用状态记录。
// 它将请求正文解析为QueryUsStatus模型，将其编组为JSON，
// 并使用由应用程序ID构造的密钥将其存储在Redis中。
// 如果成功，它会响应一条指示成功的JSON消息；否则，它返回错误响应。
func CreateApplication(w http.ResponseWriter, r *http.Request) {
	log.Println("开始创建应用状态") // 添加日志
	queryUsStatus := &models.QueryUsStatus{}
	utils.ParseBody(r, queryUsStatus)
	marshal, err := json.Marshal(queryUsStatus)
	if err != nil {
		utils.ResultError(w, err.Error(), http.StatusInternalServerError)
		log.Println("JSON编码失败:", err) // 添加错误日志
		return
	}
	err = redisClient.Set(ctx, keyPrefix+queryUsStatus.ApplicationID, marshal, 0).Err()
	if err != nil {
		utils.ResultError(w, err.Error(), http.StatusInternalServerError)
		log.Println("Redis写入失败:", err) // 添加错误日志
		return
	}
	utils.ResultJSON(w, queryUsStatus, "保存成功")
	log.Println("应用状态创建成功") // 添加成功日志
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
	log.Printf("检索应用状态成功：%s\n", appID) // 添加日志
}

// UpdateApplication 更新Redis中的应用状态记录
func UpdateApplication(w http.ResponseWriter, r *http.Request) {
	queryUsStatus := &models.QueryUsStatus{}
	utils.ParseBody(r, queryUsStatus)
	exists, err := redisClient.Exists(ctx, keyPrefix+queryUsStatus.ApplicationID).Result()
	if err != nil {
		utils.ResultError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exists == 0 {
		utils.ResultError(w, "键不存在，无法更新", http.StatusBadRequest)
		return
	}
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
	utils.ResultJSON(w, nil, "修改成功")
	log.Printf("更新应用状态成功：%s\n", queryUsStatus.ApplicationID) // 添加日志
}

// DeleteApplication 通过Application ID删除Redis中的签证申请记录
func DeleteApplication(w http.ResponseWriter, r *http.Request) {
	appID := r.URL.Query().Get("application_id")
	if appID == "" {
		utils.ResultError(w, "Application ID is required", http.StatusInternalServerError)
		return
	}
	result, err := redisClient.Del(ctx, keyPrefix+appID).Result()
	if err != nil {
		utils.ResultError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if result == 0 {
		utils.ResultError(w, "Application not found", http.StatusNotFound)
		return
	}
	utils.ResultJSON(w, nil, "删除成功")
	log.Printf("删除应用状态成功：%s\n", appID) // 添加日志
}

// RetrieveAllApplications 从Redis获取所有签证申请记录
func RetrieveAllApplications(w http.ResponseWriter, r *http.Request) {
	var applications []models.QueryUsStatus

	iter := redisClient.Scan(ctx, 0, keyPrefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		result, err := redisClient.Get(ctx, iter.Val()).Result()
		if err != nil {
			utils.ResultError(w, "Error retrieving application: "+err.Error(), http.StatusInternalServerError)
			return
		}
		application := &models.QueryUsStatus{}
		err = json.Unmarshal([]byte(result), application)
		if err != nil {
			utils.ResultError(w, "Error parsing application data: "+err.Error(), http.StatusInternalServerError)
			return
		}
		applications = append(applications, *application)
	}
	if err := iter.Err(); err != nil {
		utils.ResultError(w, "Error iterating keys: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if len(applications) == 0 {
		utils.ResultJSON(w, nil, "No applications found")
		log.Println("未找到任何应用状态") // 添加日志
	} else {
		utils.ResultJSON(w, applications, "Applications retrieved successfully")
		log.Println("成功检索所有应用状态") // 添加日志
	}
}
