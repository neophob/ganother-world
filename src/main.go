package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	data, err := ioutil.ReadFile("./assets/memlist.bin")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	resourceMap := unmarshallingMemlistBin(data)
	printStatisticsForMemlistBin(resourceMap)
}
