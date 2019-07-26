package api

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"idcs-workshop/api/controllers"
	"log"
	"net/http"
)

func Start() {
	r := mux.NewRouter()
	initializeControllers(r)
	log.Fatal(http.ListenAndServe(":8000",
		handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}))(r)))
}

func initializeControllers(r *mux.Router) {
	controllers.InitStatusController(r)
	controllers.InitUserController(r)
}
