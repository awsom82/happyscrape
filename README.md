#  JSON Web Scraper
Scrape the web, but only for demo purposes! Don't be bad.

## Installing
Just execute in terminal `git clone github.com/awsom82/happyscrape && cd happyscrape`
and run command `go build ./cmd/happyscrape && ./happyscrape`.

This will run conversion http service at port 8080.

## Using
After your run this app, you should able to send any JSON file to `http://localhost:8080/`.

Notice, there no specific path for JSON. The application will detect an input type of file by a mime-type header, or if it lacks that info. It will try to detect that by file signature [MIME Sniffing](https://mimesniff.spec.whatwg.org)

### Configuration
You can just type `./happyscrape --help` to get help message.
```
Usage of ./happyscrape:

You should able to send JSON file localhost:8080.
Notice, there no specific path for JSON, you may use any.

The application will detect an input type of file by a mime-type header.
If it lacks that info, it will try to detect that by file signature.

Examples:
> http :8080 Content-type:application/json < links.json

Version:
  Build NOBUILD at 0

Author:
  Igor A. Melekhine – 2020 © MIT License

  -hostname string
    	Bind server address (default "localhost")
  -keep-alive
    	HTTP Keep-Alive
  -links int
    	Max links number in payload (default 20)
  -max-requests int
    	Max simultaneous requests (default 100)
  -outbound-conn int
    	Max outbound requests (default 5)
  -outbound-timeout duration
    	Timeout for outgoing requests (default 500ms)
  -port int
    	Port number (default 8080)
  -read-timeout duration
    	HTTP Read timeout (default 5s)
  -shutdown-timeout duration
    	Seconds to complete requests before shutdown (default 5s)
  -write-timeout duration
    	HTTP Write timeout (default 10s)
```

### Example
```
$ http :8080 Content-type:application/json < example.json
```

### wrk
Good start point for testing is run app with next keys:

`./happyscrape -outbound-timeout 8s`.

Use `wrk -t1 -c8 -d30s -s wrkpayload.lua http://localhost:8080/` for simply load test