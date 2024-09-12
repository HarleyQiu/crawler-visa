package scheduler

import (
	"context"
	"crawler-visa/config"
	"crawler-visa/models"
	"crawler-visa/service"
	"crawler-visa/utils"
	"encoding/json"
	"fmt"
	"time"
)

var ctx = context.Background()
var redisClient = config.ConfigureRedis()

const keyPattern = "application:status:*"

func RunScheduledTasks() {

	tracker := utils.NewStatusTracker[models.UsStatus]()
	sender := utils.NewNotificationSender("https://apis.visa5i.com/wuai/system/wechat-notification/save")

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			iter := redisClient.Scan(ctx, 0, keyPattern, 0).Iterator()

			for iter.Next(ctx) {
				result, err := redisClient.Get(ctx, iter.Val()).Result()
				if err != nil {
					fmt.Printf("从Redis读取查询错误: %v\n", err)
					continue
				}

				var query models.QueryUsStatus
				err = json.Unmarshal([]byte(result), &query)
				if err != nil {
					fmt.Printf("解析查询数据错误: %v\n", err)
					continue
				}

				fmt.Printf("查询信息：%+v\n", query)
				usStatus, err := service.RunVisaStatusCheck(&query)
				usStatus.Code = 200
				if err != nil {
					fmt.Printf("检查签证状态错误: %v\n", err)
					continue
				}
				changed := tracker.UpdateStatus(query.ApplicationID, usStatus)
				if changed {
					fmt.Printf("状态变更：%s, 新状态：%+v\n", query.ApplicationID, usStatus)
					remark := utils.FormatVisaStatus(usStatus.Status, usStatus.StatusContent, usStatus.Created, usStatus.LastUpdated, query.ApplicationID, query.PassportNumber)

					notificationData := utils.NotificationData{
						Sys:        query.Location,
						ConsDist:   "美签预约状态查询",
						MonCountry: "美签预约状态查询",
						ApptTime:   usStatus.LastUpdated,
						Status:     "2",
						UserName:   query.ApplicationID,
						Remark:     remark,
					}
					err := sender.SendNotification(notificationData)
					if err != nil {
						fmt.Printf("Error sending notification: %v\n", err)
					}
				}
			}
		}
	}()
}
