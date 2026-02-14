package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetPrefix("ai-radar: ")
	// Set up channel on which to send signal notifications. We must use a
	// buffered channel or risk missing the signal if we're not ready to
	// receive when the signal is sent.
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	parser := Parser{List: make([]Element, 5)}
	socket := flag.String("name", "/tmp/ai-radar.sock", "socket's name")
	flag.Parse()

	l, err := net.Listen("unix", *socket)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := l.Close(); err != nil {
			log.Print(err)
		}
	}()

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				// This error "ErrClosed" is returned when the
				// listener is closed during shutdownw. In this
				// case we just call "return" to exit this
				// goroutine.
				if errors.Is(err, net.ErrClosed) {
					return
				}
				log.Print(err)
				continue
			}

			func() {
				defer func() {
					if err := conn.Close(); err != nil {
						log.Print(err)
					}
				}()

				buf := make([]byte, 512)
				n, err := conn.Read(buf)
				if err != nil {
					log.Print(err)
					return
				}

				text, err := parser.Parse(buf[:n])
				if err != nil {
					log.Print(err)
					return
				}

				// Waybar expects the exec-script to output its
				// data in JSON format. This should look like
				// this:
				//   {
				//     "alt": "$alt",
				//     "text": "$text",
				//     "class": "$class",
				//     "tooltip": "$tooltip",
				//     "percentage": $percentage
				//   }
				fmt.Printf("{\"text\":%q}\n", text)
			}()
		}
	}()

	log.Printf("listening to socket: %s", *socket)
	fmt.Println("{\"text\":\"-\"}")
	// Block until any signal is received.
	log.Print(<-c)
}
