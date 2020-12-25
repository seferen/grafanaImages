package grafana

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

var client = &http.Client{}

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

//Grafana config of Grafana structure
type Grafana struct {
	URL        string                         `json:"url"`
	TOKEN      string                         `json:"token"`
	Config     map[string][]map[string]string `json:"dashboards"`
	dashboards []dashboard
}

func (g *Grafana) Search() error {

	urlRes := g.URL + "/api/search/"
	log.Println(urlRes)
	req, err := http.NewRequest(http.MethodGet, urlRes, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+g.TOKEN)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&g.dashboards)
	// bite, _ := ioutil.ReadAll(resp.Body)
	// log.Println(string(bite))
	// err = json.Unmarshal(bite, &g.dashboards)
	if err != nil {
		log.Println(err)
	}
	log.Println(g.dashboards)

	return nil
}

func (g *Grafana) FindByTitle(title string) {
	for _, dash := range g.dashboards {
		if strings.ToLower(title) == strings.ToLower(dash.Title) {
			log.Println(dash)
		}

	}
}
