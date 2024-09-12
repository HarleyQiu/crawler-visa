package router

import (
	"crawler-visa/controller"
	"github.com/gorilla/mux"
)

var RegisterRouters = func(router *mux.Router) {
	router.HandleFunc("/wuai/system/crawler_visa/us-visa-status", controller.StatusCheck).Methods("POST")
	router.HandleFunc("/wuai/system/crawler_visa/us-visa-tracking", controller.EmailTracking).Methods("POST")

	router.HandleFunc("/wuai/system/crawler_visa/cn-us/create", controller.CreateApplication).Methods("POST")
	router.HandleFunc("/wuai/system/crawler_visa/cn-us/get", controller.RetrieveApplication).Methods("GET")
}
