package grafana

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
	URL  string `json:"url"`
	Test struct {
		TimeStart string `json:"timeStart"`
		TimeEnd   string `json:"timeEnd"`
	} `json:"test"`
	TOKEN      string                         `json:"token"`
	Config     map[string][]map[string]string `json:"dashboards"`
	dashboards []dashboard
}

func (g Grafana) String() string {
	return fmt.Sprintf("{url: %s, token: %s}", g.URL, g.TOKEN)
}

func (g *Grafana) Search() error {

	urlRes := g.URL + "/api/search/"
	log.Println(urlRes)
	req, err := g.NewGrafanaRequest(http.MethodGet, urlRes, nil)
	if err != nil {
		return err
	}
	// req.Header.Add("Authorization", "Bearer "+g.TOKEN)
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

	// log.Println("Config:", g.Config)
	// log.Println("Dashboards:", result)
	for _, dash := range result {
		if _, b := g.Config[dash.Title]; b {
			g.dashboards = append(g.dashboards, dash)
		}
	}
	log.Println("Result dashboards:", g.dashboards)

	return nil
}

func (g *Grafana) getDashboardByUid(uid string) (dash *DashboardFull, err error) {
	urlRes := strings.Trim(g.URL, "/") + "/api/dashboards/uid/" + uid
	log.Println(urlRes)
	req, err := g.NewGrafanaRequest(http.MethodGet, urlRes, nil)
	if err != nil {
		return dash, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return dash, err
	}
	defer resp.Body.Close()

	bodyB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return dash, err
	}

	err = json.Unmarshal(bodyB, &dash)
	if err != nil {
		return dash, err
	}
	// dash.GetPanelId()
	return dash, nil

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
