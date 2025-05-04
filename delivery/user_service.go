package delivery

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

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

func (s Server) GetProjects(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")

	claim, err := s.controller.VerifyToken(authToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	println(claim.UserID)

	p, err := s.userSvc.GetAllProject(claim.UserID)

	jsonProject, err := json.Marshal(p)

	w.Write(jsonProject)

}

func (s Server) GetOtherProjects(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")

	claim, err := s.controller.VerifyToken(authToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	println(claim.UserID)

	p, err := s.userSvc.GetAllOthersProject(claim.UserID)

	jsonProject, err := json.Marshal(p)

	w.Write(jsonProject)

}

func (s Server) GetProjectByID(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")

	claim, err := s.controller.VerifyToken(authToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	projectIdStr := r.URL.Query().Get("project_id")
	projectId, err := strconv.Atoi(projectIdStr)
	println("claim.UserID, projectId", claim.UserID, projectId)

	can := s.acl.CanViewProject(claim.UserID, projectId)

	if !can {
		fmt.Println(can)
		http.Error(w, err.Error(), http.StatusForbidden)

	}
	p, err := s.userSvc.GetProjectByID(projectId)
	jsonProject, err := json.Marshal(p)

	w.Write(jsonProject)

}
