package happyscrape

import (
	//"encoding/json"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

var requestTimeout time.Duration = 15 * time.Second

// tokens limits active outbound connections
var tokens chan struct{}

// global client to reuse existing connections
var client http.Client

func init() {
	client = http.Client{}
}

type Link struct {
	URL  string
	Data string
}

func ScrapeLinks(ctx context.Context, links []Link) ([]Link, error) {
	var wg sync.WaitGroup
	var results []Link

	errch := make(chan error)
	parsed := make(chan Link, len(links))

	ctxCancel, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		for _, link := range links {
			tokens <- struct{}{}
			wg.Add(1)
			go processLink(ctxCancel, parsed, tokens, &wg, link, errch)
		}
		wg.Wait()
		close(parsed)
	}()

	for {
		select {
		case err := <-errch:
			return nil, fmt.Errorf("Abort by fetch error: %s", ErrWhileScraping, err)
		case <-ctx.Done():
			return nil, fmt.Errorf("Aborted by client...")
		case link, ok := <-parsed:
			results = append(results, link)
			if !ok {
				return results, nil
			}
		}
	}

	//res := returnResults(parsed)

	//return res, nil
}

func processLink(ctx context.Context, parsed chan<- Link, tokens <-chan struct{}, wg *sync.WaitGroup, link Link, errch chan<- error) {
	var err error
	var data bytes.Buffer
	defer wg.Done()

	log.Printf("[HappyScrape] GET %s\n", link.URL)

	ctxTimeout, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctxTimeout, "GET", link.URL, nil)
	if err != nil {
		errch <- err
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		errch <- err
		return
	}
	defer resp.Body.Close()

	<-tokens

	if resp.StatusCode != 200 {
		err = fmt.Errorf("[HappyScrape] Bad resp code %d for %s", resp.StatusCode, link.URL)
		errch <- err
		return
	}

	_, err = io.Copy(&data, resp.Body)
	if err != nil {
		errch <- err
		return
	}

	link.Data = base64.StdEncoding.EncodeToString(data.Bytes())

	parsed <- link
}

// func returnResults(parsed <-chan Link) []Link {
//
// }
