package tcp

import (
	"io"
	"log"
	"net"
	"net/url"
	"os"

	"github.com/yaproxy/libyap/proxy"
)

// Progress indicates transfer status
type Progress struct {
	bytes uint64
}

var silentMode bool

// TransferStreams launches two read-write goroutines and waits for signal from them
func TransferStreams(con net.Conn) {
	c := make(chan Progress)

	// Read from Reader and write to Writer until EOF
	copy := func(r io.ReadCloser, w io.WriteCloser) {
		defer func() {
			r.Close()
			w.Close()
		}()
		n, err := io.Copy(w, r)
		if err != nil {
			if !silentMode {
				log.Printf("[%s]: ERROR: %s\n", con.RemoteAddr(), err)
			}
		}
		c <- Progress{bytes: uint64(n)}
	}

	go copy(con, os.Stdout)
	go copy(os.Stdin, con)
	p1 := <-c
	p2 := <-c
	if !silentMode {
		log.Printf("[%s]: Connection has been closed by remote peer, %d bytes has been received\n", con.RemoteAddr(), p1.bytes)
		log.Printf("[%s]: Local peer has been stopped, %d bytes has been sent\n", con.RemoteAddr(), p2.bytes)
	}
}

// StartServer starts TCP listener
func StartServer(proto string, port string) {
	ln, err := net.Listen(proto, port)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Listening on", proto+port)
	con, err := ln.Accept()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("[%s]: Connection has been opened\n", con.RemoteAddr())
	TransferStreams(con)
}

// StartClient starts TCP connector
func StartClient(proto string, host string, port string, proxyURL string, silent bool) error {
	var dial func(network, address string) (net.Conn, error)
	if proxyURL != "" {
		fixedURL, err := url.Parse(proxyURL)
		if err != nil {
			return err
		}
		dialer, err := proxy.FromURL(fixedURL, nil, nil)
		if err != nil {
			return err
		}
		dial = dialer.Dial
	} else {
		dial = net.Dial
	}
	silentMode = silent
	con, err := dial(proto, host+":"+port)
	if err != nil {
		return err
	}
	if !silentMode{
		log.Println("Connected to", host+port)
	}
	TransferStreams(con)
	return nil
}
