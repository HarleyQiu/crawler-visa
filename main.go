package main

import (
	"crawler-visa/router"
	"crawler-visa/scheduler"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	router.RegisterRouters(r)
	scheduler.RunScheduledTasks()

	log.Println("Server is starting on 0.0.0.0:9010...")
	log.Fatal(http.ListenAndServe("0.0.0.0:9010", setupCORS(r)))
}

// setupCORS wraps the router with CORS settings
func setupCORS(r *mux.Router) http.Handler {
	return handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),                                                                                            // 允许任何来源，生产环境中应指定明确的域名
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),                                                      // 允许的HTTP方法
		handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"}), // 允许的HTTP头部
	)(r)
}
