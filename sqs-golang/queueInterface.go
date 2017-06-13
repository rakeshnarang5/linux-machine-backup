package main

import (
	"fmt"
	"goamz/goamz/sqs"
)

//This function sends a message to an aws sqs queue
func send_msg(msg string, q *sqs.Queue) sqs.SendMessageResponse {
	res, _ := q.SendMessage(msg)
	return *res
}

//This function receives a message from the aws sqs queue
func get_msg(q *sqs.Queue) sqs.ReceiveMessageResponse {
	receiveMessageResponse, err := q.ReceiveMessage(1)
	return *receiveMessageResponse
}

//delete message
func delete_msg(receiveMessageResponse *sqs.ReceiveMessageResponse, q *sqs.Queue) sqs.DeleteMessageResponse {
	deleteMessageReponse, err := q.DeleteMessageUsingReceiptHandle(receiveMessageResponse.Messages[0].ReceiptHandle)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	return *deleteMessageReponse
}
