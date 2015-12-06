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
type KubeReloadApp struct {
	repoChan chan string
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

func (app *KubeReloadApp) handler(w http.ResponseWriter, r *http.Request) {
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
    app.repoChan <- data.Repository.RepoName
}

func (app *KubeReloadApp) reloader() {
	for {
		repo := <-app.repoChan
		log.Println(repo)
	}
}

func main() {
	app := KubeReloadApp{ repoChan: make(chan string) }
	go app.reloader()
	log.Println("Listening on *:80")
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.handler)
	http.ListenAndServe(":80", mux)
}