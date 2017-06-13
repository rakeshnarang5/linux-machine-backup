package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type XMLStrap struct {
	XMLName xml.Name `xml:"strap"`
	Key     string   `xml:"key,attr"`
	Value   string   `xml:"value,attr"`
}

type XMLStraps struct {
	XMLName xml.Name   `xml:"straps"`
	Straps  []XMLStrap `xml:"strap"`
}

func ReadStraps(reader io.Reader) ([]XMLStrap, error) {
	var xmlStraps XMLStraps
	if err := xml.NewDecoder(reader).Decode(&xmlStraps); err != nil {
		return nil, err
	}

	return xmlStraps.Straps, nil
}

func main() {
	// Build the location of the straps.xml file
	// filepath.Abs appends the file name to the default working directly
	strapsFilePath, err := filepath.Abs("straps.xml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Open the straps.xml file
	file, err := os.Open(strapsFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	// Read the straps file
	xmlStraps, err := ReadStraps(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Display The first strap
	fmt.Printf("Key: %s  Value: %s", xmlStraps[0].Key, xmlStraps[0].Value)
}
