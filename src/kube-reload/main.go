package main

import (
	"io"
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

func sayThanks(callbackUrl string) {
    var respStr = []byte(`{"state": "success", "description": "Thank you very much!"}`)
    log.Println("Going to call", callbackUrl)
    req, err := http.NewRequest("POST", callbackUrl, bytes.NewBuffer(respStr))
    // req.Header.Set("X-Custom-Header", "myvalue")
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
    sayThanks(data.CallbackUrl)
    // js, _ := json.Marshal(&data)
    // log.Println(string(js))

    // url := "http://restapi3.apiary.io/notes"
    // fmt.Println("URL:>", url)


    // log.Println(string(js))
}

func main() {
	log.Println("Listening on *:80")
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	http.ListenAndServe(":80", mux)
}