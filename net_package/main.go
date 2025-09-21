package main

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)
}

type Message struct {
	from    string
	payload []byte
}

type Server struct {
	listenAddr string
	listener   net.Listener
	msgch      chan Message
	quitch     chan struct{}
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		msgch:      make(chan Message, 10),
	}
}

func (svr *Server) Start() error {
	listener, err := net.Listen("tcp", svr.listenAddr)
	if err != nil {
		return err
	}
	defer listener.Close()
	svr.listener = listener

	go svr.acceptLoop()
	<-svr.quitch

	close(svr.msgch)
	return nil
}

func (svr *Server) acceptLoop() {
	for {
		conn, err := svr.listener.Accept()
		log.Printf("new incoming connection. addr=%s", conn.RemoteAddr())
		if err != nil {
			log.Printf("accept error. err=%v", err)
			continue
		}
		go svr.readLoop(conn)
	}
}

func (svr *Server) readLoop(conn net.Conn) {
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("conn closed, bailing out. err=%v", err)
				return
			}
			log.Printf("read error. err=%v", err)
			continue
		}
		svr.msgch <- Message{
			from:    conn.RemoteAddr().String(),
			payload: buf[:n],
		}
		conn.Write([]byte("your message has been received\n"))
	}
}

func main() {
	server := NewServer(":3000")
	go func() {
		for data := range server.msgch {
			log.Printf("message from connection (%s): %s", data.from, string(data.payload))
		}
	}()
	log.Printf("starting a new server on port=%s", server.listenAddr)
	log.Fatal(server.Start())
}
