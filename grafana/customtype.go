package grafana

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Url struct {
	UrlStr string
	url    *url.URL
}

func (u *Url) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &u.UrlStr)
	if err != nil {
		return err
	}

	u.url, err = url.Parse(u.UrlStr)
	if err != nil {
		return err
	}

	return nil

}

type fileUrl struct {
	FileName    string
	URL         *url.URL
	respStatus  int
	fileWriting bool
}

func (f fileUrl) String() string {
	return fmt.Sprintf("fileUrl: {fileName: %s, url: %s, Responce Status: %d, Writing file: %t}", f.FileName, f.URL.String(), f.respStatus, f.fileWriting)
}
