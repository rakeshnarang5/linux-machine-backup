package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	url := "http://shrouded-beyond-17924.herokuapp.com/al"
	fmt.Println("URL:>", url)

	var jsonStr = []byte(`{"name": "Narang"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-type", "application/json")
	// req.Header.Set("Host", "wuqq2umcai.execute-api.us-east-1.amazonaws.com")
	// req.Header.Set("X-Amz-Date", "20170103T182526Z")
	// req.Header.Set("Authorization", "AWS4-HMAC-SHA256 Credential=AKIAI5JZQ2MENT4QAAEA/20170103/us.east.1/execute-api/aws4_request, SignedHeaders=content-type;host;x-amz-date, Signature=85874253bf95349c262f6a5f1041d4bf943fb07251bea65ee7343c1285c59dc7")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status", resp.Status)
	fmt.Println("Response headers", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response body", string(body))
}
