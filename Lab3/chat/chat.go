// Demonstration of channels with a chat application
// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Chat is a server that lets clients chat with each other.

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type client struct{
	name string
	chl  chan<- string // an outgoing message channel
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func activeUsers(clients map[client]bool, ch chan<- string){
	if len(clients)<1{
		return
	}
	ch <- "Current active users are"
	for cl := range clients{
		ch <- cl.name
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			cliName := strings.Fields(msg)[0]
			for cli := range clients {
				if cli.name!= cliName{
					cli.chl <- msg
				}
			}

		case cli := <-entering:
			activeUsers(clients,cli.chl)
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli.chl)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	ch<- "What is your Name"
	input := bufio.NewScanner(conn)
	
	var name string
	if input.Scan(){
		name = input.Text()
	}
	ch <- "You are " + name
	messages <- name + " has arrived"
	usr := client{name:name,chl:ch}
	entering <- usr

	for input.Scan() {
		messages <- name + " : " + input.Text()
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- usr
	messages <- name + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}
