package main

import (
	"os"
	"log"
	"fmt"
	"regexp"
	"os/exec"
	"io/ioutil"
	"path/filepath"
	"github.com/jonasagx/csutils"
)

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
	if err != nil {
		log.Fatal(err)
	}
	return dir
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

func main() {
	url := csutils.ReadText("Youtube url: ")

	DownloadVideo(url)

	files := GetFilesList(GetCurrentDir())

	videos := FilterFilenames(files, `\.mkv|mp4|webm`)
	log.Println(videos)

	for _,video := range videos {
		Convert2Mp3(video)
	}
}