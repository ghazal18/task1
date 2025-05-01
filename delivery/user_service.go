package delivery

import (
	"encoding/json"
	"task1/service/userservice"

	"fmt"
	"io"
	"net/http"
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
	user,err := s.userSvc.Register(uReq)
	jsonUser,err:= json.Marshal(user)

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