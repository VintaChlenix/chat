package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_ROOMLIST:
			s.listRooms(cmd.client)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client)
		case CMD_HELP:
			s.help(cmd.client)
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("New client has joined: %s", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		nick:     "anonymous",
		commands: s.commands,
	}
	c.msg(fmt.Sprint("Welcome to the chat-server!"))
	c.readInput()
}

func (s *server) nick(c *client, args []string) {

}

func (s *server) join(c *client, args []string) {

}

func (s *server) listRooms(c *client) {

}

func (s *server) msg(c *client, args []string) {

}

func (s *server) quit(c *client) {

}

func (s *server) help(c *client) {
	c.msg("Commands:\n/nick\n/join\n/roomlist\n/msg\n/quit\n")
}
