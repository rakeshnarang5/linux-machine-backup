package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var data []byte

type parameters struct {
	Actions   string
	Devices   []string
	Locations []string
}

type fulfillment struct {
	speech string
}

type test_struct struct {
	param  parameters  `json:"parameters"`
	fulfil fulfillment `json:"fulfillment"`
}

type result struct {
	res test_struct `json:"result"`
}

func main() {

	body, err := ioutil.ReadAll([]data)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
	var r result
	err = json.Unmarshal(body, &r)
	if err != nil {
		panic(err)
	}
	log.Println(r.res.param.Actions)
	log.Println(r.res.fulfil.speech)
}
