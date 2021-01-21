package grafana

import "net/http"

var client = &http.Client{}

const timeFormat string = "2006-01-02 15:04:05" //Mon Jan 2 15:04:05 MST 2006
