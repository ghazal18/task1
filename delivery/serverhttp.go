package delivery

import (
	"net/http"
	"task1/controller"
	"task1/service/projectservice"
	"task1/service/userservice"
	uservalidator "task1/validator/user_validator"

	"github.com/gorilla/mux"
)

type Server struct {
	userSvc    userservice.Service
	projectSvc projectservice.Service
	controller controller.Service
	acl        userservice.ACLGenerator
	validator  uservalidator.Validator
}

func New(userSvc userservice.Service, control controller.Service, acl userservice.ACLGenerator, validator uservalidator.Validator,projectSvc projectservice.Service) Server {
	return Server{userSvc: userSvc,
		projectSvc: projectSvc,
		controller: control,
		acl:        acl,
		validator:  validator,
	}
}

func (s Server) Serve() {
	router := mux.NewRouter()

	router.HandleFunc("/health-check", health_check).Methods("GET")
	router.HandleFunc("/api/v1/signup", s.UserSignupHandler).Methods("POST")
	router.HandleFunc("/api/v1/login", s.UserLoginHandler).Methods("POST")
	router.HandleFunc("/api/v1/auth/project/new", s.NewProjectHandler).Methods("POST")
	router.HandleFunc("/api/v1/auth/project/all", s.GetProjectsHandler).Methods("GET")
	router.HandleFunc("/api/v1/auth/project/other/all", s.GetOtherProjectHandler).Methods("GET")
	router.HandleFunc("/api/v1/auth/project", s.GetProjectByIDHandler).Methods("GET")
	router.HandleFunc("/api/v1/auth/project", s.DeleteProjectByIDHandler).Methods("DELETE")
	router.HandleFunc("/api/v1/auth/project", s.PutProjectByIDHandler).Methods("PUT")
	router.HandleFunc("/api/v1/auth/project/join", s.JoinOtherProjectHandler).Methods("POST")
	// Start the server
	http.ListenAndServe(":8010", router)
}
