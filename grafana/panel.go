package grafana

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
)

type Panel struct {
	Id          int     `json:"id"`
	Title       string  `json:"Title"`
	Type        string  `json:"type"`
	Panels      []Panel `json:"panels"`
	Transparent bool    `json:"transparent"`
}

func (p Panel) String() string {
	return fmt.Sprintf("Panel: {id: %d, title: %s, type: %s, panels: %v, transparent: %v}", p.Id, p.Title, p.Type, p.Panels, p.Transparent)
}

func (p *Panel) GetPanelIdWithGraph(grafana *Grafana, dashboard *DashboardFull, down chan *FileUrl) {

	//если Тип является графиком то выполнить деествия в урпавляющей конструкции
	if p.Type == "graph" {
		// log.Println(p)
		//Парсим юрл Юрл Графаны
		var resultUrl = *grafana.URL.url

		resultUrl.Path = strings.ReplaceAll(dashboard.Meta.URL, "/d", "/render/d-solo")
		//формируем query для запроса
		qr := url.Values{}

		qr.Set("orgId", "1")
		qr.Set("panelId", strconv.Itoa(p.Id))

		qr.Set("from", parceTime(grafana.Test.TimeStart))
		qr.Set("to", parceTime(grafana.Test.TimeEnd))
		qr.Set("width", "1000")
		qr.Set("height", "500")
		qr.Set("tz", "Europe/Moscow")
		qr.Set("timeout", "20")

		if ls := grafana.Config[dashboard.Dashboard.Title]; len(ls) != 0 {

			for i, mapOfConfigs := range ls {
				log.Println("qr:", qr)
				qrWithConfig := qr

				for key, val := range mapOfConfigs {
					qrWithConfig.Add("var-"+key, val)
				}
				resultUrl.RawQuery = qrWithConfig.Encode()
				log.Println("qrWithConfig", qrWithConfig)

				file := FileUrl{}

				file.FileName = re.ReplaceAllString(fmt.Sprintf("%s_%d_%s", dashboard.Dashboard.Title, i, p.Title), "_")
				file.URL = &resultUrl

				down <- &file

			}
		} else {
			qrWithConfig := qr
			resultUrl.RawQuery = qrWithConfig.Encode()
			log.Println("qrWithConfig", qrWithConfig)

			file := FileUrl{}

			file.FileName = re.ReplaceAllString(fmt.Sprintf("%s_%s", dashboard.Dashboard.Title, p.Title), "_")
			file.URL = &resultUrl

			down <- &file

		}

	} else if len(p.Panels) != 0 {
		for _, panel := range p.Panels {
			panel.GetPanelIdWithGraph(grafana, dashboard, down)
		}

	}

}
