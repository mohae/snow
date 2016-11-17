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

// MP3 handles the downloading of MP3 episodes
type MP3 struct {
	startEpisode int
	stopEpisode  int
	saveDir      string
	Get          func(i int) (int64, error)
}

// Returns a MP3 processor.
func NewMP3(c Conf) *MP3 {
	var mp3 MP3
	mp3.startEpisode = c.startEpisode
	mp3.stopEpisode = c.stopEpisode
	mp3.saveDir = c.SaveDir
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
func (m *MP3) Process() (cnt int, bytes int64, err error) {
	for i := m.startEpisode; i <= m.stopEpisode; i++ {
		n, err := m.Get(i)
		bytes += n
		if err != nil {
			return cnt, bytes, err
		}
		cnt++
	}
	return cnt, bytes, nil
}

// LowQuality downloads the high quality 16Kbps version of an episode.
func (m *MP3) LowQuality(i int) (int64, error) {
	file := fmt.Sprintf("sn-%03d-lq.mp3", i)
	dest := filepath.Join(m.saveDir, file)
	// open the save file
	f, err := os.OpenFile(dest, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0664)
	if err != nil {
		return 0, fmt.Errorf("error processing %s: %s", dest, err)
	}
	defer f.Close()
	// Get the file
	resp, err := http.Get(SNURL + file)
	if err != nil {
		return 0, fmt.Errorf("error processing %s: %s", file, err)
	}
	defer resp.Body.Close()
	var b int64
	for {
		n, err := io.Copy(f, resp.Body)
		b += n
		if err != nil {
			if err == io.EOF {
				return b, nil
			}
			return b, fmt.Errorf("error processing %s: %s", file, err)
		}
		// no bytes copied == done
		if n == 0 {
			return b, nil
		}
	}
}

// HighQuality downloads the high quality 64Kbps version of an episode.
func (m *MP3) HighQuality(i int) (int64, error) {
	file := fmt.Sprintf("sn-%03d.mp3", i)
	dest := filepath.Join(m.saveDir, file)
	// open the save file
	f, err := os.OpenFile(dest, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0664)
	if err != nil {
		return 0, fmt.Errorf("error processing %s: %s", dest, err)
	}
	defer f.Close()
	// Get the file
	resp, err := http.Get(SNURL + file)
	if err != nil {
		return 0, fmt.Errorf("error processing %s: %s", file, err)
	}
	defer resp.Body.Close()
	var b int64
	for {
		n, err := io.Copy(f, resp.Body)
		b += n
		if err != nil {
			if err == io.EOF {
				return b, nil
			}
			return b, fmt.Errorf("error processing %s: %s", file, err)
		}
		// no bytes copied == done
		if n == 0 {
			return b, nil
		}
	}
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
