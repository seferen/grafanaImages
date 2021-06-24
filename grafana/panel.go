package grafana

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
)

// Panel struct contains information about Grafana's Panel, using for search graphs
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

// GetPanelIdWithGraph get search from deep panels and and create url and a file name to FileUrl struct and send resul to down chan
func (p *Panel) GetPanelIdWithGraph(grafana *Grafana, dashboard *DashboardFull, down chan *FileUrl) {

	//если Тип является графиком то выполнить деествия в урпавляющей конструкции
	if p.Type == "graph" {

		if ls := grafana.Config[dashboard.Dashboard.Title]; len(ls) != 0 {

			for i, mapOfConfigs := range ls {
				var resultUrl = *grafana.URL.url

				resultUrl.Path = strings.ReplaceAll(dashboard.Meta.URL, "/d", "/render/d-solo")
				qr := url.Values{}

				qr.Set("orgId", strconv.Itoa(grafana.Org.Id))
				qr.Set("panelId", strconv.Itoa(p.Id))

				qr.Set("from", parceTime(grafana.Test.TimeStart))
				qr.Set("to", parceTime(grafana.Test.TimeEnd))
				qr.Set("width", "1000")
				qr.Set("height", "500")
				qr.Set("tz", "Europe/Moscow")
				qr.Set("timeout", "20")

				// qrWithConfig := url.Values{}
				// qrWithConfig = qr

				for key, val := range mapOfConfigs {
					qr.Add("var-"+key, val)
				}
				log.Println("qr:", qr)

				resultUrl.RawQuery = qr.Encode()
				log.Println("qrWithConfig", qr)

				file := FileUrl{}

				file.FileName = re.ReplaceAllString(fmt.Sprintf("%s_%d_%s", dashboard.Dashboard.Title, i, p.Title), "_")
				file.URL = &resultUrl

				down <- &file

			}
		} else {
			var resultUrl = *grafana.URL.url

			resultUrl.Path = strings.ReplaceAll(dashboard.Meta.URL, "/d", "/render/d-solo")
			qr := url.Values{}

			qr.Set("orgId", "1")
			qr.Set("panelId", strconv.Itoa(p.Id))

			qr.Set("from", parceTime(grafana.Test.TimeStart))
			qr.Set("to", parceTime(grafana.Test.TimeEnd))
			qr.Set("width", "1000")
			qr.Set("height", "500")
			qr.Set("tz", "Europe/Moscow")
			qr.Set("timeout", "20")
			resultUrl.RawQuery = qr.Encode()
			log.Println("qrWithConfig", qr)

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
