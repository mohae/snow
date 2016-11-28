package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	UA           = "snow"                                // UserAgent for snow
	URL          = "https://www.grc.com/securitynow.htm" // url of main security now page.
	SNURL        = "https://media.grc.com/sn/"
	concurrentDL = 1 // default number of episodes to download concurrently
	// if a value > maxConcurrency is specified, maxConcurrency will
	// be used and a message notifying the user will be emitted.
	maxConcurrentDL = 4 // maximum number of episodes to download concurrently
)

type Conf struct {
	lastN        int    // download the last n episodes. If 0, all are downloaded unless start is specified
	startEpisode int    // episode number to start downloading from; this takes precedence over lastN
	stopEpisode  int    // episode number to stop downloading at; if 0 everything up to current will be downloaded
	lowQuality   bool   // download the low quality version
	overwrite    bool   // overwrite existing file, if one exists
	ConcurrentDL int    `json:"concurrent_downloads"` // the number of episodes to download concurrently
	SaveDir      string `json:"save_dir"`             // directory to save the downloads to; if empty, $HOME/Downloads/security-now/ will be used
}

var (
	conf         Conf
	lastN        int
	startEpisode int
	stopEpisode  int
	concurrency  int
	lowQuality   bool
	overwrite    bool
	saveDir      string

	//verbose provides more detailed output
	verbose bool
)

func (c *Conf) Concurrency(i int) {
	if i == 0 {
		c.ConcurrentDL = concurrentDL
		fmt.Printf("info: invalid download concurrency, %d was specified, snow will use it's default value: %d\n", i, concurrentDL)
		return
	}
	if i > maxConcurrentDL {
		c.ConcurrentDL = maxConcurrentDL
		fmt.Printf("info: invalid download concurrency, %d was specified, snow will use it's maximum value: %d\n", i, maxConcurrentDL)
		return
	}
	c.ConcurrentDL = i
}

func init() {
	// -1 means last episode; the default
	flag.IntVar(&lastN, "lastn", 1, "download the last n episodes; 0 means all")
	flag.IntVar(&startEpisode, "start", 0, "episode number from which to start downloading")
	flag.IntVar(&stopEpisode, "stop", 0, "episode number at which to stop downloading")
	flag.IntVar(&concurrency, "concurrency", concurrentDL, "number of episodes to concurrently download")
	flag.BoolVar(&lowQuality, "lq", false, "download the low quality version: 16Kbps mp3")
	flag.BoolVar(&verbose, "verbose", false, "verbose output")
	flag.BoolVar(&overwrite, "overwrite", false, "overwrite existing file, if one exists")
	flag.StringVar(&saveDir, "savedir", "$HOME/Downloads/security-now", "save directory")
}

func main() {
	flag.Parse()

	conf.lastN = lastN
	conf.startEpisode = startEpisode
	conf.stopEpisode = stopEpisode
	conf.lowQuality = lowQuality
	conf.SaveDir = saveDir
	conf.Concurrency(concurrency) // set via method because the checking logic is part of conf

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
	mp3.Process()

	// TODO update completion messages
	fmt.Println("done")
}

// Verbose prints out messages if verbose.
func Verbose(s string) {
	if !verbose {
		return
	}
	fmt.Println(s)
}
