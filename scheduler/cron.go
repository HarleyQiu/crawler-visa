package scheduler

import (
	"crawler-visa/models"
	"crawler-visa/service"
	"fmt"
	"time"
)

func Corn() {
	var queries []models.QueryUsStatus

	queries = append(queries, models.QueryUsStatus{
		Location:               "BEJ",
		ApplicationID:          "AA00DJIS17",
		PassportNumber:         "R1708298",
		First5LettersOfSurname: "Zheng",
	})

	queries = append(queries, models.QueryUsStatus{
		Location:               "BEJ",
		ApplicationID:          "AA00DJIGNH",
		PassportNumber:         "EJ5341909",
		First5LettersOfSurname: "Zhang",
	})

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {

			for _, query := range queries {

				fmt.Printf("%+v\n", query)
				applicationCheck, err := service.RunVisaStatusCheck(&query)
				applicationCheck.Code = 200
				if err != nil {
					return
				}
				fmt.Printf("%+v\n", applicationCheck)
			}
		}
	}()
}
