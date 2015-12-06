package main

import (
	"bytes"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	)


type UpdateData_Repository struct {
	RepoName string `json:"repo_name"`
}
type UpdateData struct {
	CallbackUrl string `json:"callback_url"`
	Repository UpdateData_Repository `json:"repository"`
}

func sayThanks(callbackUrl string) {
    var respStr = []byte(`{"state": "success", "description": "Thank you very much!"}`)
    log.Println("Going to call", callbackUrl)
    req, err := http.NewRequest("POST", callbackUrl, bytes.NewBuffer(respStr))
    req.Header.Set("Content-Type", "application/json")
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    log.Println("response Status:", resp.Status)
    log.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    log.Println("response Body:", string(body))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

    decoder := json.NewDecoder(r.Body)
    data := UpdateData{}
    err := decoder.Decode(&data)
    if err != nil {
        log.Fatal(err)
    }
    sayThanks(data.CallbackUrl)
}

func main() {
	log.Println("Listening on *:80")
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	http.ListenAndServe(":80", mux)
}