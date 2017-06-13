// package main

// import (
// 	"fmt"
// 	"io"
// 	"strconv"
// )

// type structureForLIP struct {
// 	integrationCommand string
// 	integrationID      int
// 	actionParameter    int
// 	brightnessLevel    float64
// }

// type scene struct {
// 	saved      bool
// 	components []*structureForLIP
// }

// type dataHandler struct {
// 	w          *io.Writer
// 	Collection map[string]scene
// }

// func (s *structureForLIP) MarshalText() (text []byte, err error) {
// 	return []byte("#" + s.integrationCommand + "," + strconv.Itoa(s.integrationID) + "," + strconv.Itoa(s.actionParameter) + "," + strconv.FormatFloat(s.brightnessLevel, 'f', -1, 64) + "\r\n"), nil
// }

// func (s *scene) MarshalText() (text []byte, err error) {
// 	retVal := make([]byte, 0, 0)
// 	for _, val := range s.components {
// 		text, _ := val.MarshalText()
// 		for i := 0; i < len(text); i++ {
// 			retVal = append(retVal, text[i])
// 		}
// 	}
// 	return []byte(retVal), nil
// }

// func NewDataHandler() *dataHandler {
// 	var myDataHandler dataHandler
// 	mapForDataHandler := make(map[string]scene)
// 	myDataHandler.Collection = mapForDataHandler
// 	return &myDataHandler
// }

// func main() {

// 	h := NewDataHandler()

// 	dummyEntryToDataHandler(h)

// 	scene1 := h.Collection["Test Scene"]

// 	text, err := scene1.MarshalText()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(string(text))

// }

// // takes a pointer to dataHandler
// // adds a dummy entry to it
// func dummyEntryToDataHandler(h *dataHandler) {
// 	handlerMap := h.Collection
// 	structSlice := make([]*structureForLIP, 0, 0)

// 	val1 := &structureForLIP{"OUTPUT", 98, 1, 25.00}
// 	val2 := &structureForLIP{"OUTPUT", 99, 1, 50.00}
// 	val3 := &structureForLIP{"OUTPUT", 100, 1, 75.00}

// 	structSlice = append(structSlice, val1)
// 	structSlice = append(structSlice, val2)
// 	structSlice = append(structSlice, val3)

// 	var x scene
// 	x.saved = true
// 	x.components = structSlice
// 	handlerMap["Test Scene"] = x
// }

//
//
//
//

// package main

// import (
// 	"fmt"
// 	"strconv"
// 	"strings"
// )

// type structureForLIP struct {
// 	integrationCommand string
// 	integrationID      int
// 	actionParameter    int
// 	brightnessLevel    float64
// }

// func (s *structureForLIP) UnmarshalText(text []byte) error {

// 	stringText := string(text)
// 	stringSlice := strings.Split(stringText, ",")

// 	q1 := stringSlice[0]
// 	s.integrationCommand = q1[1:]
// 	s.integrationID, _ = strconv.Atoi(stringSlice[1])
// 	s.actionParameter, _ = strconv.Atoi(stringSlice[2])
// 	s.brightnessLevel, _ = strconv.ParseFloat(stringSlice[3], 64)

// 	return nil
// }

// func main() {

// 	val := []byte("#OUTPUT,98,1,100")
// 	var val1 structureForLIP
// 	err := val1.UnmarshalText(val)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(val1)

// }
