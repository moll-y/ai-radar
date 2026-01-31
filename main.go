package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
)

func main() {
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	l, err := net.Listen("unix", "/tmp/test.sock")
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
				log.Print(err)
				return
			}

			func() {
				defer func() {
					if err := conn.Close(); err != nil {
						log.Print(err)
					}
				}()

				buf := make([]byte, 1024)
				n, err := conn.Read(buf)
				if err != nil {
					log.Print(err)
					return
				}

				cmd, err := parse(buf[:n])
				if err != nil {
					log.Print(err)
					return
				}

				fmt.Println(cmd)
			}()
		}
	}()

	// Block until any signal is received.
	// And then, clean up.
	<-c
}
