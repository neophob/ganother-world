package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	log.Println("- load memlist.bin")
	data := readFile("./assets/memlist.bin")
	resourceMap := unmarshallingMemlistBin(data)
	printStatisticsForMemlistBin(resourceMap)

	bankFilesMap := createBankMap("./assets/")
	loadEntryFromBank(resourceMap, bankFilesMap, 21)
}

func readFile(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("File reading error", err)
		os.Exit(1)
	}
	return data
}

func createBankMap(assetPath string) map[int][]byte {
	bankFilesMap := make(map[int][]byte)
	for i := 0x01; i < 0x0e; i++ {
		name := fmt.Sprintf("%sbank0%x", assetPath, i)
		log.Println("- load file", name)
		entry := readFile(name)
		bankFilesMap[i] = entry
	}
	return bankFilesMap
}
