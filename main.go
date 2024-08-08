package main

import (
	"crawler-visa/router"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	// 设置CORS策略
	corsObj := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),                                                                                            // 允许任何来源，生产环境中应指定明确的域名
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),                                                      // 允许的HTTP方法
		handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"}), // 允许的HTTP头部
	)
	router.RegisterRouters(r)
	//http.Handle("/", corsObj(r))
	log.Println("Server is starting on localhost:9010...")
	log.Fatal(http.ListenAndServe("0.0.0.0:9010", corsObj(r)))
}
