package grafana

import (
	"encoding/json"
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
