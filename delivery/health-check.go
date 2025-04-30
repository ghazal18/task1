package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func health_check(w http.ResponseWriter, r *http.Request) {

	rawText := map[string]string{
		"massage": "everyThing is fine!",
	}
	resp, err := json.Marshal(rawText)

	if err != nil {
		fmt.Println(err)
	}

	w.Write(resp)
}
