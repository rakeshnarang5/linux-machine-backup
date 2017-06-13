// package main

// import (
// 	"encoding/xml"
// 	"fmt"
// 	"io/ioutil"
// 	//"os"
// )

// type Data struct {
// 	Series      Show
// 	EpisodeList []Episode `xml:"Episode>"`
// }

// type Show struct {
// 	Title    string `xml:"SeriesName>"`
// 	SeriesID int    `xml:"SeriesID>"`
// }

// type Episode struct {
// 	SeasonNumber  int    `xml:"SeasonNumber>"`
// 	EpisodeNumber int    `xml:"EpisodeNumber>"`
// 	EpisodeName   string `xml:"EpisodeName>"`
// 	FirstAired    string `xml:"FirstAired>"`
// }

// func (s Show) String() string {
// 	return fmt.Sprintf("%s - %d", s.Title, s.SeriesID)
// }

// func (e Episode) String() string {
// 	return fmt.Sprintf("S%02dE%02d - %s - %s", e.SeasonNumber, e.EpisodeNumber, e.EpisodeName, e.FirstAired)
// }

// func main() {
// 	/*xmlFile, err := os.Open("data.xml")
// 	if err != nil {
// 		fmt.Printf("error: %v\n", err)
// 		return
// 	}
// 	defer xmlFile.Close()

// 	var q Data

// 	file, err := ioutil.ReadFile("data.xml")
// 	if err != nil {
// 		fmt.Printf("error: %v\n", err)
// 		return
// 	}

// 	fmt.Println(string(file))

// 	xml.Unmarshal(file, &q)

// 	fmt.Println(q.Series)
// 	for _, episode := range q.EpisodeList {
// 		fmt.Printf("\t%s\n", episode)
// 	}
// }
// */