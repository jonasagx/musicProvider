package main

import (
	"os"
	"log"
	"fmt"
	"path"
	"time"
	"regexp"
	"os/exec"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"path/filepath"

	// "github.com/jonasagx/id3tags"
	"github.com/julienschmidt/httprouter"
)

const (
	supportedVideosFormats = `\.mkv|\.mp4|\.webm`
	httpPort = ":8000"
	audioOutputFolder = "musicFiles"
	videoOutputFolder = "videoFiles"
	videoFilenameFormat = "%(title)s-%(id)s.%(ext)s"
)

type Song struct {
	Title string `json:title`
	Artist string `json:artist`
	Album string `json:album`
	Url string `json:url`
}

func (s *Song) GetVideoId() string {
	regex := MustCompile("v=")
	parts := regex.Split(s.Url, 2)

	if len(parts) == 2 {
		return parts[1]
	} else {
		return ""
	}
}

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
	envPath, err := exec.LookPath(command)
    if err != nil {
    	log.Fatal(command, "not in environment path")
    	return false, ""
    }
    
    log.Println(command, "is available at ", envPath)
    return true, envPath
}

func GetCurrentDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	checkErr(err)
	return dir
}

func makePath(absolutePath, specialDir, filename string) string {
	var builtPath = path.Join(absolutePath, specialDir)
	err := os.MkdirAll(builtPath, os.ModePerm)
	checkErr(err)

	if filename != "" {
		builtPath = path.Join(builtPath, filename)
	}

	return builtPath
}

func GetAudioOutputDir(filename string) string {
	return makePath(GetCurrentDir(), audioOutputFolder, filename)
}

func GetVideoOutputDir(filename string) string {
	return makePath(GetCurrentDir(), videoOutputFolder, filename)
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
	args := []string{url, "-o", GetVideoOutputDir("") + "/" + videoFilenameFormat}

	log.Println("Downloading video")
	runCommand(command, args)
}

func StringParserToMp3(videoFilename string) string {
	checker := regexp.MustCompile(`\.mkv|\.mp4|\.webm`)

	if ok := checker.MatchString(videoFilename); !ok {
		log.Fatal("File name without known extension " + videoFilename)
	}

	return checker.ReplaceAllLiteralString(videoFilename, ".mp3")
}

func Convert2Mp3(videoFilename string){
	log.Println("Converting", videoFilename)

	command := "ffmpeg"

	audioFilename := StringParserToMp3(videoFilename)
	audioOutputFilePath := GetAudioOutputDir(audioFilename)

	videoInputFilePath := GetVideoOutputDir(videoFilename)

	args := []string{"-i", videoInputFilePath, "-vn", "-acodec", "libmp3lame", "-ac", "2", "-ab", "160k", "-ar", "48000", audioOutputFilePath}
	log.Println("Converting to mp3")
	runCommand(command, args)
}

// func SetMp3Tags(song Song) {
// 	var audioFile id3tags.Mp3
// 	audioFile.

// }

func ConvertVideoToMp3(song Song) {
	t0 := time.Now()
	log.Println("Converting to mp3")
	DownloadVideo(song.Url)

	files := GetFilesList(GetVideoOutputDir(""))

	videos := FilterFilenames(files, supportedVideosFormats)
	log.Println(videos)

	for _,video := range videos {
		Convert2Mp3(video)
	}

	t1 := time.Now()

	log.Println("Request took %v", t1.Sub(t0))
}

// --------------------------------- HTTP Facede ---------------------------------

func postVideoToMp3(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	checkErr(err)

	var song Song
	err = json.Unmarshal(body, &song)
	checkErr(err)

	json.NewEncoder(rw).Encode(song)

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

// ----------------------------- start up calling -----------------------------

func main() {
	StartHTTPServer()
}