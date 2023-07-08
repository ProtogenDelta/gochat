package main

import (
	"log"
	"net"

	gochat "github.com/protogendelta/gochat/lib"
)

type connection struct {
	sock *net.Conn
	nick string
}

func (c connection) setNick(nick string) connection {
	oldnick := c.nick
	c.nick = nick
	switch oldnick {
	case "":
		publish(makeJoin(nick))
	default:
		publish(makeNickChange(oldnick, nick))
	}
	return c
}

func makeMessage(cmsg *gochat.C2SMessage_ChatMessage_, user string) *gochat.S2CMessage {
	if cmsg.ChatMessage.IsAction {
		log.Printf("[#%s] %s %s", cmsg.ChatMessage.Channel, user, cmsg.ChatMessage.Content)
	} else {
		log.Printf("[#%s] %s: %s", cmsg.ChatMessage.Channel, user, cmsg.ChatMessage.Content)
	}
	return &gochat.S2CMessage{
		Content: &gochat.S2CMessage_ChatMessage_{
			ChatMessage: &gochat.S2CMessage_ChatMessage{
				Name:     user,
				Channel:  cmsg.ChatMessage.Channel,
				Content:  cmsg.ChatMessage.Content,
				IsAction: cmsg.ChatMessage.IsAction,
			},
		},
	}
}

func makeJoin(name string) *gochat.S2CMessage {
	log.Printf("+ %s", name)
	return &gochat.S2CMessage{
		Content: &gochat.S2CMessage_Join_{
			Join: &gochat.S2CMessage_Join{
				Name: name,
			},
		},
	}
}

func makePart(name string) *gochat.S2CMessage {
	log.Printf("- %s", name)
	return &gochat.S2CMessage{
		Content: &gochat.S2CMessage_Part_{
			Part: &gochat.S2CMessage_Part{
				Name: name,
			},
		},
	}
}

func makeNickChange(from string, to string) *gochat.S2CMessage {
	log.Printf("* %s -> %s", from, to)
	return &gochat.S2CMessage{
		Content: &gochat.S2CMessage_NickChange_{
			NickChange: &gochat.S2CMessage_NickChange{
				From: from,
				To:   to,
			},
		},
	}
}
