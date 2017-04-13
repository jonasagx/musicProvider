package backend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func postVideo(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	checkErr(err)

	var s Song
	err = json.Unmarshal(body, &s)
	checkErr(err)

	json.NewEncoder(rw).Encode(s)

	go DownloadVideo(s.Url)
	// go ConvertVideoToMp3(s)
}

func getHomePage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "Welcome to the MusicProvider!")
	fmt.Println("Endpoint Hit: homePage")
}

func StartHTTPServer() {
	log.Println("Starting server at http://localhost:" + Conf.HttpPort)

	router := httprouter.New()
	router.GET("/", getHomePage)
	router.POST("/download", postVideo)

	log.Fatal(http.ListenAndServe(":"+Conf.HttpPort, router))
}
