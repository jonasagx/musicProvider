package backend

type Config struct {
	Version string
	HttpPort string
	AudioOutputFolder string
	VideoOutputFolder string
	VideoFilenameFormat	string
	SupportedVideosFormats string
}

func LoadConfig() *Config {
	return &Config {
		Version: "0.0.1",
		HttpPort: "8000",
		AudioOutputFolder: "musicFiles",
		VideoOutputFolder: "videoFiles",
		VideoFilenameFormat: "%(id)s.%(ext)s",
		SupportedVideosFormats: `\.mkv|\.mp4|\.webm`}
}

var Conf *Config