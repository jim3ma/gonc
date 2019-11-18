package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dddpaul/gonc/tcp"
	"github.com/dddpaul/gonc/udp"
)

func main() {
	var host, port, proto, proxy string
	var listen, silent bool
	flag.StringVar(&host, "host", "", "Remote host to connect, i.e. 127.0.0.1")
	flag.StringVar(&proto, "proto", "tcp", "TCP/UDP mode")
	flag.StringVar(&proxy, "proxy", "", "Proxy mode")
	flag.BoolVar(&listen, "listen", false, "Listen mode")
	flag.BoolVar(&silent, "silence", true, "Silent or quiet mode. Don't show progress meter or error messages")
	flag.StringVar(&port, "port", "22", "Port to listen on or connect to, i.e. 22")
	flag.Parse()

	switch proto {
	case "tcp":
		if listen {
			tcp.StartServer(proto, port)
		} else if host != "" {
			err := tcp.StartClient(proto, host, port, proxy, silent)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			flag.Usage()
		}
	case "udp":
		if listen {
			udp.StartServer(proto, port, silent)
		} else if host != "" {
			udp.StartClient(proto, host, port, silent)
		} else {
			flag.Usage()
		}
	default:
		flag.Usage()
	}
}
