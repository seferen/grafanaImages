package grafana

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var client = &http.Client{}

const timeFormat string = "2006-01-02 15:04:05" //Mon Jan 2 15:04:05 MST 2006

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

//Grafana config of Grafana structure
type Grafana struct {
	URL  Url `json:"url"`
	Test struct {
		TimeStart string `json:"timeStart"`
		TimeEnd   string `json:"timeEnd"`
	} `json:"test"`
	TOKEN      string                         `json:"token"`
	Config     map[string][]map[string]string `json:"dashboards"`
	dashboards []dashboard
}

func (g Grafana) String() string {

	return fmt.Sprintf("{url: %s, token: %s}", g.URL.UrlStr, g.TOKEN)
}

func (g *Grafana) Search() error {

	urlRes := g.URL.UrlStr + "/api/search/"
	log.Println(urlRes)
	req, err := g.NewGrafanaRequest(http.MethodGet, urlRes, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	result := make([]dashboard, 10)
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&result)
	if err != nil {
		log.Println(err)
	}

	for _, dash := range result {
		if _, b := g.Config[dash.Title]; b {
			g.dashboards = append(g.dashboards, dash)
		}
	}
	log.Println("Result dashboards:", g.dashboards)

	return nil
}

func (g *Grafana) getDashboardByUid(uid string) (*DashboardFull, error) {
	urlRes := strings.Trim(g.URL.UrlStr, "/") + "/api/dashboards/uid/" + uid
	log.Println(urlRes)
	dash := DashboardFull{}
	req, err := g.NewGrafanaRequest(http.MethodGet, urlRes, nil)
	if err != nil {
		return nil, err
	}
	log.Println("Request:", urlRes, "was fomed")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	log.Println("Request:", urlRes, "was recived")
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	err = dec.Decode(&dash)
	if err != nil {
		return nil, err
	}
	log.Println("Request:", urlRes, "was write to dash as JSON")
	return &dash, nil

}

func (g *Grafana) GetImages() error {
	for _, dash := range g.dashboards {

		if dashboard, err := g.getDashboardByUid(dash.UID); err != nil {
			log.Println(err)
		} else {
			dashboard.GetUrls(g)
		}
	}
	return nil
}

func (g *Grafana) NewGrafanaRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", "Bearer "+g.TOKEN)
	return req, err

}
