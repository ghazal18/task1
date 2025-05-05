package delivery

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"task1/entity"
	"task1/service/userservice"
)

func (s Server) UserSignupHandler(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request - Go away!", http.StatusBadRequest)
	}

	var req userservice.SignUpRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}
	if err, fieldErrors := s.validator.ValidateRegisterRequest(req); err != nil {
		errormasage, err := json.Marshal(fieldErrors)
		if err != nil {
			fmt.Errorf("can't marshall error")
		}

		http.Error(w, string(errormasage), http.StatusBadRequest)
		return
	}
	user, err := s.userSvc.SignUpService(req)
	if err != nil {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)

	}
	jsonUser, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
	}

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

func (s Server) DeleteProjectByID(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")

	claim, err := s.controller.VerifyToken(authToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	projectIdStr := r.URL.Query().Get("project_id")
	projectId, err := strconv.Atoi(projectIdStr)
	println("claim.UserID, projectId", claim.UserID, projectId)

	can := s.acl.CanEditProject(claim.UserID, projectId)

	if !can {
		fmt.Println(can)
		http.Error(w, err.Error(), http.StatusForbidden)

	}
	p, err := s.userSvc.DeleteProjectByID(projectId)
	fmt.Println("dude", p)

	jsonProject, err := json.Marshal(p)

	w.Write(jsonProject)

}

func (s Server) PutProjectByID(w http.ResponseWriter, r *http.Request) {
	/*
		data, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Print(err)
		}

		var uReq entity.Project
		err = json.Unmarshal(data, &uReq)
		if err != nil {
			w.Write([]byte(
				fmt.Sprintf(`{"error": "%s"}`, err.Error()),
			))

		}
		fmt.Println("this is the requesti object",uReq)
		projectIdStr := r.URL.Query().Get("project_id")
		s.userSvc.UpdateProjectByID(projectIdStr,uReq)
	*/

	var input map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	jsonData, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}

	// Convert JSON to struct
	var p entity.Project
	err = json.Unmarshal(jsonData, &p)
	if err != nil {
		panic(err)
	}
	fmt.Println(p)

}
func (s Server) JoinOtherProject(w http.ResponseWriter, r *http.Request) {

	authToken := r.Header.Get("Authorization")

	claim, err := s.controller.VerifyToken(authToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	projectIdStr := r.URL.Query().Get("project_id")

	s.userSvc.JoinProjectByID(projectIdStr, strconv.Itoa(claim.UserID))

}
