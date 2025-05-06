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
	user, err := s.userSvc.SignUp(req)
	if err != nil {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)

	}
	jsonUser, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
	}

	w.Write(jsonUser)

}

func (s Server) UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
		return
	}

	var req userservice.LoginRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		http.Error(w, "Bad request - Go away!", http.StatusBadRequest)
		return

	}
	resp, err := s.userSvc.Login(req)
	if err != nil {
		http.Error(w, "BadRequest", http.StatusBadRequest)
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
	}

	w.Write(jsonResp)

}

func (s Server) NewProjectHandler(w http.ResponseWriter, r *http.Request) {

	authToken := r.Header.Get("Authorization")

	claim, err := s.controller.VerifyToken(authToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Print(err)
	}

	Req := userservice.NewProjectRequest{
		ID: claim.UserID,
	}
	err = json.Unmarshal(data, &Req)
	if err != nil {
		w.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		))

	}

	resp, err := s.userSvc.NewProject(Req, claim.UserID)
	if err != nil {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
	}

	w.Write(jsonResp)

}

func (s Server) GetProjectsHandler(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")

	claim, err := s.controller.VerifyToken(authToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	req := userservice.AllProjectRequest{
		ID: claim.UserID,
	}

	p, exist, err := s.userSvc.GetAllProject(req)
	if !exist {
		noProject := "You should create or join project first"
		w.Write([]byte(noProject))
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	jsonProject, err := json.Marshal(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(jsonProject)

}

func (s Server) GetOtherProjectHandler(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")

	claim, err := s.controller.VerifyToken(authToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	req := userservice.AllOtherProjectRequest{
		ID: claim.UserID,
	}

	p, exist, err := s.userSvc.GetAllOthersProject(req)
	if !exist {
		noProject := "You should create or join project first"
		w.Write([]byte(noProject))
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	jsonProject, err := json.Marshal(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(jsonProject)

}

func (s Server) GetProjectByIDHandler(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")

	claim, err := s.controller.VerifyToken(authToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	projectIdStr := r.URL.Query().Get("project_id")
	projectId, err := strconv.Atoi(projectIdStr)
	if err!= nil {
		http.Error(w, "Please enter number for ID", http.StatusBadRequest)

	}

	can := s.acl.CanViewProject(claim.UserID, projectId)

	if !can {
		//err.Error()
		http.Error(w, "we cannot do what you want", http.StatusForbidden)
		return

	}
	req := userservice.GetProjectByIDRequest{
		UserID: claim.UserID,
		ProjectID: projectId,

	}
	p, err := s.userSvc.GetProjectByID(req)
	if err != nil {
		http.Error(w, "Something unexpected happend", http.StatusInternalServerError)
	}
	jsonProject, err := json.Marshal(p)
	if err != nil {
		http.Error(w, "Something unexpected happend", http.StatusInternalServerError)
	}

	w.Write(jsonProject)

}

func (s Server) DeleteProjectByID(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")

	claim, err := s.controller.VerifyToken(authToken)
	if err != nil {
		//err.Error()
		http.Error(w,"You don’t have permission to access" , http.StatusForbidden)
		return
	}
	
	projectIdStr := r.URL.Query().Get("project_id")
	projectId, err := strconv.Atoi(projectIdStr)
	if err!= nil {
		http.Error(w,"internal server error" , http.StatusInternalServerError)
		
	}

	can := s.acl.CanEditProject(claim.UserID, projectId)

	if !can {
		http.Error(w, "You don’t have permission to access, you should be owner to delete project", http.StatusForbidden)
		return

	}
	_, affected, err := s.userSvc.DeleteProjectByID(projectId)
	if err != nil {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)

	}
	if !affected {
		http.Error(w, "couldnt delete", http.StatusBadRequest)
		return
	}

	jsonProject, err := json.Marshal("Project deleted")
	if err!= nil {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)

	}
	

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
