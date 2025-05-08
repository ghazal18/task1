package delivery

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"task1/response"
	"task1/service/userservice"
	// errormsg "task1/error"
)

func WriteError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response.ErrorResponse{
		Error: message,
	})
}

func (s Server) UserSignupHandler(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		WriteError(w, 500, "Internal Server Error ")
		return
	}

	var req userservice.SignUpRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		WriteError(w, 400, "Bad Request")
		return
	}
	fieldErrors, err := s.validator.ValidateRegisterRequest(req)
	if fieldErrors != nil {
		if err != nil {
			WriteError(w, 500, err.Error())
			return
		}
		errormasage, err := json.Marshal(fieldErrors)
		if err != nil {
			WriteError(w, 500, err.Error())
			return
		}
		WriteError(w, 500, string(errormasage))

	}

	user, err := s.userSvc.SignUp(req)
	if err != nil {
		WriteError(w, 500, err.Error())
		return

	}
	jsonUser, err := json.Marshal(user)
	if err != nil {
		WriteError(w, 500, err.Error())
		return
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
	fmt.Println("UNMARSHALLAZED", resp)

	jsonResp, err := json.Marshal(resp)
	fmt.Println("MARSHALLAZED", string(jsonResp))
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
	fmt.Println("claim.UserID", claim.UserID)

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
	if err != nil {
		http.Error(w, "Please enter number for ID", http.StatusBadRequest)

	}

	can := s.acl.CanViewProject(claim.UserID, projectId)

	if !can {
		//err.Error()
		http.Error(w, "we cannot do what you want", http.StatusForbidden)
		return

	}
	req := userservice.GetProjectByIDRequest{
		UserID:    claim.UserID,
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

func (s Server) DeleteProjectByIDHandler(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")

	claim, err := s.controller.VerifyToken(authToken)
	if err != nil {
		//err.Error()
		http.Error(w, "You don’t have permission to access", http.StatusForbidden)
		return
	}

	projectIdStr := r.URL.Query().Get("project_id")
	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)

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
	if err != nil {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)

	}

	w.Write(jsonProject)

}

func (s Server) JoinOtherProjectHandler(w http.ResponseWriter, r *http.Request) {

	authToken := r.Header.Get("Authorization")

	claim, err := s.controller.VerifyToken(authToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	projectIdStr := r.URL.Query().Get("project_id")
	req := userservice.JoinProjectByIDRequest{
		ProjectID: projectIdStr,
		UserID:    strconv.Itoa(claim.UserID),
	}

	done, err := s.userSvc.JoinProjectByID(req)
	if err != nil {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)

	}
	if !done {
		http.Error(w, "couldnt Join this project", http.StatusBadRequest)
		return
	}
	res := "YOU join the project"
	jsonRes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)

	}

	w.Write(jsonRes)

}

func (s Server) PutProjectByIDHandler(w http.ResponseWriter, r *http.Request) {

	authToken := r.Header.Get("Authorization")

	claim, err := s.controller.VerifyToken(authToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	projectIdStr := r.URL.Query().Get("project_id")
	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)

	}

	can := s.acl.CanEditProject(claim.UserID, projectId)

	if !can {
		http.Error(w, "You don’t have permission to access, you should be owner to Edit project", http.StatusForbidden)
		return

	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request - Go away!", http.StatusBadRequest)
	}

	var p userservice.PutProjectByIDRequest
	err = json.Unmarshal(data, &p)
	if err != nil {
		panic(err)
	}
	fmt.Println(p)
	res, ok, err := s.userSvc.UpdateProjectByID(projectIdStr, p)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)

	}
	if !ok {
		http.Error(w, "couldn't update project", http.StatusBadRequest)
		return
	}
	jsonRes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)

	}

	w.Write(jsonRes)

}
