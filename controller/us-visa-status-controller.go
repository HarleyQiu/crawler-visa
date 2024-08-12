package controller

import (
	"crawler-visa/models"
	"crawler-visa/service"
	"crawler-visa/utils"
	"encoding/json"
	"net/http"
)

func StatusCheck(w http.ResponseWriter, r *http.Request) {
	queryUsStatus := &models.QueryUsStatus{}
	utils.ParseBody(r, queryUsStatus)
	applicationCheck, err := service.RunVisaStatusCheck(queryUsStatus)
	applicationCheck.Code = 200
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, _ := json.Marshal(applicationCheck)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func EmailTracking(w http.ResponseWriter, r *http.Request) {
	queryUsStatus := &models.QueryUsStatus{}
	utils.ParseBody(r, queryUsStatus)
	applicationCheck, err := service.RunVisaEmailTracking(queryUsStatus)
	applicationCheck.Code = 200
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, _ := json.Marshal(applicationCheck)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
