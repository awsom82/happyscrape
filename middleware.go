package happyscrape

import (
	"log"
	"net/http"
)

// Log colors
const (
	endString    = " %0.2f kB \u2192 %s"
	InfoColor    = "\033[1;34m%s\033[0m" + endString
	WarningColor = "\033[1;33m%s\033[0m" + endString
)

// ScrapeLogMiddleware req limiter and logs http requests
func ScrapeLogMiddleware(next http.Handler, concurrentReqs int) http.Handler {
	sem := make(chan struct{}, concurrentReqs)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sem <- struct{}{}
		defer func() { <-sem }()

		color := InfoColor

		if r.Method != "POST" {
			color = WarningColor
		}

		log.Printf(color, r.Method, float64(r.ContentLength)/1024.0, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}
