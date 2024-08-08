package controller

import (
	"crawler-visa/models"
	"crawler-visa/service"
	"crawler-visa/utils"
	"encoding/json"
	"net/http"
)

var status models.QueryUsStatus

func StatusCheck(w http.ResponseWriter, r *http.Request) {
	queryUsStatus := &models.QueryUsStatus{}
	utils.ParseBody(r, queryUsStatus)
	applicationCheck, err := service.RunVisaApplicationCheck(queryUsStatus)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	usStatus := &models.UsStatus{Status: applicationCheck}
	res, _ := json.Marshal(usStatus)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
