package main

import (
	"crawler-visa/router"
	"crawler-visa/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	corsMiddleware := utils.SetupCORS() // 设置CORS策略
	router.RegisterRouters(r)

	//http.Handle("/", corsObj(r))
	log.Println("Server is starting on 0.0.0.0:9010...")
	log.Fatal(http.ListenAndServe("0.0.0.0:9010", corsMiddleware(r)))
}
