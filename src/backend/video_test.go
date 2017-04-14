package backend

import (
	"log"
	"testing"
)

func TestExtractVideoId(t *testing.T) {
	text := "2017/04/13 11:30:02 [youtube] 75Ju0eM5T2c: Downloading webpage\n" +
	"[youtube] 75Ju0eM5T2c: Downloading video info webpage\n" + 
	"[youtube] 75Ju0eM5T2c: Extracting video information\n" +
	"[youtube] 75Ju0eM5T2c: Downloading MPD manifest\n" +
	"[download] Destination: /home/jonas/projects/musicProvider/bin/videoTmp/75Ju0eM5T2c.f243.webm\n" + 
	"[download] 100% of 7.69MiB in 00:0235MiB/s ETA 00:006\n" + 
	"[download] Destination: /home/jonas/projects/musicProvider/bin/videoTmp/75Ju0eM5T2c.f140.m4a\n" + 
	"[download] 100% of 3.24MiB in 00:0038MiB/s ETA 00:006\n" +
	"[ffmpeg] Merging formats into \"/home/jonas/projects/musicProvider/bin/videoTmp/75Ju0eM5T2c.mkv\"\n" + 
	"Deleting original file /home/jonas/projects/musicProvider/bin/videoTmp/75Ju0eM5T2c.f243.webm (pass -k to keep)\n" + 
	"Deleting original file /home/jonas/projects/musicProvider/bin/videoTmp/75Ju0eM5T2c.f140.m4a (pass -k to keep)\n"

	id := extractVideoId(text)
	log.Println(id)
}