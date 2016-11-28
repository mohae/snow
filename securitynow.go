package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/net/html"
)

// download holds information about a given download
type Download struct {
	Name    string // the name of the thing downloaded
	Path    string // the path of the save file; including name
	skipped bool
	n       int64 // number of bytes downloaded
	err     error // error incountered, if any
}

func (d *Download) skip(s string) {
	d.skipped = true
	fmt.Printf("skipped %s: %s at %s\n", d.Name, s, d.Path)
}

// MP3 handles the downloading of MP3 episodes
type MP3 struct {
	concurrency int
	// start/stop are inclusive
	startEpisode int
	stopEpisode  int
	saveDir      string
	workCh       chan int      // channel for sending work to
	resultCh     chan Download // channel for sending result of download to
	downloads    []Download    // results of the downloads
	Get          func(i int) Download
}

// Returns a MP3 processor.
func NewMP3(c Conf) *MP3 {
	var mp3 MP3
	mp3.concurrency = c.ConcurrentDL
	mp3.startEpisode = c.startEpisode
	mp3.stopEpisode = c.stopEpisode
	mp3.saveDir = c.SaveDir
	mp3.workCh = make(chan int)
	mp3.resultCh = make(chan Download)
	if c.lowQuality {
		mp3.Get = mp3.LowQuality
		return &mp3
	}
	mp3.Get = mp3.HighQuality
	return &mp3
}

// Process processes the episodes to download. If an error is encountered,
// the number of successfully processed episodes, the total bytes downloaded
// and the error are returned. If the process completes without an error, the
// number of episodes downloaded along with the bytes downloaded are returned.
func (m *MP3) Process() {
	for i := 0; i < m.concurrency; i++ {
		go m.GetEpisodes()
	}

	fmt.Println("downloading...")

	go func() {
		for i := m.startEpisode; i <= m.stopEpisode; i++ {
			fmt.Println(i)
			m.workCh <- i
		}
	}()

	// we know how many results we're going to get so we just count the results
	for i := 0; i < (m.stopEpisode - m.startEpisode + 1); i++ {
		fmt.Printf("waiting for result %d\n", i+1)
		v := <-m.resultCh
		m.downloads = append(m.downloads, v)
		fmt.Println(v)
	}

	fmt.Println("complete...")

	return
}

// GetEpisodes downloads episodes.
func (m *MP3) GetEpisodes() {
	// work until work channel is closed
	for {
		fmt.Println("get episodes")
		i, ok := <-m.workCh
		fmt.Println(i, ok)
		if !ok {
			return
		}
		fmt.Println(i)
		m.resultCh <- m.Get(i)
		fmt.Println("result sent")
	}
}

// LowQuality downloads the high quality 16Kbps version of an episode.
func (m *MP3) LowQuality(i int) Download {
	var d Download
	d.Name = fmt.Sprintf("sn-%03d-lq.mp3", i)
	d.Path = filepath.Join(m.saveDir, d.Name)
	fmt.Println("download:" + d.Name)
	return m.Download(d)
}

// HighQuality downloads the high quality 64Kbps version of an episode.
func (m *MP3) HighQuality(i int) Download {
	var d Download
	d.Name = fmt.Sprintf("sn-%03d-lq.mp3", i)
	d.Path = filepath.Join(m.saveDir, d.Name)
	fmt.Println("download:" + d.Name)
	return m.Download(d)
}

// Download handles the actual download.
func (m *MP3) Download(d Download) Download {
	// if it already exists; don't do anything
	if fileExists(d.Path) {
		d.skip("file exists")
		return d
	}
	// open the save file
	f, err := os.OpenFile(d.Path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0664)
	if err != nil {
		d.err = err
		return d
	}
	defer f.Close()
	// Get the file
	resp, err := http.Get(SNURL + d.Name)
	if err != nil {
		d.err = err
		return d
	}
	defer resp.Body.Close()
	for {
		n, err := io.Copy(f, resp.Body)
		d.n += n
		if err != nil {
			if err == io.EOF {
				return d
			}
			d.err = err
			return d
		}
		// no bytes copied == done
		if n == 0 {
			break
		}
	}
	return d
}

// GetLastEpisodenNumber returns the number of the most recent episode. This
// represents the current upper bound.
func GetLastEpisodeNumber() (int, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("GET of %q resulted in an unexpected status: %q", URL, resp.Status)
	}
	// tokenize the response
	tokens := getTokens(resp.Body)
	i, err := lastEpisodeFromTokens(tokens)
	if err != nil {
		return 0, fmt.Errorf("%s: %s", URL, err)
	}
	return i, nil
}

// lastEpisodeFromTokens checks the tokens for the max episode by checking all
// anchor tags for a value that translates into an int returning the largest
// value found. The last episode is the first anchor with the key "name" that can be
// converted to an int, but this goes through all of them just in case.
func lastEpisodeFromTokens(tokens []html.Token) (int, error) {
	var i int
	for _, token := range tokens {
		if token.Type == html.StartTagToken && token.DataAtom.String() == "a" {
			for _, attr := range token.Attr {
				if attr.Key == "name" {
					v, err := strconv.Atoi(attr.Val)
					if err == nil && v > i {
						i = v
					}
				}
			}
		}
	}
	if i == 0 {
		return 0, errors.New("no episode numbers found")
	}
	return i, nil
}

// getTokens returns all tokens in the body
func getTokens(body io.Reader) []html.Token {
	var tokens []html.Token
	page := html.NewTokenizer(body)
	for {
		typ := page.Next()
		if typ == html.ErrorToken {
			return tokens
		}
		tokens = append(tokens, page.Token())
	}
}

func setEpisodeRange(i int, cnf *Conf) error {

	// if there's a startEpisode make sure it's within range
	if cnf.startEpisode > 0 {
		if cnf.startEpisode > i {
			return fmt.Errorf("Nothing to do: the start episode, %d, does not yet exist. The last episode was %d.", cnf.startEpisode, i)
		}

		if cnf.stopEpisode > i || cnf.stopEpisode == 0 {
			cnf.stopEpisode = i
		}

		return nil
	}

	// lastN processing means we'll always stop at current episode
	cnf.stopEpisode = i

	switch cnf.lastN {
	case 1: // -1 means last episode
		cnf.startEpisode = i
		return nil
	case 0: // all episodes
		cnf.startEpisode = 1
		return nil
	}
	// otherwise calculate n episodes ago
	cnf.startEpisode = i - cnf.lastN + 1
	// make sure it's within range
	if cnf.startEpisode < 0 {
		cnf.startEpisode = 1
	}
	return nil
}

func printDownloadMessage(episode int, n int64, name string) {
	fmt.Printf("downloaded episode %d, totalling %d bytes, as %s\n", episode, n, name)
}

// technically speaking this is racy, but if you're using snow and mucking with
// security now episodes in the target dir...well don't blame snow for what
// does or does not happen. If any error, other than IsNotExist occurs, a true
// will be returned; this may be incorrect handling, but this is what happens
// when only a bool is returned.
func fileExists(name string) bool {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
