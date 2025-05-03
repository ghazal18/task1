package delivery

import (
	"net/http"
	"task1/controller"
	"task1/service/userservice"

	"github.com/gorilla/mux"
)

type Server struct {
	userSvc    userservice.Service
	controller controller.Service
	acl        userservice.ACLGenerator
}

func New(userSvc userservice.Service, control controller.Service, acl userservice.ACLGenerator) Server {
	return Server{userSvc: userSvc, controller: control, acl: acl}
}

func (s Server) Serve() {
	router := mux.NewRouter()

	router.HandleFunc("/health-check", health_check).Methods("GET")
	router.HandleFunc("/api/v1/signup", s.UserSignup).Methods("POST")
	router.HandleFunc("/api/v1/login", s.UserLogin).Methods("POST")
	router.HandleFunc("/api/v1/auth/project/new", s.NewProject).Methods("POST")
	router.HandleFunc("/api/v1/auth/project/all", s.GetProjects).Methods("GET")
	router.HandleFunc("/api/v1/auth/project/other/all", s.GetOtherProjects).Methods("GET")
	router.HandleFunc("/api/v1/auth/project", s.GetProjectByID).Methods("GET")

	// Start the server
	http.ListenAndServe(":8010", router)
}
