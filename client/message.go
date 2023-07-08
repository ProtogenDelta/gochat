package main

// ----- Generic User Message -----
type userMsg struct {
	User     string
	Channel  string
	Content  string
	IsAction bool
}

// ----- System Message -----
type systemMsg struct {
	content string
}

// ----- User Join & Leave -----
type userJoinMsg struct {
	User string
}

type userPartMsg struct {
	User string
}

// ----- User Name Change (/nick) -----
type userNickMsg struct {
	From string
	To   string
}
