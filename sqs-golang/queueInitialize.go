package main

import (
	"fmt"
	"goamz/goamz/sqs"
)

func queue_init() sqs.Queue {
	var access_key = "AKIAI3H6N3N5DBUNGVRQ"
	var secret_key = "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"
	var queue_url = "https://sqs.us-east-1.amazonaws.com/159364195612/MyQueue"
	var region_name = "us.east.1"
	var queue_name = "MyQueue"

	conn := establish_connection(access_key, secret_key, region_name)

	q, _ := conn.GetQueue(queue_name)

	fmt.Println("queue does not exist")
	q = conn.QueueFromArn(queue_url)

	return *q
}

// This function returns a pointer to the queue object
func establish_connection(access_key string, secret_key string, region_name string) *sqs.SQS {
	conn, err := sqs.NewFrom(access_key, secret_key, region_name)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	return conn
}
