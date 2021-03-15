package grafana

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// URl ia a struct contains information about Grafana Url after parcing.
type Url struct {
	UrlStr string
	url    *url.URL
}

func (u Url) String() string {
	return u.UrlStr

}

func (u *Url) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.UrlStr)
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

// FileUrl is a struct for saving information about downloading file
type FileUrl struct {
	// FileName is a variable that contains information about a downloading name of file
	FileName string
	// URL ia a variable for downloading a file with FileName
	URL *url.URL
	// respStatus is a variable contains information about a status of respoce after downloading file
	respStatus int
	// fileWriting is an information about success writing file with FileName
	fileWriting bool
}

func (f FileUrl) String() string {
	return fmt.Sprintf("fileUrl: {fileName: %s, url: %s, Responce Status: %d, Writing file: %t}", f.FileName, f.URL.String(), f.respStatus, f.fileWriting)
}
