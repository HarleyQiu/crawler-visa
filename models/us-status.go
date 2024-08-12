package models

type QueryUsStatus struct {
	Location               string `json:"location"`
	ApplicationID          string `json:"application_id"`
	PassportNumber         string `json:"passport_number"`
	First5LettersOfSurname string `json:"first_5_letters_of_surname"`
}

type UsStatus struct {
	Status      string `json:"status"`
	Created     string `json:"created"`
	LastUpdated string `json:"last_updated"`
	Code        int    `json:"code"`
}
