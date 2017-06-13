package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type test_struct struct {
	Name string
	Age  string
	City string
}

var m map[string]test_struct

func (t test_struct) String() string {
	return fmt.Sprintf("%s\t%s\t%s\n", t.Name, t.Age, t.City)
}

var count = 1

func test(rw http.ResponseWriter, req *http.Request) {
	var t test_struct
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	log.Println(t)
	m[req.URL.String()] = t
	log.Println(m)
}

func main() {
	m = make(map[string]test_struct)
	http.HandleFunc("/test", test)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
