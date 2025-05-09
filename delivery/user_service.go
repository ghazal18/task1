package delivery

import (
	"encoding/json"
	"io"
	"net/http"
	


	"task1/service/userservice"
)


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
		WriteError(w, 500, "Internal Server Error")
		return
	}

	var req userservice.LoginRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		WriteError(w, 400, "Bad Request")
		return
	}
	resp, err := s.userSvc.Login(req)
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
