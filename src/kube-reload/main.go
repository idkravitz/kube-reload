package main

import (
	"io"
	"os"
	"log"
	"net/http"
	"encoding/json"
	)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "" || r.Method == "GET" {
		f, err := os.Open("index.html")
		if err != nil { log.Fatal(err) }
		io.Copy(w, f)
		return
	}

	r.ParseMultipartForm(10 << 20)

	js, _ := json.Marshal(r.PostForm)
    // decoder := json.NewDecoder(r.Body)
    // t := map[string]interface{}{}
    // err := decoder.Decode(&t)
    // if err != nil {
    //     log.Fatal(err)
    // }
    // ddd, _ := json.Marshal(&t)
    // log.Println(string(ddd))
    log.Println(string(js))
}

func main() {
	log.Println("Listening on *:80")
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	http.ListenAndServe(":80", mux)
}