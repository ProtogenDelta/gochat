package main

import (
	"log"
	"net"

	tea "github.com/charmbracelet/bubbletea"
	gochat "github.com/protogendelta/gochat/lib/gochat/v1"
	"google.golang.org/protobuf/proto"
)

func handleSend(c *net.TCPConn) {
	for {
		data, err := proto.Marshal(<-sendOut)
		if err != nil {
			continue
		}
		_, err = c.Write(data)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func handleRecv(c *net.TCPConn, p *tea.Program) {
	c.SetKeepAlive(true)

	for {
		data := make([]byte, 4096)
		length, err := c.Read(data)
		if err != nil {
			return
		}

		var payload gochat.S2CMessage
		err = proto.Unmarshal(data[:length], &payload)

		switch msg := payload.Content.(type) {
		case *gochat.S2CMessage_ChatMessage_:
			p.Send(userMsg{
				User:     msg.ChatMessage.Name,
				Channel:  msg.ChatMessage.Channel,
				Content:  msg.ChatMessage.Content,
				IsAction: msg.ChatMessage.IsAction,
			})
		case *gochat.S2CMessage_Join_:
			p.Send(userJoinMsg{
				User: msg.Join.Name,
			})
		case *gochat.S2CMessage_Part_:
			p.Send(userPartMsg{
				User: msg.Part.Name,
			})
		case *gochat.S2CMessage_NickChange_:
			p.Send(userNickMsg{
				From: msg.NickChange.From,
				To:   msg.NickChange.To,
			})
		}
	}
}
