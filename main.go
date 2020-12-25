package main

import (
	"log"

	"./config"
	"./grafana"
)

func main() {

	graf := grafana.Grafana{}
	config.GetConfigFromFile("config.json", &graf)
	log.Println(graf)
	graf.Search()
	graf.FindByTitle("gatling")
}
