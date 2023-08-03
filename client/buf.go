package main

import (
	gochat "github.com/protogendelta/gochat/lib/gochat/v1"
)

func makeIdent(name string) *gochat.C2SMessage {
	return &gochat.C2SMessage{
		Content: &gochat.C2SMessage_Ident_{
			Ident: &gochat.C2SMessage_Ident{
				Name: name,
			},
		},
	}
}

func makeMessage(channel, content string, action bool) *gochat.C2SMessage {
	return &gochat.C2SMessage{
		Content: &gochat.C2SMessage_ChatMessage_{
			ChatMessage: &gochat.C2SMessage_ChatMessage{
				Channel:  channel,
				Content:  content,
				IsAction: action,
			},
		},
	}
}
