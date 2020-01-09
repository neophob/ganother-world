package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	data := readFile("./assets/memlist.bin")
	resourceMap := unmarshallingMemlistBin(data)
	printStatisticsForMemlistBin(resourceMap)

	loadEntryFromBank(resourceMap, 21)
}

func readFile(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("File reading error", err)
		os.Exit(1)
	}
	return data
}
