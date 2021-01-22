package grafana

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

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

	// urlRes := g.URL.UrlStr + "/api/search/"
	// log.Println(urlRes)
	req, err := g.NewGrafanaRequest(http.MethodGet, g.URL.url.String()+"/api/search/", nil)
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

	dash := DashboardFull{}
	req, err := g.NewGrafanaRequest(http.MethodGet, g.URL.url.String()+"/api/dashboards/uid/"+uid, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	log.Println("Request:", req, "was recived")

	dec := json.NewDecoder(resp.Body)

	err = dec.Decode(&dash)
	if err != nil {
		return nil, err
	}
	log.Println("Request:", req, "was write to dash as JSON")
	return &dash, nil

}

func (g *Grafana) GetImages() error {
	for _, dash := range g.dashboards {
		log.Println("Dashboard:", dash)

		if dashboard, err := g.getDashboardByUid(dash.UID); err != nil {
			log.Println(err)
		} else {
			urls := dashboard.GetUrls(g)
			for _, u := range urls {
				log.Println(u.String())
				req, err := g.NewGrafanaRequest(http.MethodGet, u.URL.String(), nil)
				if err != nil {
					log.Println(err)
				}
				resp, err := client.Do(req)
				if err != nil {
					log.Println(err)
				}
				defer resp.Body.Close()
				u.respStatus = resp.StatusCode

				if resp.StatusCode == 200 {
					forFile, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						log.Println(err)
					}
					err = ioutil.WriteFile(u.FileName, forFile, os.ModeAppend)
					if err != nil {
						log.Println(err)
					}
					u.fileWriting = true
				}

				log.Println(u.String())

			}
		}
	}
	return nil
}

func (g *Grafana) NewGrafanaRequest(method, Url string, body io.Reader) (*http.Request, error) {

	log.Println("Url:", Url)
	req, err := http.NewRequest(method, Url, body)
	req.Header.Add("Authorization", "Bearer "+g.TOKEN)

	return req, err

}
