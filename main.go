package main

import (
	"log"

	"./config"
	"./grafana"
)

func main() {
	downChan := make(chan *grafana.FileUrl)
	graf := grafana.Grafana{}
	config.GetConfigFromFile("config.json", &graf)

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
