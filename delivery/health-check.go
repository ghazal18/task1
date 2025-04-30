package delivery

import (
	"net/http"
)

func health_check(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("SERVER is alright"))
}
