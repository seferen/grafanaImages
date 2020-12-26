package grafana

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Panel struct {
	Id     int     `json:"id"`
	Title  string  `json:"Title"`
	Type   string  `json:"type"`
	Panels []Panel `json:"panels"`
}

func (p Panel) String() string {
	return fmt.Sprintf("{id: %d, title: %s, type: %s, panels: %v}", p.Id, p.Title, p.Type, p.Panels)
}

func (p *Panel) GetPanelIdWithGraph(grafana *Grafana, dashboard *DashboardFull) []*url.URL {
	panelIDArray := make([]*url.URL, 0)

	//если Тип является графиком то выполнить деествия в урпавляющей конструкции
	if p.Type == "graph" {
		//Парсим юрл Юрл Графаны
		resultUrl, err := url.Parse(grafana.URL)
		if err != nil {
			log.Println("Неверный формат URL GRAFANA")

		}

		resultUrl.Path = strings.ReplaceAll(dashboard.Meta.URL, "/d", "/render/d-solo")
		//формируем query для запроса
		qr := resultUrl.Query()

		for _, v := range dashboard.Dashboard.Templating.List {
			log.Println(v)
			if v.UseTags == true {
				switch t := v.Current.Value.(type) {
				case []interface{}:
					log.Println("array", t)
				default:
					log.Println("notarray", t)

				}

			}

		}
		qr.Add("orgId", "1")
		qr.Add("panelId", strconv.Itoa(p.Id))

		qr.Add("from", parceTime(grafana.Test.TimeStart))
		qr.Add("to", parceTime(grafana.Test.TimeEnd))

		resultUrl.RawQuery = qr.Encode()

		log.Println(resultUrl.String())

		panelIDArray = append(panelIDArray, resultUrl)
	} else if len(p.Panels) != 0 {
		for _, panel := range p.Panels {
			test := panel.GetPanelIdWithGraph(grafana, dashboard)
			// log.Println(">>>>", test)
			panelIDArray = append(panelIDArray, test...)
		}

	}

	return panelIDArray

}

type Variables struct {
	Current struct {
		Value interface{} `json:"value"`
	} `json:"current"`
	Name    string `json:"name"`
	UseTags bool   `json:"useTags"`
}

func (v *Variables) UnmarshalJSON(b []byte) error {
	v.UseTags = true
	return json.Unmarshal(b, v)

}

type DashboardFull struct {
	Meta struct {
		Slug string `json:"slug"`
		URL  string `json:"url"`
	} `json:"meta"`
	Dashboard struct {
		Annotations   interface{} `json:"annotations"`
		Editable      bool        `json:"editable"`
		GnetId        int         `json:"gnetId"`
		GraphTooltip  int         `json:"graphTooltip"`
		Id            int         `json:"id"`
		Iteration     int         `json:"iteration"`
		Links         []string    `json:"links"`
		Panels        []Panel     `json:"panels"`
		Refresh       bool        `json:"refresh"`
		SchemaVersion int         `json:"schemaVersion"`
		Style         string      `json:"style"`
		Tags          []string    `json:"tags"`
		Templating    struct {
			List []Variables `json:"list"`
		} `json:"templating"`
		Time       interface{} `json:"time"`
		Timepicker interface{} `json:"timepicker"`
		Timezone   string      `json:"timezone"`
		Title      string      `json:"title"`
		Uid        string      `json:"uid"`
		Version    int         `json:"version"`
	} `json:"dashboard"`
}

func (d DashboardFull) String() string {
	return fmt.Sprintf("{title: %s}", d.Dashboard.Title)
}

func (d *DashboardFull) GetUrls(grafana *Grafana) {
	urls := make([]*url.URL, 0)
	for _, dash := range d.Dashboard.Panels {
		urls = append(urls, dash.GetPanelIdWithGraph(grafana, d)...)
	}
	for _, u := range urls {
		log.Println(u.String())
	}

}
func parceTime(timeStr string) string {
	resultTime, err := time.Parse(timeFormat, timeStr)
	if err != nil {
		log.Println(err)
	}
	return strconv.FormatInt(resultTime.UnixNano()/1000000, 10)

}
