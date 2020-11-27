package happyscrape

import (
	// "context"
	"fmt"
	// "net"
	"net/http"
	// rate limiter
	// "golang.org/x/net/netutil"
)

// NewServer creates new server and limiter
func NewServer(conf *Config) *http.Server {

	h := http.HandlerFunc(ParserHandler)

	srv := http.Server{
		Addr:         fmt.Sprintf("%s:%d", conf.Hostname, conf.Port),
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
		Handler:      http.TimeoutHandler(ScrapeLogMiddleware(h, conf.SimultaneousReqs), conf.WriteTimeout, "Timeout!\n"), // add middleware to see requests and limit by set simply timeout
	}

	srv.SetKeepAlivesEnabled(conf.KeepAlive)

	return &srv
}
