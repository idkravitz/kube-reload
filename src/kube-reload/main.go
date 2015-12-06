package main

import (
	"os"
	"fmt"
	"bytes"
	"log"
	"os/exec"
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
	KubeMasterHost string
	KubeMasterPort string
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
	if app.KubeMasterPort == "" { log.Fatal("kube_master port is not set") }
	if app.KubeMasterHost == "" { log.Fatal("kube_master host is not set")}
	os.Setenv("no_proxy", app.KubeMasterHost)
	for {
		repo := <-app.RepoChan
		log.Println("Was notified that repo updated:", repo)
		rc := app.Config[repo]
		log.Println("And I gonna reload rc:", rc)
		cmd := exec.Command("./kubectl", "-s", fmt.Sprintf("%v:%v", app.KubeMasterHost, app.KubeMasterPort), "get", "pods")
		log.Println("exec reload")
		out, _ := cmd.CombinedOutput()
		log.Println(string(out))
		// if err != nil { log.Println(err) }
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

	app := KubeReloadApp{ RepoChan: make(chan string), Config: config, KubeMasterHost: os.Getenv("KUBE_MASTER_HOST"), KubeMasterPort: os.Getenv("KUBE_MASTER_PORT") }
	go app.reloader()
	log.Println("Listening on *:80")
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.handler)
	http.ListenAndServe(":80", mux)
}