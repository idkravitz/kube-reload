package main

import (
	"os"
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
	RepoChan chan string
	Config map[string]string
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
    app.RepoChan <- data.Repository.RepoName
}

func (app *KubeReloadApp) reloader() {
	for {
		repo := <-app.RepoChan
		log.Println("Was notified that repo updated:", repo)
		rc := app.Config[repo]
		log.Println("And I gonna reload rc:", rc)
	}
}

func main() {
	file, err := os.Open("./config.json")
	if err != nil {
		log.Fatal(err)
	}
    decoder := json.NewDecoder(file)
    config := map[string]string{}
    err = decoder.Decode(&config)
	if err != nil {
        log.Fatal(err)
    }
    js, _ := json.Marshal(config)
    log.Println(string(js))

	app := KubeReloadApp{ RepoChan: make(chan string), Config: config }
	go app.reloader()
	log.Println("Listening on *:80")
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.handler)
	http.ListenAndServe(":80", mux)
}