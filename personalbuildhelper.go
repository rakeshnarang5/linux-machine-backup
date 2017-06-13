package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	paramChan chan *Parameters
)

func init() {
	paramChan = make(chan *Parameters, 20)
}

type request struct {
	Result *Result `json:"result"`
}

type Result struct {
	Parameters *Parameters `json:"parameters"`
}

type Parameters struct {
	ServiceType  string
	Status       int
	Code         string `json:"build-source-type"`
	BuildType    string `json:"build-type"`
	AlexaMessage string
}

type Fulfillment struct {
	Speech      string `json:"speech"`
	Data        string `json:"data"`
	DisplayText string `json:"displayText"`
}

func hello(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return
	}
	log.Printf("method: %v\n", r.Method)
	log.Printf("URL: %+v\n", r.URL)
	header := w.Header()
	header.Add("Content-type", "application/json")
	req := &request{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return
	}
	ffm := &Fulfillment{
		Speech: "your " + req.Result.Parameters.Code + " is compiled for " + req.Result.Parameters.BuildType,
	}
	req.Result.Parameters.ServiceType = "Google Home"
	paramChan <- req.Result.Parameters
	respBody, err := json.Marshal(ffm)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return
	}
	w.Write(respBody)
	log.Printf("request body: %s\n", body)
}

func helloAlexa(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ALEXA ERROR: %v\n", err)
		return
	}
	log.Printf("ALEXA method: %v\n", r.Method)
	log.Printf("ALEXA URL: %+v\n", r.URL)
	req := &request{}
	req.Result.Parameters.ServiceType = "Alexa"
	req.Result.Parameters.AlexaMessage = string(body)
	paramChan <- req.Result.Parameters
	/*header := w.Header()
	header.Add("Content-type", "application/json")
	req := &request{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return
	}
	ffm := &Fulfillment{
		Speech: "your " + req.Result.Parameters.Code + " is compiled for " + req.Result.Parameters.BuildType,
	}
	paramChan <- req.Result.Parameters
	respBody, err := json.Marshal(ffm)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return
	}
	w.Write(respBody)*/
	w.Write([]byte("This is hisenberg from your heroku cloud"))
	log.Printf("ALEXA request body: %s\n", body)
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
	mux["/gh"] = hello
	mux["/al"] = helloAlexa
	mux["/system"] = handleSystemClient
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

func handleSystemClient(w http.ResponseWriter, r *http.Request) {
	timer := time.NewTimer(20 * time.Second)
	var parameters *Parameters
	select {
	case parameters = <-paramChan:
	case <-timer.C:
		parameters = &Parameters{
			Status: -1,
		}
	}
	w.Header().Add("Content-type", "application/json")
	data, err := json.Marshal(parameters)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return
	}
	log.Printf("response body: %s\n", data)
	w.Write(data)
}
