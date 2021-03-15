package grafana

import (
	"net/http"
	"regexp"
)

// http client for sendig request to grafana
var client = &http.Client{
	Transport: http.DefaultTransport,
}

// timeFormat is const of timeFormat for parsing from a string
const timeFormat string = "2006-01-02 15:04:05" //Mon Jan 2 15:04:05 MST 2006
var re = regexp.MustCompile(`[\W]+`)

// Dir - the variable for creating a directory
var Dir *string

// Prefix - the variable of string for adding to prefix of a file name.
var Prefix *string
