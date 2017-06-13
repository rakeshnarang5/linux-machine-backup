package main

import (
	"fmt"
	"strings"
)

func main() {
	str := "(((akjdlfa)))&&(((adsfadfasdfadfa)))&&(((adsfadfadfadsadsfadfa)))&&"
	retVal := splitter(str)
	fmt.Println(retVal)
}

func splitter(str string) []string {
	if len(str) == 0 {
		retVal := make([]string, 0, 0)
		return retVal
	}
	n := 0
	i := 0
	for str[i] == 40 {
		n++
		i++
	}
	fmt.Println(n, i)
	xyz := createString(n)
	fmt.Println(xyz)
	idx := strings.Index(str, xyz)
	fmt.Println(idx)
	retVal := splitter(str[(idx + n + 2):])
	value := str[1:(idx + 2)]
	value = value
	fmt.Println(value)
	retVal = append(retVal, value)
	return retVal

}

func createString(n int) string {
	retVal := ""
	for n > 0 {
		retVal += ")"
		n--
	}
	return retVal
}
