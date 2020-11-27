package happyscrape

import (
	"fmt"
	"net/http"
)

// NewServer creates new server and limiter
func NewServer(conf *Config) *http.Server {

	h := http.HandlerFunc(ParserHandler)

	srv := http.Server{
		Addr:         fmt.Sprintf("%s:%d", conf.Hostname, conf.Port),
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
		Handler:      ScrapeLogMiddleware(h), // add middleware to see requests
	}

	srv.SetKeepAlivesEnabled(conf.KeepAlive)

	return &srv
}
