package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

var queriesToWriteOnTelnetConnection = [5][]byte{
	[]byte("eng\r\n"),
	[]byte("eng1\r\n"),
	[]byte("?OUTPUT,98,1\r\n"),
	[]byte("?OUTPUT,99,1\r\n"),
	[]byte("?OUTPUT,100,1\r\n"),
}

type structureForLIP struct {
	integrationID   int
	actionParameter int
	brightnessLevel float64
}

func main() {

	makeConnectionAndQueryData()

	// sliceOfLIPStructures := createStructureFromLIPResponse(dataQueriedFromQSGZones)

	// mapOfSceneNameVsSliceOfLIPStruct := make(map[string][]structureForLIP)
	// mapOfSceneNameVsSliceOfLIPStruct["Test Scene"] = sliceOfLIPStructures

	// stringSliceOfLIPCommands := createLIPCommandsFromLIPStructure(mapOfSceneNameVsSliceOfLIPStruct["Test Scene"])

	// fmt.Println(stringSliceOfLIPCommands)
}

func makeConnectionAndQueryData() {
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

	fmt.Println(dataQueriedFromQSGZones)
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

func createLIPCommandsFromLIPStructure(structureForLIPCommandCreation []structureForLIP) []string {
	prefix := "#OUTPUT"
	stringSliceToReturn := make([]string, 0, 0)
	for _, oneLIPStruct := range structureForLIPCommandCreation {
		sliceOfLIPCommmands := make([]string, 0, 0)
		sliceOfLIPCommmands = append(sliceOfLIPCommmands, prefix)
		sliceOfLIPCommmands = append(sliceOfLIPCommmands, strconv.Itoa(oneLIPStruct.integrationID))
		sliceOfLIPCommmands = append(sliceOfLIPCommmands, strconv.Itoa(oneLIPStruct.actionParameter))
		sliceOfLIPCommmands = append(sliceOfLIPCommmands, strconv.FormatFloat(oneLIPStruct.brightnessLevel, 'f', -1, 64))
		stringSliceToReturn = append(stringSliceToReturn, strings.Join(sliceOfLIPCommmands, ","))
	}
	return stringSliceToReturn
}

func createStructureFromLIPResponse(dataQueriedFromQSGZones []string) []structureForLIP {

	sliceOfLIPStructures := make([]structureForLIP, 0, 0)

	for _, oneQuery := range dataQueriedFromQSGZones {
		var tempStruct structureForLIP
		extractedLIPCommandFields := strings.Split(oneQuery, ",")
		tempStruct.integrationID, _ = strconv.Atoi(extractedLIPCommandFields[0])
		tempStruct.actionParameter, _ = strconv.Atoi(extractedLIPCommandFields[1])
		tempStruct.brightnessLevel, _ = strconv.ParseFloat(extractedLIPCommandFields[2], 64)

		fmt.Println(tempStruct)

		sliceOfLIPStructures = append(sliceOfLIPStructures, tempStruct)

	}

	return sliceOfLIPStructures

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
