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
	Id          int     `json:"id"`
	Title       string  `json:"Title"`
	Type        string  `json:"type"`
	Panels      []Panel `json:"panels"`
	Transparent bool    `json:"transparent"`
}

func (p Panel) String() string {
	return fmt.Sprintf("{id: %d, title: %s, type: %s, panels: %v, transparent: %v}", p.Id, p.Title, p.Type, p.Panels, p.Transparent)
}

func (p *Panel) GetPanelIdWithGraph(grafana *Grafana, dashboard *DashboardFull) []*url.URL {
	panelIDArray := make([]*url.URL, 0)

	//если Тип является графиком то выполнить деествия в урпавляющей конструкции
	if p.Type == "graph" {
		log.Println(p)
		//Парсим юрл Юрл Графаны
		resultUrl, err := url.Parse(grafana.URL)
		if err != nil {
			log.Println("Неверный формат URL GRAFANA")

		}

		resultUrl.Path = strings.ReplaceAll(dashboard.Meta.URL, "/d", "/render/d-solo")
		//формируем query для запроса
		qr := resultUrl.Query()

		// for _, v := range dashboard.Dashboard.Templating.List {
		// 	// if v.UseTags == true {
		// 	switch t := v.Current.Value.(type) {
		// 	case []interface{}:
		// 		log.Println("tag struct:", v, " array", t)

		// 		for _, v1 := range t {
		// 			qr.Add("var-"+v.Name, strings.Trim(v1.(string), "$__"))

		// 		}
		// 	case string:
		// 		log.Println("tag struct:", v, " notarray", t)
		// 		qr.Add("var-"+v.Name, strings.Trim(t, "$__"))
		// 	default:
		// 		log.Println("ERROR", t)

		// 	}
		// 	// }
		// }

		qr.Set("orgId", "1")
		qr.Set("panelId", strconv.Itoa(p.Id))

		qr.Set("from", parceTime(grafana.Test.TimeStart))
		qr.Set("to", parceTime(grafana.Test.TimeEnd))
		qr.Set("width", "1000")
		qr.Set("height", "500")
		qr.Set("tz", "Europe/Moscow")
		qr.Set("timeout", "20")

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
	type aliasVariables Variables
	var all aliasVariables
	all.UseTags = true
	err := json.Unmarshal(b, &all)
	*v = Variables(all)
	return err

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
	// for _, u := range urls {
	// 	log.Println(u.String())
	// }

}
func parceTime(timeStr string) string {
	resultTime, err := time.Parse(timeFormat, timeStr)
	if err != nil {
		log.Println(err)
	}
	return strconv.FormatInt(resultTime.UnixNano()/1000000, 10)

}
