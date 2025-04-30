package delivery

import (
	"net/http"
	"task1/service/userservice"

	"github.com/gorilla/mux"
)

type Server struct {
	userSvc userservice.Service
}

func New(userSvc userservice.Service) Server {
	return Server{userSvc:userSvc}
}

func (s Server) Serve() {
	router := mux.NewRouter()

	router.HandleFunc("/health-check", health_check).Methods("GET")
	router.HandleFunc("/api/v1/signup", s.ServeHTTP).Methods("POST")

	// Start the server
	http.ListenAndServe(":8010", router)
}
