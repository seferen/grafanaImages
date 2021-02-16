package grafana

import (
	"log"
	"net/http"
	"os"
	"regexp"
)

var client = &http.Client{
	Transport: http.DefaultTransport,
}

const timeFormat string = "2006-01-02 15:04:05" //Mon Jan 2 15:04:05 MST 2006
var re = regexp.MustCompile(`[\W]+`)
var dir string = "result/"

func init() {
	err := os.MkdirAll(dir, os.ModeAppend)

	if err != nil {
		log.Panicln("Dir did't create. Please try the path:", dir)
	}
	log.Println("Directory was tryed. Path:", dir)
}
