package backend

import (
	"log"
	"regexp"
)

type Song struct {
	Title  string `json:title`
	Artist string `json:artist`
	Album  string `json:album`
	Url    string `json:url`
	Year   string `json:year`
}

func (s *Song) GetVideoId() string {
	regex := regexp.MustCompile("v=")

	if !regex.MatchString(s.Url) {
		log.Fatal("Invalid Url" + s.Url)
	}

	parts := regex.Split(s.Url, 2)
	return parts[1]
}