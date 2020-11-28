package happyscrape

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

const maxBodyBytes = int64(65536) // limit payload to 64kb

var maxLinks int = 20 // maximum links per request

var (
	ErrPayload         = errors.New("Bad payload")
	ErrWhileScraping   = errors.New("Error while scraping")
	ErrWhileMarshaling = errors.New("Error while marshaling")
)

// ParserHandler a http.handler function
func ParserHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		var err error
		var payload bytes.Buffer
		var links []Link

		// Copy payload
		r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

		_, err = io.Copy(&payload, r.Body)
		if err != nil {
			err = fmt.Errorf("[HappyScrape] %w: %s", ErrPayload, "Input load error: Data not provided or corrupted")
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check content-type
		ct := r.Header.Get("Content-type")
		if ct == "" || ct == "application/octet-stream" {
			ct = http.DetectContentType(payload.Bytes())
		}

		if ct != "application/json" {
			err = fmt.Errorf("[HappyScrape] %w: %s", ErrPayload, "Content-type mismatch: Looks like its not a json file")
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Load data
		err = json.Unmarshal(payload.Bytes(), &links)
		if err != nil {
			err = fmt.Errorf("[HappyScrape] %w: %s -> %s", ErrPayload, "Cant unmarshal json, please check you payload format", err)
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check links not exceed maximum length
		if len(links) > maxLinks {
			err = fmt.Errorf("[HappyScrape] %w: Too many links; do not send more than %d per request", ErrPayload, maxLinks)
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Scrape data
		res, err := ScrapeLinks(r.Context(), links)
		if err != nil {
			err = fmt.Errorf("[HappyScrape] %w: %s", ErrWhileScraping, err)
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		js, err := json.Marshal(res)
		if err != nil {
			err = fmt.Errorf("[HappyScrape] %w: %s", ErrWhileMarshaling, err)
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

	}

}
