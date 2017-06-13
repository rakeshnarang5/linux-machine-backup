package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"sync"
)

var queriesToWriteOnTelnetConnection = [5][]byte{
	[]byte("eng\r\n"),
	[]byte("eng1\r\n"),
	// []byte("#OUTPUT,98,1,10\r\n#OUTPUT,98,1,15\r\n#OUTPUT,98,1,20\r\n#OUTPUT,98,1,30\r\n#OUTPUT,98,1,40\r\n#OUTPUT,98,1,50\r\n#OUTPUT,98,1,60\r\n#OUTPUT,98,1,70\r\n#OUTPUT,98,1,80\r\n#OUTPUT,98,1,90\r\n"),
	// []byte("#OUTPUT,99,1,10\r\n#OUTPUT,99,1,15\r\n#OUTPUT,99,1,20\r\n#OUTPUT,99,1,30\r\n#OUTPUT,99,1,40\r\n#OUTPUT,99,1,50\r\n#OUTPUT,99,1,60\r\n#OUTPUT,99,1,70\r\n#OUTPUT,99,1,80\r\n#OUTPUT,99,1,90\r\n"),
	// []byte("#OUTPUT,100,1,10\r\n#OUTPUT,100,1,15\r\n#OUTPUT,100,1,20\r\n#OUTPUT,100,1,30\r\n#OUTPUT,100,1,40\r\n#OUTPUT,100,1,50\r\n#OUTPUT,100,1,60\r\n#OUTPUT,100,1,70\r\n#OUTPUT,100,1,80\r\n#OUTPUT,100,1,90\r\n"),
	[]byte("?OUTPUT,98,1\r\n"),
	[]byte("?OUTPUT,99,1\r\n"),
	[]byte("?OUTPUT,100,1\r\n"),
}

func main() {

	dataQueriedFromQSGZones := make([][]byte, 0, 0)

	conn, err := net.Dial("tcp", "10.4.3.240:23")
	if err != nil {
		fmt.Println(err)
	}

	scanner := bufio.NewScanner(conn)
	scanner.Split(readUntilSplitter([]byte("login: "), []byte("password: "), []byte("QNET> "), []byte("\r\n")))

	writer := bufio.NewWriter(conn)
	var responseFromTelnetConnection []byte
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		for {
			scanner.Scan()

			fmt.Printf("%s\n", scanner.Bytes())
			//fmt.Println(scanner.Bytes())
		}
		wg.Done()
	}()
	for i := 0; i < 5; i++ {

		writer.Write(queriesToWriteOnTelnetConnection[i])
		writer.Flush()
		if i < 2 {
			continue
		}

		if i > 2 {
			dataQueriedFromQSGZones = append(dataQueriedFromQSGZones, responseFromTelnetConnection)
		}
	}

	scanner.Scan()
	responseFromTelnetConnection = scanner.Bytes()

	dataQueriedFromQSGZones = append(dataQueriedFromQSGZones, responseFromTelnetConnection)

	// for _, oneQuery := range dataQueriedFromQSGZones {
	// 	fmt.Println("----------------------start---------------------------")
	// 	fmt.Println(oneQuery)
	// 	fmt.Println(string(oneQuery))
	// 	fmt.Println("-----------------------end--------------------------")
	// }

	// scanner.Scan()

	// response := scanner.Bytes()

	// fmt.Println(string(response))

	// writer := bufio.NewWriter(conn)

	// username := []byte("eng\r\n")
	// writer.Write(username)
	// writer.Flush()

	// scanner.Scan()
	// response = scanner.Bytes()
	// fmt.Println(string(response))

	// writer.Write([]byte("eng1\r\n"))
	// writer.Flush()

	// scanner.Scan()
	// response = scanner.Bytes()
	// fmt.Println("-------start-------")
	// fmt.Println(response)
	// fmt.Println(string(response))
	// fmt.Println("-------end---------")

	// writer.Write([]byte("?OUTPUT,98,1\r\n"))
	// writer.Flush()

	// scanner.Scan()
	// response = scanner.Bytes()
	// fmt.Println("-------start-------")
	// fmt.Println(response)
	// fmt.Println(string(response))
	// fmt.Println("-------end---------")

	// writer.Write([]byte("?OUTPUT,99,1\r\n"))
	// writer.Flush()

	// scanner.Scan()
	// response = scanner.Bytes()
	// fmt.Println("-------start-------")
	// fmt.Println(response)
	// fmt.Println(string(response))
	// fmt.Println("-------end---------")

	// writer.Write([]byte("?OUTPUT,100,1\r\n"))
	// writer.Flush()

	// scanner.Scan()
	// response = scanner.Bytes()
	// fmt.Println("-------start-------")
	// fmt.Println(response)
	// fmt.Println(string(response))
	// fmt.Println("-------end---------")
	wg.Wait()
}

func readUntilSplitter(delims ...[]byte) bufio.SplitFunc {
	split := func(data []byte, atEOF bool) (int, []byte, error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		//fmt.Printf("Splitter : %s\n", data)
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
