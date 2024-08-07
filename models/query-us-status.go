package models

type usStatus struct {
	ApplicationID          string `json:"application_id"`
	PassportNumber         string `json:"passport_number"`
	First5LettersOfSurname string `json:"first_5_letters_of_surname"`
}
