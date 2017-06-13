package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

type ProjectName struct {
	ProjectName string `xml:",attr"`
	UUID        uint32 `xml:",attr"`
}

type Area struct {
	UUID  uint32  `xml:",attr"`
	Name  string  `xml:",attr"`
	Areas []*Area `xml:">Area"`
}

type Project struct {
	XMLName     xml.Name `xml:"Project"`
	ProjectName *ProjectName
	Areas       []*Area `xml:">Area"`
}

func (area *Area) print() {
	fmt.Print(area.Name + "\t")
	fmt.Println(area.UUID)
	if area.Areas == nil {
		return
	}

	for _, innerArea := range area.Areas {
		innerArea.print()
	}
}

func main() {
	absFilePath, err := filepath.Abs("multiserver-data.xml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	file, err := os.Open(absFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	var v Project

	if err := xml.NewDecoder(file).Decode(&v); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(v.ProjectName.ProjectName)
	fmt.Println(v.ProjectName.UUID)

	for _, outerArea := range v.Areas {
		outerArea.print()
	}

}
