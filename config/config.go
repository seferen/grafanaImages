package config

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
)

//GetConfigFromFile return struct from a file with filename and write it to struct configStruct
func GetConfigFromFile(fileName string, configStruct interface{}) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Panicln(err)
	}
	defer file.Close()
	buff := bytes.Buffer{}
	_, err = buff.ReadFrom(file)
	if err != nil {
		log.Println(err)
	}
	json.Unmarshal(buff.Bytes(), configStruct)
	log.Println(configStruct)

}
