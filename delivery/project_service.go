package delivery

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"task1/service/userservice"
)

func (s Server) NewProjectHandler(w http.ResponseWriter, r *http.Request) {

	authToken := r.Header.Get("Authorization")

	claim, err := s.controller.VerifyToken(authToken)
	if err != nil {
		WriteError(w, 401, "Unauthorized")
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		WriteError(w, 500, "Internal Server Error")
		return
	}

	Req := userservice.NewProjectRequest{
		ID: claim.UserID,
	}
	err = json.Unmarshal(data, &Req)
	if err != nil {
		WriteError(w, 400, "Bad Request")
		return

	}

	resp, err := s.userSvc.NewProject(Req, claim.UserID)
	if err != nil {
		WriteError(w, 400, err.Error())
		return

	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		WriteError(w, 500, err.Error())
		return
	}

	w.Write(jsonResp)

}

func (s Server) GetProjectsHandler(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")

	claim, err := s.controller.VerifyToken(authToken)
	if err != nil {
		WriteError(w, 401, "Unauthorized")
		return
	}

	req := userservice.AllProjectRequest{
		ID: claim.UserID,
	}
	p, err := s.userSvc.GetAllProject(req)

	if err != nil {
		WriteError(w, 400, err.Error())
		return
	}
	jsonProject, err := json.Marshal(p)
	if err != nil {
		WriteError(w, 500, err.Error())
		return
	}

	w.Write(jsonProject)

}

func (s Server) GetOtherProjectHandler(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")

	claim, err := s.controller.VerifyToken(authToken)
	if err != nil {
		WriteError(w, 401, "Unauthorized")
		return
	}

	req := userservice.AllOtherProjectRequest{
		ID: claim.UserID,
	}

	p, err := s.userSvc.GetAllOthersProject(req)
	if err != nil {
		WriteError(w, 400, err.Error())
		return
	}
	jsonProject, err := json.Marshal(p)
	if err != nil {
		WriteError(w, 500, err.Error())
		return
	}
	w.Write(jsonProject)

}

func (s Server) GetProjectByIDHandler(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")

	claim, err := s.controller.VerifyToken(authToken)
	if err != nil {
		WriteError(w, 401, "Unauthorized")
		return
	}

	projectIdStr := r.URL.Query().Get("project_id")
	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		WriteError(w, 400, "please enter validate project id")
		return
	}

	can := s.acl.CanViewProject(claim.UserID, projectId)

	if !can {
		WriteError(w, 403, "Forbidden")
		return
	}
	req := userservice.GetProjectByIDRequest{
		UserID:    claim.UserID,
		ProjectID: projectId,
	}
	p, err := s.userSvc.GetProjectByID(req)
	if err != nil {
		WriteError(w, 400, err.Error())
		return
	}
	jsonProject, err := json.Marshal(p)
	if err != nil {
		WriteError(w, 500, err.Error())
		return
	}
	w.Write(jsonProject)

}

func (s Server) DeleteProjectByIDHandler(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")

	claim, err := s.controller.VerifyToken(authToken)
	if err != nil {
		WriteError(w, 401, "Unauthorized")
		return
	}

	projectIdStr := r.URL.Query().Get("project_id")
	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		WriteError(w, 400, "please enter validate project id")
		return
	}

	can := s.acl.CanEditProject(claim.UserID, projectId)

	if !can {
		WriteError(w, 403, "Forbidden")
		return
	}
	_, affected, err := s.userSvc.DeleteProjectByID(projectId)
	if err != nil {
		WriteError(w, 400, err.Error())
		return
	}
	if !affected {
		WriteError(w, 400, "couldnt't delete")
		return
	}
	respon := map[string]string{"status":"Project deleted"}

	jsonProject, err := json.Marshal(respon)
	if err != nil {
		WriteError(w, 500, "Internal Server Error")
		return
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
	_, err = strconv.Atoi(projectIdStr)
	if err != nil {
		WriteError(w, 400, "please enter validate project id")
		return
	}
	
	req := userservice.JoinProjectByIDRequest{
		ProjectID: projectIdStr,
		UserID:    strconv.Itoa(claim.UserID),
	}

	done, err := s.userSvc.JoinProjectByID(req)
	if err != nil {
		WriteError(w, 500, err.Error())
		return
	}
	if !done {
		WriteError(w, 500, "you couldn't join project")
		return
	}

	respon := map[string]string{"status":"joined the project"}
	jsonRes, err := json.Marshal(respon)
	if err != nil {
		WriteError(w, 500, "Internal Server Error")
		return
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
		WriteError(w, 400, "please enter validate project id")
		return
	}

	can := s.acl.CanEditProject(claim.UserID, projectId)

	if !can {
		WriteError(w, 403, "Forbidden")
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		WriteError(w, 400, "Bad Request")
		return
	}

	var p userservice.PutProjectByIDRequest
	err = json.Unmarshal(data, &p)
	if err != nil {
		WriteError(w, 400, "Bad Request")
		return

	}
	
	res, ok, err := s.userSvc.UpdateProjectByID(projectId, p)
	if err != nil {
		WriteError(w, 500, err.Error())
		return
	}
	if !ok {
		WriteError(w, 400, "couldn't update project")
		return
	}
	jsonRes, err := json.Marshal(res)
	if err != nil {
		WriteError(w, 500, err.Error())
		return
	}

	w.Write(jsonRes)

}
