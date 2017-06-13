package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type Output struct {
	Name string `xml:",attr"`
	UUID uint32 `xml:",attr"`
}

type Device struct {
	Name         string `xml:",attr"`
	UUID         uint32 `xml:",attr"`
	SerialNumber uint32 `xml:",attr"`
}

type DeviceGroup struct {
	Device []*Device `xml:">Device"`
}

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
	XMLName      xml.Name `xml:"Project"`
	ProjectName  *ProjectName
	Areas        []*Area        `xml:">Area"`
	DeviceGroups []*DeviceGroup `xml:">DeviceGroup"`
	Outputs      []*Output      `xml:">Output"`
}

func (pr *Area) print() {
	// fmt.Printf("%s\t%d\n", pr.Name, pr.UUID)
	// if pr.Areas != nil {
	// 	for _, ar := range pr.Areas {
	// 		ar.print()
	// 	}
	// }

	fmt.Printf("%s\t%d\n", pr.Name, pr.UUID)

	if pr.Areas == nil {
		return
	}

	for _, ar := range pr.Areas {
		ar.print()
	}
}

func (device *Device) print() {
	fmt.Printf("%s\t%d\t%d\t", device.Name, device.UUID, device.SerialNumber)
}

func (op *Output) print() {
	fmt.Print(op.Name)
	fmt.Print("\t")
	fmt.Print(op.UUID)
}

func main() {

	v := Project{}

	data, err := ioutil.ReadFile("multi-server-data.xml")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	err = xml.Unmarshal(data, &v)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	fmt.Println(v.XMLName)
	fmt.Printf("%s\t%d\n", v.ProjectName.ProjectName, v.ProjectName.UUID)

	for _, ar := range v.Areas {
		ar.print()
	}

	for _, dev := range v.DeviceGroups {
		for _, device := range dev.Device {
			device.print()
		}
	}

	for _, op := range v.Outputs {
		op.print()
	}

}
