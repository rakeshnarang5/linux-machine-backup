package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter Command")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)

	userCommand := strings.Split(text, ",")

}
