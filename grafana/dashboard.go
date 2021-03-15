package grafana

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

// Dashboard is struct from Grafana, using for serach after Search function.
type Dashboard struct {
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

func (d Dashboard) String() string {
	return fmt.Sprintf("dashboard: {id: %d, uid: %s, title: %s}", d.ID, d.UID, d.Title)
}

// Variables is a struct
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

// DashboardFull is a struct from Grafana Api, contains information about default variables of Grafana's dashboard and panels.
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
	return fmt.Sprintf("DashboardFull: {title: %s}", d.Dashboard.Title)
}

func parceTime(timeStr string) string {
	resultTime, err := time.ParseInLocation(timeFormat, timeStr, time.Local)
	if err != nil {
		log.Println(err)
	}
	result := strconv.FormatInt(resultTime.UnixNano()/1000000, 10)
	return result

}
