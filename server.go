package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"strings"

	_ "github.com/go-sql-driver/mysql"
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
	if len(args) < 2 {
		c.msg("Nick is required. usage: /nick NAME")
		return
	}

	c.nick = args[1]
	c.msg(fmt.Sprintf("OK, now you're %s", c.nick))
}

func (s *server) join(c *client, args []string) {
	if len(args) < 2 {
		c.msg("Room name is required. usage: /join ROOM_NAME")
		return
	}

	roomName := args[1]

	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}
	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s joined the room", c.nick))

	c.msg(fmt.Sprintf("welcome to %s", roomName))
}

func (s *server) listRooms(c *client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.msg(fmt.Sprintf("available rooms: %s", strings.Join(rooms, ", ")))
}

func (s *server) msg(c *client, args []string) {
	if len(args) < 2 {
		c.msg("message is required, usage: /msg MSG")
		return
	}

	msg := strings.Join(args[1:], " ")
	c.room.broadcast(c, c.nick+": "+msg)

	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/chat")
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	defer db.Close()

	insert, err := db.Query("INSERT INTO test (name, room, msg) VALUES (?, ?, ?)", c.nick, c.room.name, msg)
	if err != nil {
		log.Fatalf("Cannot upload to database: %s", err)
	}
	defer insert.Close()
}

func (s *server) quit(c *client) {
	log.Printf("client has left the chat: %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)

	c.msg("sad to see you go =(")
	c.conn.Close()
}

func (s *server) help(c *client) {
	c.msg("Commands:\n/nick\n/join\n/roomlist\n/msg\n/quit\n")
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		oldRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
		oldRoom.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}
