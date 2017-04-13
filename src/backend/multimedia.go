package backend

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/jonasagx/id3tags"
)

func ConvertVideoToMp3(song Song) {
	t0 := time.Now()
	log.Println("Converting to mp3")
	DownloadVideo(song.Url)

	files := GetFilesList(GetVideoOutputDir(""))

	videos := FilterFilenames(files, Conf.SupportedVideosFormats)
	log.Println(videos)

	for _, video := range videos {
		Convert2Mp3(video)
		// DeleteFile(video)
	}

	SetMp3Tags(song)

	t1 := time.Now()
	log.Println("Request took", t1.Sub(t0))
}

func DownloadVideo(url string) {
	command := "youtube-dl"
	args := []string{url, "-o", GetVideoOutputDir("") + "/" + Conf.VideoFilenameFormat}

	log.Println("Downloading video")
	runCommand(command, args)
}

func Convert2Mp3(videoFilename string) {
	log.Println("Converting", videoFilename)

	command := "ffmpeg"

	audioFilename := StringParserToMp3(videoFilename)
	audioOutputFilePath := GetAudioOutputDir(audioFilename)

	videoInputFilePath := GetVideoOutputDir(videoFilename)

	args := []string{"-i", videoInputFilePath, "-vn", "-acodec", "libmp3lame", "-ac", "2", "-ab", "160k", "-ar", "48000", audioOutputFilePath}
	log.Println("Converting to mp3")
	runCommand(command, args)
}

func SetMp3Tags(song Song) {
	var audioFile id3tags.Mp3
	listOfAudioFiles := FilterFilenames(GetFilesList(Conf.AudioOutputFolder), song.GetVideoId())
	log.Println(listOfAudioFiles)

	if listOfAudioFiles != nil && len(listOfAudioFiles) == 1 {
		filename := listOfAudioFiles[0]
		filePath := GetAudioOutputDir(filename)

		audioFile.FilePath = filePath
		audioFile.GetID3Tags()
		audioFile.Artist = song.Artist
		audioFile.Title = song.Title
		audioFile.Album = song.Album
		audioFile.Year = song.Year
		audioFile.SetID3Tags()
	}

}

func runCommand(command string, args []string) {
	log.Println("Running", command, args)
	err := exec.Command(command, args...).Run()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}

func StringParserToMp3(videoFilename string) string {
	checker := regexp.MustCompile(`\.mkv|\.mp4|\.webm`)

	if ok := checker.MatchString(videoFilename); !ok {
		log.Fatal("File name without known extension " + videoFilename)
	}

	return checker.ReplaceAllLiteralString(videoFilename, ".mp3")
}
