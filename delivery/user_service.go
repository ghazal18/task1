package delivery

import(
	"fmt"
	"io"
	"net/http"
)

func (h Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("This is my home page"))
	str,err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(str)
}