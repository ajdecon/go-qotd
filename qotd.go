package main

import (
	"flag"
	"fmt"
	"github.com/ajdecon/go-qotd/quotes"
	"github.com/ajdecon/go-qotd/server"
	"time"
)

var port = flag.Int("port", 17, "Port to run QOTD on")
var quotesFile = flag.String("file", "/tmp/quotes", "File to get quotes from")
var debug = flag.Bool("debug", false, "Print debug messages")

func main() {
	// Get options
	flag.Parse()

	// Get a quote generator
	qchannel := quotes.FileGenerator(*quotesFile)

	udps := server.NewUDP()
	udps.SetDebug(*debug)
	tcps := server.NewTCP()
	tcps.SetDebug(*debug)

	udps.Start(*port, qchannel)
	tcps.Start(*port, qchannel)

	if *debug {
		fmt.Printf("QOTD listening on port %d\n", *port)
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
