package backend

import (
	"os"
	"log"
	"path"
	"regexp"
	"os/exec"
	"path/filepath"
	"io/ioutil"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func DeleteFile(filePath string) {
	var err = os.Remove(filePath)
	checkErr(err)
}

func GetFilesList(dir string) []string {
	var filenames []string

	infoFiles, _ := ioutil.ReadDir(dir)

	for _, filename := range infoFiles {
		filenames = append(filenames, filename.Name())
	}

	return filenames
}

func FilterFilenames(filenames []string, expr string) []string {
	var matches []string
	checker := regexp.MustCompile(expr)

	for _, file := range filenames {
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
	return makePath(GetCurrentDir(), Conf.AudioOutputFolder, filename)
}

func GetVideoOutputDir(filename string) string {
	return makePath(GetCurrentDir(), Conf.VideoOutputFolder, filename)
}