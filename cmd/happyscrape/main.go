package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	// graceful shutdown
	"context"
	// "net"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/awsom82/happyscrape"
)

var (
	gitHash string = "NOBUILD"
	gitTime string = "0"
)

var happyscrapeUsage = func() {
	var useText string = `You should be able send JSON payloads to localhost:8080.
Notice, there no specific path for JSON, you may use any.

The application will detect an input type by a mime-type header.
If it lacks that data, it will try to detect that by file signature.

Examples:
> http :8080 Content-type:application/json < example.json`

	appVersion := fmt.Sprintf("Version:\n  Build %s at %s\n\nAuthor:\n  Igor A. Melekhine — 2020 © MIT License\n\n", strings.ToUpper(gitHash[:7]), gitTime)

	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n\n"+useText+"\n\n"+appVersion, os.Args[0])
	flag.PrintDefaults()
}

func main() {

	flag.Usage = happyscrapeUsage

	conf := happyscrape.NewConfig()

	flag.StringVar(&conf.Hostname, "hostname", conf.Hostname, "Bind server address")
	flag.IntVar(&conf.Port, "port", conf.Port, "Port number")
	flag.IntVar(&conf.LinksLimit, "links", conf.LinksLimit, "Max links number in payload")
	flag.IntVar(&conf.SimultaneousReqs, "max-requests", conf.SimultaneousReqs, "Max simultaneous requests")
	flag.IntVar(&conf.MaxConcurrentReqs, "outbound-conn", conf.MaxConcurrentReqs, "Max outbound requests")
	flag.DurationVar(&conf.ClientRequestTimeout, "outbound-timeout", conf.ClientRequestTimeout, "Outgoing requests timeout")
	flag.BoolVar(&conf.KeepAlive, "keep-alive", conf.KeepAlive, "HTTP Keep-Alive")
	flag.DurationVar(&conf.ReadTimeout, "read-timeout", conf.ReadTimeout, "HTTP Read timeout")
	flag.DurationVar(&conf.WriteTimeout, "write-timeout", conf.WriteTimeout, "HTTP Write timeout")
	flag.DurationVar(&conf.ConnTimeout, "conn-timeout", conf.ConnTimeout, "HTTP Connection timeout")
	flag.DurationVar(&conf.ShutdownTimeout, "shutdown-timeout", conf.ShutdownTimeout, "Timeout to complete requests before shutdown")

	flag.Parse()

	// run main context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("Happyscrape service starting...")

	srv := happyscrape.NewServer(conf)
	conf.Init() // TODO: not a best solution, but...
	defer conf.Close()

	// run server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Happyscrape serve error:", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Happyscrape shutting down...")

	// The context is used to inform the server it has conf.ShutdownTimeout seconds to finish
	ctxTimeout, cancelTimeout := context.WithTimeout(ctx, conf.ShutdownTimeout)
	defer cancelTimeout()
	if err := srv.Shutdown(ctxTimeout); err != nil {
		log.Fatal("Happyscrape forced to shutdown:", err)
	}

	log.Println("Happyscrape exiting")
}
