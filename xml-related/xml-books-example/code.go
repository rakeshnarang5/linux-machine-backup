package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

func main() { //sorted
	type OneBook struct {
		XMLName     xml.Name `xml:"book"`
		Id          string   `xml:"id,attr"`
		Author      string   `xml:"author"`
		Title       string   `xml:"title"`
		Price       float32  `xml:"price"`
		description string   `xml:"description"`
	}

	type AllBooks struct {
		BooksList []OneBook `xml:"book"`
	}

	v := AllBooks{}

	data, err := ioutil.ReadFile("books.xml")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	err = xml.Unmarshal(data, &v)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	for _, book := range v.BooksList {
		fmt.Println(book.Author)
		fmt.Println(book.Id)
		fmt.Println(book.Price)
		fmt.Println(book.Title)
		fmt.Println(book.XMLName)
		fmt.Println(book.description)
	}
}
