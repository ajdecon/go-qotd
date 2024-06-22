package server

import (
	"fmt"
	"log"
	"log/syslog"
	"net/http"
)

// External interface to the HttpQotdServer
type HttpQotdServer interface {
	Start(port int, quotes chan string) error
	Close()
	SetMaxLength(ml int)
	SetDebug(debug bool)
}

type httpServer struct {
	Port      int
	Quotes    chan string
	Done      chan bool
	Debug     bool
	MaxLength int
}

func NewHTTP() HttpQotdServer {
	return &httpServer{Port: 0, Quotes: nil, MaxLength: 512, Done: make(chan bool), Debug: false}
}

func (https *httpServer) Start(port int, quotes chan string) error {
	https.Port = port
	https.Quotes = quotes
	go httpRespondToQuotes(https)
	return nil
}

func (https *httpServer) Close() {
	https.Done <- true
}

func (https *httpServer) SetDebug(debug bool) {
	https.Debug = debug
}

func (https *httpServer) SetMaxLength(ml int) {
	https.MaxLength = ml
}

func (https *httpServer) LogAndDie(msg string) {
	https.Close()
	log.Fatalln(msg)
}

func quotesRoute(https *httpServer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := <-https.Quotes
		if len(q) > https.MaxLength {
			q = string([]byte(q)[0 : https.MaxLength-3]) // MaxLength - 3 for ellipses
			q = q + "..."
		}
		fmt.Fprintf(w, q)
	})
}

func httpRespondToQuotes(https *httpServer) {
	port := fmt.Sprintf(":%d", https.Port)

	_, err := syslog.NewLogger(syslog.LOG_INFO, 1)
	if err != nil {
		https.LogAndDie("Could not create a new syslog logger")
	}

	if https.Debug {
		fmt.Printf("[Will listen on %s]\n", port)
	}

	// Set up default route and listen for connections
	http.Handle("/", quotesRoute(https))
	http.ListenAndServe(port, nil)
}
