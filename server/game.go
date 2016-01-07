package server

import "github.com/lologarithm/survival/server/messages"

type GameMessage struct {
	net    messages.Net
	client *Client
	mtype  messages.MessageType
}
