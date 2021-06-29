package grafana

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

//Grafana config of Grafana structure hold all configurations needed for a work of the application.
type Grafana struct {
	URL  Url `json:"url"`
	Org  Org
	Test struct {
		TimeStart string `json:"timeStart"`
		TimeEnd   string `json:"timeEnd"`
	} `json:"test"`
	TOKEN      string                         `json:"token"`
	Config     map[string][]map[string]string `json:"dashboards"`
	Dashboards []Dashboard
}

func (g Grafana) String() string {

	return fmt.Sprintf("Grafana: {url: %s, token: %s}", g.URL.UrlStr, g.TOKEN)
}

func (g *Grafana) Ping() error {

	req, err := g.NewGrafanaRequest(http.MethodGet, g.URL.url.String()+"/api/login/ping/", nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("Please check your grafana url: " + g.URL.String() + ". It was return status code: " + resp.Status)
	}

	return nil

}

//Search - function that used the API Grafana from it documentation and insert into the Grafana struct dashboards with information about elements
func (g *Grafana) Search() ([]Dashboard, error) {

	req, err := g.NewGrafanaRequest(http.MethodGet, g.URL.url.String()+"/api/search/", nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result := make([]Dashboard, 10)
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&result)
	if err != nil {
		log.Println(err)
	}

	return result, nil
}

func (g *Grafana) GetOrgId() error {
	req, err := g.NewGrafanaRequest(http.MethodGet, g.URL.url.String()+"/api/org/", nil)

	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&g.Org)
	if err != nil {
		return err
	}

	return nil
}

func (g *Grafana) GetDashboardByUid(uid string) (*DashboardFull, error) {

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

	log.Println("Request:", req.RequestURI, "was recived")

	dec := json.NewDecoder(resp.Body)

	err = dec.Decode(&dash)
	if err != nil {
		return nil, err
	}
	log.Println("Request:", req.RequestURI, "was write to dash as JSON")
	return &dash, nil

}

func (g *Grafana) NewGrafanaRequest(method, Url string, body io.Reader) (*http.Request, error) {

	log.Println("Url:", Url)
	req, err := http.NewRequest(method, Url, body)
	req.Header.Add("Authorization", "Bearer "+g.TOKEN)

	return req, err

}

// DownloadFile function for downloading of a file from url from information about FileUrl
func (g *Grafana) DownloadFile(u *FileUrl) {
	log.Println("START DOWNLOAD:", u.String())
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
		writeFile(forFile, u, ".png")

	}

	log.Println("STOP DOWNLOAD:", u.String())

}

func writeFile(fileData []byte, u *FileUrl, endFile string) {
	for i := 0; true; i++ {

		filePath := filepath.Join(*Dir, fmt.Sprintf("%s%s_%d%s", *Prefix, u.FileName, i, endFile))
		log.Println("file will writed to", filePath)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			if err := os.WriteFile(filePath, fileData, os.ModePerm); err != nil {
				log.Println("Problem with writing file:", u.FileName, "message:", err)
			} else {
				u.fileWriting = true
			}
			break
		} else {
			log.Println("file is exist:", filePath)
		}

	}

}
