package main

import (
	"bufio"
	"bytes"
	"encoding"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var queriesToWriteOnTelnetConnection = [6][]byte{
	[]byte("eng\r\n"),
	[]byte("eng1\r\n"),
	[]byte("?OUTPUT,98,1\r\n"),
	[]byte("?OUTPUT,99,1\r\n"),
	[]byte("?OUTPUT,100,1\r\n"),
	[]byte("\r\n"),
}

type structureForLIP struct {
	integrationCommand string
	integrationID      int
	actionParameter    int
	brightnessLevel    float64
}

type dataHandler struct {
	Collection map[string]scene
	//Collection map[string]scene
}

type scene struct {
	saved     bool
	comonents []*structureForLIP
}

// func (s *scene) update(b []byte) {
// 	//subscribe for the inputs

// }

func (d *dataHandler) processData(sceneName string, data []byte, integrationCommand string) {

	fmt.Println(sceneName)
	fmt.Println(string(data))

	stringData := string(data)
	querySlice := strings.Split(stringData, ",")

	fmt.Println(querySlice)

	var tempStruct structureForLIP
	tempStruct.integrationCommand = integrationCommand
	tempStruct.integrationID, _ = strconv.Atoi(querySlice[1])
	tempStruct.actionParameter, _ = strconv.Atoi(querySlice[2])
	tempStruct.brightnessLevel, _ = strconv.ParseFloat(querySlice[3], 64)

	fmt.Println(tempStruct)

	sliceOfStructs := d.Collection[sceneName]
	sliceOfStructs = append(sliceOfStructs, tempStruct)
	d.Collection[sceneName] = sliceOfStructs

	fmt.Println(d.Collection[sceneName])
}

func (s structureForLIP) MarshalText() (text []byte, err error) {
	return []byte("#" + s.integrationCommand + "," + strconv.Itoa(s.integrationID) + "," + strconv.Itoa(s.actionParameter) + "," + strconv.FormatFloat(s.brightnessLevel, 'f', -1, 64) + "\r\n"), nil
}

func (s *structureForLIP) UnmarshalText(text []byte) error {

	stringText := string(text)
	stringSlice := strings.Split(stringText, ",")

	q1 := stringSlice[0]
	s.integrationCommand = q1[1:]
	s.integrationID, _ = strconv.Atoi(stringSlice[1])
	s.actionParameter, _ = strconv.Atoi(stringSlice[2])
	s.brightnessLevel, _ = strconv.ParseFloat(stringSlice[3], 64)

	return nil
}

// func (d *dataHandler) executeScene(sceneName string) {
// 	sceneValues := d.Collection[sceneName]
// 	for _, val := range sceneValues {

// 	}

// }

func Sender(i encoding.TextMarshaler, w *bufio.Writer) {
	text, err := i.MarshalText()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(text))
	w.Write(text)
	w.Flush()
}

func main() {

	var myDataHandler dataHandler

	mapForMyDataHandler := make(map[string][]structureForLIP)

	myDataHandler.Collection = mapForMyDataHandler

	conn, err := net.Dial("tcp", "10.4.3.240:23")
	if err != nil {
		fmt.Println(err)
	}

	sceneName := "Test Scene"

	structSlice := make([]structureForLIP, 0, 0)

	myDataHandler.Collection[sceneName] = structSlice

	scanner := bufio.NewScanner(conn)
	scanner.Split(readUntilSplitter([]byte("login: "), []byte("password: "), []byte("QNET> "), []byte("\r\n")))

	writer := bufio.NewWriter(conn)

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		for {

			scanner.Scan()

			response := scanner.Bytes()

			//trimspace here

			// fmt.Println(response)
			// fmt.Printf("%s\n", response)

			emptyResponse := ""
			delimiterResponse := "\r\n"

			if string(response) != emptyResponse && string(response) != delimiterResponse[:] {
				myDataHandler.processData(sceneName, response, "OUTPUT")
			}
		}
		wg.Done()
	}()

	for i := 0; i < 6; i++ {
		writer.Write(queriesToWriteOnTelnetConnection[i])
		writer.Flush()
	}

	time.Sleep(60 * time.Second)

	fmt.Println("I reached here")

	for _, val := range myDataHandler.Collection[sceneName] {
		fmt.Println("calling the sender function")
		Sender(val, writer)
	}

	//fmt.Println(myDataHandler.Collection)

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
				if i == 0 {
					return tokenEnd, nil, nil
				}
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

// type subscriber struct {
// 	//slice of subscription
// }

// // subscribe
// // unsubscribe
// // send to client

// type subscription struct {
// 	//string OUTPUT
// 	//chaneel of bytes
// }
