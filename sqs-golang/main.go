package main

import (
	"fmt"
)

func main() {
	q := queue_init()
	fmt.Println(q)

	send_msg_response := send_msg("hello", &q)
	fmt.Println(send_msg_response)

	receive_msg_response := get_msg(&q)
	fmt.Println(receive_msg_response.Message[0].Body)

	delete_msg_response := delete_msg(&receiveMessageResponse, &q)
	fmt.Println(delete_msg_response)
}
