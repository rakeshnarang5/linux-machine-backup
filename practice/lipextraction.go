package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {

	type zoneParameters struct {
		integrationID   int
		actionParameter int
		brightnessLevel float64
	}

	responseFromLIPConnection := "~OUTPUT,98,1,85.00"

	LIPResponseComponents := strings.Split(responseFromLIPConnection, ",")

	fmt.Println(LIPResponseComponents)

	var LIPParameters98 zoneParameters

	LIPParameters98.integrationID, _ = strconv.Atoi(LIPResponseComponents[1])
	LIPParameters98.actionParameter, _ = strconv.Atoi(LIPResponseComponents[2])
	LIPParameters98.brightnessLevel, _ = strconv.ParseFloat(LIPResponseComponents[3], 64)

	fmt.Println(LIPParameters98)

	LIPCommandComponents := make([]string, 0, 0)
	LIPCommandComponents
}
