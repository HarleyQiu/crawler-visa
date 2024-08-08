package router

import (
	"crawler-visa/controller"
	"github.com/gorilla/mux"
)

var RegisterRouters = func(router *mux.Router) {
	router.HandleFunc("/api/us-visa-status", controller.StatusCheck).Methods("POST")
	router.HandleFunc("/api/us-visa-tracking", controller.EmailTracking).Methods("POST")
}
