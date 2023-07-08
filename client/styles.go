package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const (
	UserColor   = "4"
	SystemColor = "8"
	JoinColor   = "10"
	PartColor   = "9"
)

// ----- Styles -----
func color(s, color string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(s)
}

// ----- Message formatters -----
func (m userMsg) Draw() string {
	if m.IsAction {
		return color(fmt.Sprintf("* %s %s", color(m.User, "4"), m.Content), "8")
	} else {
		return fmt.Sprintf("%s %s", color(fmt.Sprintf("%s:", m.User), "4"), m.Content)
	}
}

func (m systemMsg) Draw() string {
	return color(m.content, SystemColor)
}

func (m userJoinMsg) Draw() string {
	return color(fmt.Sprintf("+ %s", color(m.User, UserColor)), JoinColor)
}

func (m userPartMsg) Draw() string {
	return color(fmt.Sprintf("- %s", color(m.User, UserColor)), PartColor)
}

func (m userNickMsg) Draw() string {
	return color(fmt.Sprintf("* %s is now known as %s", color(m.From, UserColor), color(m.To, UserColor)), SystemColor)
}
