package main

import (
	"io"
	"os"
	"log"
	"net/http"
	"encoding/json"
	)


type UpdateData_Repository {
	RepoName string `json:"repo_name"`

}
type UpdateData struct {
	CallbackUrl string `json:"callback_url"`
	Repository UpdateData_Repository `json:"repository"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "" || r.Method == "GET" {
		f, err := os.Open("index.html")
		if err != nil { log.Fatal(err) }
		io.Copy(w, f)
		return
	}

	// r.ParseMultipartForm(10 << 20)

	// js, _ := json.Marshal(r.PostForm)
    decoder := json.NewDecoder(r.Body)
    data := UpdateData{}
    // t := map[string]interface{}{}
    err := decoder.Decode(&data)
    if err != nil {
        log.Fatal(err)
    }
    js, _ := json.Marshal(&data)
    log.Println(string(js))
    // log.Println(string(js))
}

func main() {
	log.Println("Listening on *:80")
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	http.ListenAndServe(":80", mux)
}