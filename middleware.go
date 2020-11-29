package happyscrape

import (
	"log"
	"net/http"
)

// Log colors
const (
	endString    = " %0.2f kB \u2192 %s"
	InfoColor    = "[HappyScrape] \033[1;34m%s\033[0m" + endString
	WarningColor = "[HappyScrape] \033[1;33m%s\033[0m" + endString
)

var requests chan struct{}

// ScrapeLogMiddleware req limiter and logs http requests
func ScrapeLogMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case requests <- struct{}{}:
			defer func() { <-requests }()

			color := InfoColor
			if r.Method != "POST" {
				color = WarningColor
			}
			log.Printf(color, r.Method, float64(r.ContentLength)/1024.0, r.URL.Path)

			next.ServeHTTP(w, r)

		default:
			log.Println("LIMIT:", r.URL.Path)
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		}
	})

}
