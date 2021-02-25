package main

import (
	"flag"
	"log"
	"os"

	"./config"
	"./grafana"
)

var (
	configFile *string
	prefix     *string
)

func main() {
	downChan := make(chan *grafana.FileUrl)
	graf := grafana.Grafana{}
	config.GetConfigFromFile(*configFile, &graf)

	go func() {
		log.Println("chan was started")
		for {
			select {
			case x := <-downChan:
				log.Println(">>>>>", x)
				graf.DownloadFile(x)

			}
		}
	}()

	dashs, err := graf.Search()
	if err != nil {
		log.Panicln(err)
	}

	for _, dash := range dashs {
		if _, b := graf.Config[dash.Title]; b {
			log.Println(dash)
			graf.Dashboards = append(graf.Dashboards, dash)
			dashF, err := graf.GetDashboardByUid(dash.UID)
			if err != nil {
				log.Panicln(err)
			}

			for _, dash := range dashF.Dashboard.Panels {
				dash.GetPanelIdWithGraph(&graf, dashF, downChan)
			}

		}
	}
	close(downChan)
	log.Println("Channel was closed")

}

func init() {
	log.Println("Main init")
	configFile = flag.String("f", "config.json", "a file with configeration for app")
	grafana.Dir = flag.String("dir", "result", "configure the path of a result directory where will download files of grapthics")
	grafana.Prefix = flag.String("prefix", "", "prefix for downloading files")
	flag.Parse()

	err := os.MkdirAll(*grafana.Dir, os.ModeAppend)

	if err != nil {
		log.Panicln("Dir did't create. Please try the path:", *grafana.Dir)
	}
	log.Println("Directory was tryed. Path:", *grafana.Dir)

}
