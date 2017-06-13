package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	resp, err := http.Get("http://shrouded-beyond-17924.herokuapp.com/")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	fmt.Println(string(body))
}
