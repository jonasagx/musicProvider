package main

import (
	"os"
	"log"
	"fmt"
	"regexp"
	"os/exec"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"path/filepath"

	// "github.com/jonasagx/id3tags"
	"github.com/julienschmidt/httprouter"
)

var(
	supportedVideosFormats = `\.mkv|\.mp4|\.webm`
	httpPort = ":8000"
	audioOutputFolder = "musicFiles"
)

func checkErr(err error){
	if err != nil {
		panic(err)
	}
}

func GetFilesList(dir string) []string {
	var filenames []string

	infoFiles, _ := ioutil.ReadDir(dir)
    
    for _,filename := range infoFiles {
    	filenames = append(filenames, filename.Name())
    }

    return filenames
}

func FilterFilenames(filenames []string, expr string) []string {
	var matches []string
	checker := regexp.MustCompile(expr)

	for _,file := range filenames {
		if checker.MatchString(file) {
			matches = append(matches, file)
		}
	}

	return matches
}

func CheckPath(command string) (bool, string) {
	path, err := exec.LookPath(command)
    if err != nil {
    	log.Fatal(command, "not in path")
    	return false, ""
    }
    
    log.Println(command, "is available at ", path)
    return true, path
}

func GetCurrentDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	checkErr(err)
	return dir
}

func GetAudioOutputDir() string {
	var currentDir = GetCurrentDir()
	var path = path.Join(currentDir, audioOutputFolder)

	_, err := os.Stat(path)
	checkErr(err)
	err = os.MkdirAll(path, os.ModePerm)
	checkErr(err)
}

func runCommand(command string, args []string){
	log.Println("Running", command, args)
	err := exec.Command(command, args...).Run()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}

func DownloadVideo(url string) {
	command := "youtube-dl"
	args := []string{url}

	log.Println("Downloading video")

	runCommand(command, args)
}

func StringParserToMp3(videoFilename string) string {
	checker := regexp.MustCompile(`mkv|mp4|webm`)

	if ok := checker.MatchString(videoFilename); !ok {
		log.Fatal("File name without known extension " + videoFilename)
	}

	return checker.ReplaceAllLiteralString(videoFilename, "mp3")
}

func Convert2Mp3(filenameInput string){
	command := "ffmpeg"

	audioFilename := StringParserToMp3(filenameInput)
	args := []string{"-i", filenameInput, "-vn", "-acodec", "libmp3lame", "-ac", "2", "-ab", "160k", "-ar", "48000", audioFilename}
	log.Println("Converting to mp3")
	runCommand(command, args)
}

// func SetMp3Tags(song Song) {
// 	var audioFile id3tags.Mp3
// 	audioFile.

// }

func ConvertVideoToMp3(song Song){
	log.Println("Processing", song)
	DownloadVideo(song.Url)

	files := GetFilesList(GetCurrentDir())

	videos := FilterFilenames(files, supportedVideosFormats)
	log.Println(videos)

	for _,video := range videos {
		Convert2Mp3(video)
	}
}

// --------------------------------- HTTP Facede ---------------------------------

type Song struct {
	Title string `json:title`
	Artist string `json:artist`
	Album string `json:album`
	Url string `json:url`
}

func postVideoToMp3(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	checkErr(err)

	var song Song
	err = json.Unmarshal(body, &song)
	checkErr(err)

	ConvertVideoToMp3(song)
}

func getHomePage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "Welcome to the MusicProvider!")
	fmt.Println("Endpoint Hit: homePage")
}

func StartHTTPServer() {
	log.Println("Starting server at http://localhost" + httpPort)
	
	router := httprouter.New()
	router.GET("/", getHomePage)
	router.POST("/", postVideoToMp3)

	log.Fatal(http.ListenAndServe(httpPort, router))
}

func main() {
	StartHTTPServer()
}