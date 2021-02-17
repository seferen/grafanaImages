package config

import (
	"encoding/json"
	"log"
	"os"
)

//GetConfigFromFile return struct from a file with filename and write it to struct configStruct
func GetConfigFromFile(fileName string, configStruct interface{}) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	dec.Decode(configStruct)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("file", fileName, "was readed.", configStruct)

}
