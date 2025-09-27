package chapter_3

import (
	"io"
	"net"
	"testing"
)

func TestListener(t *testing.T) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = listener.Close() }()

	t.Logf("bound to %q", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			t.Fatal(err)
		}
		go func(conn net.Conn) {
			defer conn.Close()
			// handle logic
		}(conn)
	}
}

func TestDial(t *testing.T) {
	// create a listener on a random port
	listener, err := net.Listen("tcp", ":")
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{})

	go func() {
		defer func() { done <- struct{}{} }()
		for {
			conn, err := listener.Accept()
			if err != nil {
				t.Log(err)
				return
			}
			go func(conn net.Conn) {
				defer func() {
					conn.Close()
					done <- struct{}{}
				}()
				buf := make([]byte, 1024)
				for {
					n, err := conn.Read(buf)
					if err != nil {
						if err != io.EOF {
							t.Error(err)
						}
						return
					}
					t.Logf("received: %q", buf[:n])
				}
			}(conn)
		}
	}()
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	conn.Close()
	<-done
	listener.Close()
	<-done
}
