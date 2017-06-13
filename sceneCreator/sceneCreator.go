package sceneCreator

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type structureForLIP struct {
	integrationCommand string
	integrationID      int
	actionParameter    int
	brightnessLevel    float64
}

const (
	device1            = "90"
	device2            = "94"
	device3            = "106"
	integrationCommand = "OUTPUT"
)

type scene struct {
	saved      bool
	components []*structureForLIP
}

type dataHandler struct {
	w          io.Writer
	Collection map[string]scene
}

type SceneHandler interface {
	ProcessData([]byte)
	HandleScene(string, string)
}

func NewDataHandler(w io.Writer) SceneHandler {
	var myDataHandler dataHandler
	myDataHandler.w = w
	mapForDataHandler := make(map[string]scene)
	myDataHandler.Collection = mapForDataHandler
	return &myDataHandler
}

func (s *structureForLIP) MarshalText() (text []byte, err error) {
	return []byte("#" + s.integrationCommand + "," + strconv.Itoa(s.integrationID) + "," + strconv.Itoa(s.actionParameter) + "," + strconv.FormatFloat(s.brightnessLevel, 'f', -1, 64) + "\r\n"), nil
}

func (s *scene) MarshalText() (text []byte, err error) {
	retVal := make([]byte, 0, 0)
	for _, val := range s.components {
		text, _ := val.MarshalText()
		for i := 0; i < len(text); i++ {
			retVal = append(retVal, text[i])
		}
	}
	return []byte(retVal), nil
}

func (s *structureForLIP) UnmarshalText(text []byte) error {

	stringText := string(text)
	stringSlice := strings.Split(stringText, ",")

	var err error

	if !strings.Contains(stringText, integrationCommand) {
		return errors.New("Unsupported integration command")
	}

	fmt.Println(stringSlice[1] + ", " + device1 + device2 + device3)

	if !(stringSlice[1] == device1 || stringSlice[1] == device2 || stringSlice[1] == device3) {
		return errors.New("Unsupported device integration id")
	}

	s.integrationCommand = integrationCommand
	s.integrationID, err = strconv.Atoi(stringSlice[1])
	if err != nil {
		return err
	}
	s.actionParameter, err = strconv.Atoi(stringSlice[2])
	if err != nil {
		return err
	}
	s.brightnessLevel, err = strconv.ParseFloat(stringSlice[3], 64)
	if err != nil {
		return err
	}

	return nil
}

func (d *dataHandler) ProcessData(data []byte) {

	fmt.Println("reached process data")
	fmt.Println(string(data))

	if data[0] == 0 {
		data = data[1:]
	}

	sceneMap := d.Collection
	var sceneToModifyString string

	for k, v := range sceneMap {
		if v.saved == false {
			sceneToModifyString = k
		}

	}

	fmt.Println("Processing data for: " + sceneToModifyString)

	sceneToModify := sceneMap[sceneToModifyString]

	var x structureForLIP
	err := x.UnmarshalText(data)
	if err != nil {
		fmt.Printf("Error occured: %v\n", err)
		return
	}

	sceneToModify.components = append(sceneToModify.components, &x)

	if len(sceneToModify.components) == 3 {
		sceneToModify.saved = true
	}

	fmt.Println(sceneToModify)

	sceneMap[sceneToModifyString] = sceneToModify

}

func (d *dataHandler) HandleScene(action string, sceneName string) {
	switch strings.ToUpper(action) {
	case "CREATE", "UPDATE":
		fmt.Println("Creating a new scene.")
		structSlice := make([]*structureForLIP, 0, 0)
		var myScene scene
		myScene.saved = false
		myScene.components = structSlice
		d.Collection[sceneName] = myScene
		queriesForState := []byte("?OUTPUT," + device1 + ",1\r\n?OUTPUT," + device2 + ",1\r\n?OUTPUT," + device3 + ",1\r\n")
		w := d.w
		w.Write(queriesForState)
	case "DELETE":
		delete(d.Collection, sceneName)
	case "ACTIVATE":
		fmt.Println("Activating scene")
		scene := d.Collection[sceneName]
		fmt.Println(scene)
		text, err := scene.MarshalText()
		fmt.Println(text)
		if err != nil {
			fmt.Println(err)
		}
		w := d.w
		w.Write(text)
	// case "SHOWMAP":
	// 	fmt.Println(d.Collection)
	default:

	}
}
