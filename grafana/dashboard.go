package grafana

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

type dashboard struct {
	ID        int      `json:"id"`
	UID       string   `json:"uid"`
	Title     string   `json:"title"`
	URI       string   `json:"uri"`
	URL       string   `json:"url"`
	Slug      string   `json:"slug"`
	Type      string   `json:"type"`
	Tags      []string `json:"tags"`
	IsStarred bool     `json:"isStarred"`
}

func (d dashboard) String() string {
	return fmt.Sprintf("{id: %d, uid: %s, title: %s}", d.ID, d.UID, d.Title)
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

func (d *DashboardFull) GetUrls(grafana *Grafana) (urls []fileUrl) {
	// urls := make([]fileUrl, 0)
	for _, dash := range d.Dashboard.Panels {
		urls = append(urls, dash.GetPanelIdWithGraph(grafana, d)...)
	}

	return urls

}
func parceTime(timeStr string) string {
	// log.Println("time:", timeStr)
	resultTime, err := time.ParseInLocation(timeFormat, timeStr, time.Local)
	if err != nil {
		log.Println(err)
	}
	result := strconv.FormatInt(resultTime.UnixNano()/1000000, 10)
	// log.Println("timeresult:", result)
	return result

}
