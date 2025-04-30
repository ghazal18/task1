package delivery

import (
	"encoding/json"
	"task1/service/userservice"

	"fmt"
	"io"
	"net/http"
)

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	s.userSvc.Register(uReq)

}
