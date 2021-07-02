package main

type commandID int

const (
	CMD_NICK commandID = iota
	CMD_JOIN
	CMD_ROOMLIST
	CMD_MSG
	CMD_QUIT
	CMD_HELP
)

type command struct {
	id     commandID
	client *client
	args   []string
}
