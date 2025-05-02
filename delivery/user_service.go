package delivery

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"task1/service/userservice"
)

func (s Server) UserSignup(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Print(err)
	}

	var uReq userservice.RegisterRequest
	err = json.Unmarshal(data, &uReq)
	if err != nil {
		w.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		))

	}
	user, err := s.userSvc.Register(uReq)
	jsonUser, err := json.Marshal(user)

	w.Write(jsonUser)

}

func (s Server) UserLogin(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Print(err)
	}

	var uReq userservice.LoginRequest
	err = json.Unmarshal(data, &uReq)
	if err != nil {
		w.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		))

	}
	s.userSvc.Login(uReq)

}

func (s Server) NewProject(w http.ResponseWriter, r *http.Request) {

	authToken := r.Header.Get("Authorization")

	_, err := s.controller.VerifyToken(authToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// resp, err := s.userSvc.NewProject(userservice.NewProjectRequest{UserID: claims.UserID})
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusForbidden)

	// }

	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Print(err)
	}

	var Req userservice.NewProjectRequest
	err = json.Unmarshal(data, &Req)
	if err != nil {
		w.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		))

	}
	s.userSvc.NewProject(Req)
}
