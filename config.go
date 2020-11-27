package happyscrape

import (
	"time"
)

type Config struct {
	Hostname             string
	Port                 int           // pflag doesnt support uint16 out of the box
	LinksLimit           int           // according to business request
	SimultaneousReqs     int           // according to business request
	MaxConcurrentReqs    int           // max concurrent requests NOTICE: made as client, not server!
	ClientRequestTimeout time.Duration // CleintRequestTimeout sets timeout for url requests as client
	KeepAlive            bool
	ReadTimeout          time.Duration
	WriteTimeout         time.Duration
	ShutdownTimeout      time.Duration // number of time to complete existing requests before shutdown
}

func NewConfig() *Config {
	return &Config{"localhost", 8080, 20, 100, 5, 5e8, false, 5e9, 10e9, 5e9}
}

// Init must be called before server run
func (c *Config) Init() {
	tokens = make(chan struct{}, c.MaxConcurrentReqs) // max concurrent requests
	maxLinks = c.LinksLimit
	requestTimeout = c.ClientRequestTimeout
}

// Close uncessary but let it be
func (c *Config) Close() {
	close(tokens)
}
