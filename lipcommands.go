package main

import (
	"fmt"
	"strings"
)

type structureForLIP struct {
	integrationID   int
	actionParameter int
	brightnessLevel int
}

// func (s *structureForLIP) string() string {
// 	return string(s.integrationID) + ", " + string(s.actionParameter) + ", " + string(s.brightnessLevel)
// }

func main() {

	brightnessLevelSlice := []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	sliceOfStructsForLIPWithDummyData := createSliceOfStructsWithDummyData(97, 1, brightnessLevelSlice)
	sceneNameVsStructSliceMap := make(map[string][]structureForLIP)
	sceneNameVsStructSliceMap["Test Scene"] = sliceOfStructsForLIPWithDummyData
	fmt.Println(sceneNameVsStructSliceMap["Test Scene"])

	sliceOfLIPCommands := createSliceOfLIPCommandsFromLIPStructure(sceneNameVsStructSliceMap["Test Scene"])
	fmt.Println(sliceOfLIPCommands)
	//	exectureLIPCommands(sliceOfLIPCommands)

	// sliceOfLIPCommandFields := []string{"#OUTPUT", "97", "1", "50"}
	// LIPCommand := strings.Join(sliceOfLIPCommandFields, ",")
	// fmt.Println(LIPCommand)

}

func createSliceOfStructsWithDummyData(integrationID int, actionParameter int, brightnessLevelSlice []int) []structureForLIP {
	sliceOfPointersToStructToReturn := make([]structureForLIP, 0, 0)
	for i := 0; i < 10; i++ {
		var oneStruct structureForLIP
		oneStruct.integrationID = integrationID
		oneStruct.actionParameter = actionParameter
		oneStruct.brightnessLevel = brightnessLevelSlice[i]
		sliceOfPointersToStructToReturn = append(sliceOfPointersToStructToReturn, oneStruct)
	}
	return sliceOfPointersToStructToReturn
}

func createSliceOfLIPCommandsFromLIPStructure(structureForLIPCommandCreation []structureForLIP) []string {
	prefix := "#OUTPUT"
	stringSliceToReturn := make([]string, 0, 0)
	for _, oneLIPStruct := range structureForLIPCommandCreation {
		sliceOfLIPCommmands := make([]string, 0, 0)
		sliceOfLIPCommmands = append(sliceOfLIPCommmands, prefix)
		sliceOfLIPCommmands = append(sliceOfLIPCommmands, string(oneLIPStruct.integrationID))
		sliceOfLIPCommmands = append(sliceOfLIPCommmands, string(oneLIPStruct.actionParameter))
		sliceOfLIPCommmands = append(sliceOfLIPCommmands, string(oneLIPStruct.brightnessLevel))
		stringSliceToReturn = append(stringSliceToReturn, strings.Join(sliceOfLIPCommmands, ","))
	}
	return stringSliceToReturn
}

// func exectureLIPCommands([]string) {
// 	//TODO
// }
