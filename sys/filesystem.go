package sys

import (
	"io"
	"log"
	"os"
)

var location = "/home/amethystmac/Amethyst/"

func ReadFromDisk(fileName string) ([]byte) {
	file, err := os.Open(location + fileName)
	if err != nil {
		log.Panicln(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Panicln(err)
	}
	
	return data
}

func WriteToDisk(fileName string, data []byte) {
	file, err := os.Create(location + fileName)
	if err != nil {
		log.Panicln(err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		log.Panicln(err)
	}
}
