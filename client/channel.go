package main

import (
	"fmt"
)

type message interface {
	Draw() string
}

type channelMap map[string][]message

func (cm channelMap) AddMsg(name string, msg message) {
	cm[name] = append(cm.Get(name), msg)
}

func (cm channelMap) Clear(name string) {
	cm[name] = []message{systemMsg{content: "Channel history was cleared."}}
}

func (cm channelMap) Get(name string) []message {
	msgs, ok := cm[name]
	if !ok {
		msgs = make([]message, 0)
	}
	return msgs
}

func (cm channelMap) Draw(name string) string {
	msgs := cm.Get(name)
	if len(msgs) == 0 {
		return color(" There are no Messages in this channel.", SystemColor)
	}
	content := ""
	for _, msg := range msgs {
		content += fmt.Sprintf(" %s\n", msg.Draw())
	}
	return content
}
