package main

import (
	"flag"
	"fmt"
	"github.com/ajdecon/go-qotd/quotes"
	"github.com/ajdecon/go-qotd/server"
	"time"
)

var port = flag.Int("port", 17, "Port to run QOTD on")
var udpEnable = flag.Bool("udpserver", false, "Run UDP server")
var tcpEnable = flag.Bool("tcpserver", true, "Run TCP server")

var httpEnable = flag.Bool("httpserver", false, "Run HTTP server")
var httpPort = flag.Int("httpport", 8080, "Port to run HTTP server on")

var quotesFile = flag.String("file", "./sample.data", "File to get quotes from")
var debug = flag.Bool("debug", false, "Print debug messages")
var maxLength = flag.Int("maxlen", 512, "Maximum length of quote to return, longer are trimmed")

func main() {
	// Get options
	flag.Parse()

	// Get a quote generator
	qchannel := quotes.FileGenerator(*quotesFile)

	if *udpEnable {
		udps := server.NewUDP()
		udps.SetDebug(*debug)
		udps.SetMaxLength(*maxLength)
		udps.Start(*port, qchannel)

		if *debug {
			fmt.Printf("QOTD listening on UDP port %d\n", *port)
		}
	}

	if *tcpEnable {
		tcps := server.NewTCP()
		tcps.SetDebug(*debug)
		tcps.SetMaxLength(*maxLength)
		tcps.Start(*port, qchannel)

		if *debug {
			fmt.Printf("QOTD listening on TCP port %d\n", *port)
		}
	}

	if *httpEnable {
		https := server.NewHTTP()
		https.SetDebug(*debug)
		https.SetMaxLength(*maxLength)
		https.Start(*httpPort, qchannel)

		if *debug {
			fmt.Printf("QOTD launched HTTP server on port %d\n", *httpPort)
		}
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
