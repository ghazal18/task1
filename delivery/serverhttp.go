package delivery

import (
	"net/http"

	"github.com/gorilla/mux"
)


type Server struct {
}

func New() Server {
	return Server{}
}

func(s Server) Serve() {
	router := mux.NewRouter()

	router.HandleFunc("/health-check", health_check).Methods("GET")
	router.HandleFunc("/api/v1/signup", s.ServeHTTP).Methods("POST")

	// Start the server
	http.ListenAndServe(":8010", router)
}
