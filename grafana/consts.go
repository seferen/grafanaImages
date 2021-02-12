package grafana

import (
	"net/http"
	"regexp"
)

var client = &http.Client{
	Transport: http.DefaultTransport,
}

const timeFormat string = "2006-01-02 15:04:05" //Mon Jan 2 15:04:05 MST 2006
var re = regexp.MustCompile(`[\W]+`)
