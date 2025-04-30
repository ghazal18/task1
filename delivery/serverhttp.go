package delivery

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)


type homeHandler struct{}

func Serve2(){
	router := mux.NewRouter()

    // Register the route
    home := homeHandler{}

    router.HandleFunc("/api/v1/signup", home.ServeHTTP).Methods("POST")

    // Start the server
    http.ListenAndServe(":8010", router)
}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("This is my home page"))
	str,err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(str)
}

