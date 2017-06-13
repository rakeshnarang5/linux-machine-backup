package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	paramChan chan *request
	ffmChan   chan []byte
)

func init() {
	paramChan = make(chan *request, 20)
	ffmChan = make(chan []byte, 20)
}

type request struct {
	Result *Result `json:"result"`
}

type Result struct {
	ServiceType   string
	ResolvedQuery string      `json:"resolvedQuery"`
	MetaData      *MetaData   `json:"metadata"`
	Parameters    *Parameters `json:"parameters"`
	//Fulfillment   *Fulfillment `json: "fulfillment"`
}

type MetaData struct {
	IntentName string `json:"intentName"`
	IntenID    string
}

type Parameters struct {
	Status       int
	Action       string
	Location     string
	Device       string
	AlexaMessage string
	SceneName    string
}

type Fulfillment struct {
	Speech   string   `json:"speech"`
	Messages []string `json:"message`
}

func sceneHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("HANDLING GOOGLE HOME SCENE FROM SERVER")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
	}
	log.Printf("request body: %s\n", body)
	header := w.Header()
	header.Add("Content-type", "application/json")
	req := &request{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return
	}
	i := strings.Index(req.Result.ResolvedQuery, "scene")
	req.Result.Parameters.SceneName = req.Result.ResolvedQuery[i+6:]

	log.Println("ACTION IDENTIFIED: " + req.Result.Parameters.Action)
	log.Println("SCENE NAME IDENTIFIED: " + req.Result.Parameters.SceneName)

	paramChan <- req
	ffm := <-ffmChan
	log.Printf("received fullfillment: %s\n", ffm)
	w.Write(ffm)
}

var mux map[string]func(http.ResponseWriter, *http.Request)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	server := http.Server{
		Addr:    ":" + port,
		Handler: &myHandler{},
	}
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/post"] = sceneHandler
	mux["/get"] = handleSystemClientScene
	server.ListenAndServe()
}

type myHandler struct{}

//comment
func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}
	io.WriteString(w, "My server: "+r.URL.String())
}

func handleSystemClientScene(w http.ResponseWriter, r *http.Request) {
	log.Println("HANDLING GOOGLE HOME SCENE FROM CLIENT")
	w.Header().Add("Content-type", "application/json")
	timer := time.NewTimer(20 * time.Second)
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("System Client ERROR: %v\n", err)
		}
		ffmChan <- body
	}
	var req *request
	select {
	case req = <-paramChan:
	case <-timer.C:
		req = &request{
			Result: &Result{
				Parameters: &Parameters{
					Status: -1,
				},
			},
		}
	}
	data, err := json.Marshal(req)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return
	}
	log.Printf("response body: %s\n", data)
	w.Write(data)
}
