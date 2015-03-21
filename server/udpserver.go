package server

import (
    "net"
    "fmt"
    "log"
)

// External interface to the UdpQotdServer
type UdpQotdServer interface {
    // To start the server, provide a port on which to listen and a channel
    // "quotes" from which to obtain new quotes
    Start(port int, quotes chan string) error

    // Close the server
    Close()

    // Turn on debugging
    SetDebug(debug bool)
}

type udpServer struct {
    Port int
    Quotes chan string
    Done chan bool
    Debug bool
}

func NewUDP() UdpQotdServer {
    return &udpServer{Port:0, Quotes: nil, Done: make(chan bool)}
}

func (udps *udpServer) Start(port int, quotes chan string) error {
    udps.Port = port
    udps.Quotes = quotes
    go respondToQuotes(udps)
    return nil
}

func (udps *udpServer) Close() {
    udps.Done <- true
}

func (udps *udpServer) SetDebug(debug bool) {
    udps.Debug = debug
}

func (udps *udpServer) LogAndDie(msg string) {
    udps.Close()
    log.Fatalln(msg)
}

// ****************************************************************************
// Actual server logic follows
// ****************************************************************************

func respondToQuotes(udps *udpServer) {
    // Start listening for connections
    port := fmt.Sprintf(":%d", udps.Port)

    if udps.Debug {
        fmt.Printf("[Will listen on %s]\n", port)
    }

    udpAddr, err := net.ResolveUDPAddr("udp", port)
    if err != nil {
        udps.LogAndDie("Could not resolve local port")
        return
    }

    conn, err := net.ListenUDP("udp", udpAddr)
    if err != nil {
        udps.LogAndDie("Could not bind to address")
        return
    }

    // Start a goroutine to listen for Done signals
    go func() {
        var status bool
        for {
            status = <-udps.Done
            if status == true {
                conn.Close()
                return
            }
        }
    }()

    // Listen for new connections and send them a new quote
    var buf []byte
    for {
        _, addr, err := conn.ReadFromUDP(buf[0:])
        if err != nil {
            udps.LogAndDie("Could not read from UDP")
        }
        // QOTD don't care what you say!
        // Get a new quote from the channel!
        q := <- udps.Quotes
        buf = []byte(q)

        if udps.Debug {
            log.Printf("Sending: %s\n" , q)
        }

        conn.WriteToUDP(buf[0:], addr)
    }
}



