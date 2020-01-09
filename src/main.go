package main

import (
	"fmt"
	"os"
	"io/ioutil"
)

func main() {
	data := readFile("./assets/memlist.bin")
	resourceMap := unmarshallingMemlistBin(data)
	printStatisticsForMemlistBin(resourceMap)
}

func readFile(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("File reading error", err)
		os.Exit(1)
	}
	return data
}