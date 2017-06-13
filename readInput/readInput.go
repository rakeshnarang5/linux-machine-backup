package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"sceneCreator"
	"strings"
	"sync"
)

func main() {

	var b bytes.Buffer
	dataStore := sceneCreator.NewDataHandler(&b)

	conn, err := net.Dial("tcp", "10.4.3.240:23")
	if err != nil {
		fmt.Println(err)
	}

	reader := bufio.NewReader(os.Stdin)

	scanner := bufio.NewScanner(conn)
	scanner.Split(readUntilSplitter([]byte("login: "), []byte("password: "), []byte("QNET> "), []byte("\r\n")))

	scanner.Scan()
	//fmt.Println(string(scanner.Bytes()))
	conn.Write([]byte("eng\r\n"))

	scanner.Scan()
	//fmt.Println(string(scanner.Bytes()))
	conn.Write([]byte("eng1\r\n"))

	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func() {
		for {
			scanner.Scan()
			response := scanner.Bytes()

			emptyResponse := ""
			delimiterResponse := "\r\n"

			if string(response) != emptyResponse && string(response) != delimiterResponse[:] {
				dataStore.ProcessData(response)
			}
		}
		wg.Done()
	}()

	for {
		commandFromConsole, _ := reader.ReadString('\n')
		//fmt.Println(commandFromConsole)
		command := strings.Split(commandFromConsole, ",")
		//fmt.Println(command)
		dataStore.HandleScene(command[0], command[1])

		fmt.Println("From main: " + string(b.Bytes()))

		conn.Write(b.Bytes())

		b.Reset()

	}

	fmt.Println(dataStore)

	wg.Wait()

}

func readUntilSplitter(delims ...[]byte) bufio.SplitFunc {
	split := func(data []byte, atEOF bool) (int, []byte, error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		for _, delim := range delims {
			if i := bytes.Index(data, delim); i >= 0 {
				tokenEnd := i + len(delim)
				return tokenEnd, data[0:i], nil
			}
		}

		if atEOF {
			return len(data), nil, nil
		}

		return 0, nil, nil
	}
	return split
}
