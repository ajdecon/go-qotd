package server

import (
    "net"
    "fmt"
    "log"
    "log/syslog"
)

// External interface to the UdpQotdServer
type TcpQotdServer interface {
    // To start the server, provide a port on which to listen and a channel
    // "quotes" from which to obtain new quotes
    Start(port int, quotes chan string) error

    // Close the server
    Close()

    // Turn on debugging
    SetDebug(debug bool)
}

type tcpServer struct {
    Port int
    Quotes chan string
    Done chan bool
    Debug bool
}

func NewTCP() TcpQotdServer {
    return &tcpServer{Port:0, Quotes: nil, Done: make(chan bool)}
}

func (tcps *tcpServer) Start(port int, quotes chan string) error {
    tcps.Port = port
    tcps.Quotes = quotes
    go tcpRespondToQuotes(tcps)
    return nil
}

func (tcps *tcpServer) Close() {
    tcps.Done <- true
}

func (tcps *tcpServer) SetDebug(debug bool) {
    tcps.Debug = debug
}

func (tcps *tcpServer) LogAndDie(msg string) {
    tcps.Close()
    log.Fatalln(msg)
}

// ****************************************************************************
// Actual server logic follows
// ****************************************************************************

func tcpRespondToQuotes(tcps *tcpServer) {
    // Start listening for connections
    port := fmt.Sprintf(":%d", tcps.Port)

    // Start a syslog logger
    logMe, err := syslog.NewLogger(syslog.LOG_INFO, 1)
    if err != nil {
        tcps.LogAndDie("Could not create a new syslog logger")
    }

    if tcps.Debug {
        fmt.Printf("[Will listen on %s]\n", port)
    }

    // Begin listening for connections
    tcpAddr, err := net.ResolveTCPAddr("tcp4", port)
    if err != nil {
        tcps.LogAndDie("Could not resolve port")
    }
    listener, err := net.ListenTCP("tcp", tcpAddr)
    if err != nil {
        tcps.LogAndDie("Could not bind to port")
    }
    defer listener.Close()

    // Start a goroutine to listen for Done signals
    go func() {
        var status bool
        for {
            status = <-tcps.Done
            if status == true {
                listener.Close()
                return
            }
        }
    }()

    // Respond to connections
    for {
        conn, err := listener.Accept()
        logMe.Print(fmt.Sprintf("Received request from: %s", conn.RemoteAddr().String()))
        if err != nil {
            continue
        }
        go handleConnection(conn, tcps)
    }
}

func handleConnection(conn net.Conn, tcps *tcpServer) {
    defer conn.Close()
    q := <-tcps.Quotes
    _, _ = conn.Write([]byte(q))
}

