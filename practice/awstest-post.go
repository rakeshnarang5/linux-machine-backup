package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", "https://wuqq2umcai.execute-api.us-east-1.amazonaws.com/heisenberg/mydemoresource", nil)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Host", "wuqq2umcai.execute-api.us-east-1.amazonaws.com")
	request.Header.Set("X-Amz-Date", "20170103T182526Z")
	request.Header.Set("Authorization", "AWS4-HMAC-SHA256 Credential=AKIAI5JZQ2MENT4QAAEA/20170103/us.east.1/execute-api/aws4_request, SignedHeaders=content-type;host;x-amz-date, Signature=85874253bf95349c262f6a5f1041d4bf943fb07251bea65ee7343c1285c59dc7")
	resp, _ := client.Do(request)

	//fmt.Printf("%v", resp)

	type myDemoStruct struct {
		Hello string
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var data myDemoStruct
	json.Unmarshal(body, &data)
	fmt.Println(data.Hello)

}
