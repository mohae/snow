package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"golang.org/x/net/html"
)

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

	// if lastN processing is being done, set the start point
	if cnf.lastN > 0 {
		cnf.startEpisode = i - cnf.lastN + 1
		// make sure it's within range
		if cnf.startEpisode < 0 {
			cnf.startEpisode = 1
		}
	} else {
		// otherwise everything will be downloaded, start at episode 1
		cnf.startEpisode = 1
	}
	return nil
}
