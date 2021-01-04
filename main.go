package main

import (
	"log"

	"./config"
	"./grafana"
)

func main() {

	graf := grafana.Grafana{}
	config.GetConfigFromFile("config.json", &graf)

	graf.Search()
	err := graf.GetImages()
	if err != nil {
		log.Println(err)
	}
}
