package scheduler

import (
	"crawler-visa/models"
	"crawler-visa/service"
	"crawler-visa/utils"
	"fmt"
	"time"
)

func RunScheduledTasks() {
	queryLoader := utils.NewQueryLoader[models.QueryUsStatus]("configuration.json")
	queries, err := queryLoader.LoadQueries()
	if err != nil {
		fmt.Printf("加载查询错误: %v\n", err)
		return
	}

	tracker := utils.NewStatusTracker[models.UsStatus]()

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {

			for _, query := range queries {

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
				}
			}
		}
	}()
}
