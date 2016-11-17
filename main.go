package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	UA    = "snow"                                // UserAgent for snow
	URL   = "https://www.grc.com/securitynow.htm" // url of main security now page.
	SNURL = "https://media.grc.com/sn/"
)

type Conf struct {
	lastN        int    // download the last n episodes. If 0, all are downloaded unless start is specified
	startEpisode int    // episode number to start downloading from; this takes precedence over lastN
	stopEpisode  int    // episode number to stop downloading at; if 0 everything up to current will be downloaded
	lowQuality   bool   // download the low quality version
	SaveDir      string `json:"save_dir"` // directory to save the downloads to; if empty, $HOME/Downloads/security-now/ will be used
}

var conf Conf

func init() {
	flag.IntVar(&conf.lastN, "lastn", 0, "download the last n episodes; 0 means all")
	flag.IntVar(&conf.startEpisode, "start", 0, "episode number from which to start downloading")
	flag.IntVar(&conf.stopEpisode, "stop", 0, "episode number at which to stop downloading")
	flag.BoolVar(&conf.lowQuality, "lq", false, "download the low quality version: 16kbps mp3")
	flag.StringVar(&conf.SaveDir, "savedir", "$HOME/Downloads/security-now", "save directory")
}

func main() {
	flag.Parse()

	// check flags for validity
	if conf.SaveDir == "" {
		fmt.Println("must specify a save directory; to use the default do not use the -savedir flag")
		return
	}

	if conf.startEpisode > 0 && conf.stopEpisode < conf.startEpisode {
		fmt.Printf("episode at which to stop downloading, %d, must be either greater than the start episode, %d, or 0\n", conf.stopEpisode, conf.startEpisode)
		return
	}

	// resolve home dir
	conf.SaveDir = os.ExpandEnv(conf.SaveDir)

	// make the dir (if necessary)
	err := os.MkdirAll(conf.SaveDir, 764)
	if err != nil {
		fmt.Printf("error making save dir: %s\n", err)
		return
	}

	// check the latest episode number; this will be the limit
	i, err := GetLastEpisodeNumber()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	// if the last episode is 0, something went wrong.
	if i == 0 {
		fmt.Println("error: snow encountered an unknown problem while processing episode information, the last episode was 0")
		return
	}

	// set the Start Stop info
	err = setEpisodeRange(i, &conf)
	if err != nil {
		fmt.Println(err)
		return
	}

	// download
	// TODO: add concurrency
	mp3 := NewMP3(conf)
	cnt, n, err := mp3.Process()
	if err != nil {
		fmt.Printf("processing error: %s\n", err)
		fmt.Printf("%d of %d episodes were successfully downloaded\n", cnt, mp3.stopEpisode-mp3.startEpisode)
		fmt.Printf("%d bytes downloaded\n", n)
	}
	fmt.Printf("processing complete: %d episodes totalling %d bytes were downloaded\n", cnt, n)
}
