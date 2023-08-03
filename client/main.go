package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	gochat "github.com/protogendelta/gochat/lib/gochat/v1"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var sendOut = make(chan *gochat.C2SMessage)

func main() {
	defer close(sendOut)

	if len(os.Args) != 3 {
		fmt.Printf("Syntax:\n\t%s [username] [host]:[port]\n", os.Args[0])
		os.Exit(1)
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", os.Args[2])
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	m := initialModel(os.Args[1])
	p := tea.NewProgram(m, tea.WithAltScreen())

	go handleSend(conn)
	go handleRecv(conn, p)

	go func(m *model) {
		sendOut <- makeIdent(m.nick)
	}(&m)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
	channel  string
	nick     string
	channels channelMap
	viewport viewport.Model
	textarea textarea.Model
	err      error
}

func initialModel(nick string) model {
	ta := textarea.New()
	ta.Placeholder = "Message #general..."
	ta.Focus()

	ta.Prompt = " ┃ "
	ta.CharLimit = 280

	ta.SetWidth(80)
	ta.SetHeight(3)

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(80, 30)
	vp.SetContent(` Welcome to the chat room!
 Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		channel: "general",
		nick:    nick,
		channels: channelMap{
			"general": make([]message, 0),
		},
		viewport: vp,
		textarea: ta,
		err:      nil,
	}
}

func addMsgAll(m *model, msg message) {
	for channel := range m.channels {
		m.channels.AddMsg(channel, msg)
	}
	redraw(m)
}

func addMsg(m *model, channel string, msg message) {
	m.channels.AddMsg(channel, msg)
	redraw(m)
}

func redraw(m *model) {
	m.viewport.SetContent(m.channels.Draw(m.channel))
	m.viewport.GotoBottom()
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		taCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, taCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case userMsg:
		addMsg(&m, msg.Channel, msg)
	case userJoinMsg, userPartMsg, userNickMsg:
		addMsgAll(&m, msg.(message))
	case tea.WindowSizeMsg:
		m.textarea.SetWidth(msg.Width - 2)
		m.viewport.Width = msg.Width - 2
		m.viewport.Style.Width(msg.Width - 2)
		m.viewport.Height = msg.Height - 6
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			content := m.textarea.Value()
			m.textarea.Reset()
			if len(content) == 0 {
				break
			} else if strings.HasPrefix(content, "/") {
				return m.HandleCommand(content)
			} else {
				sendOut <- makeMessage(m.channel, content, false)
			}
		}
	case error:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(taCmd, vpCmd)
}

func (m model) HandleCommand(content string) (tea.Model, tea.Cmd) {
	words := strings.Fields(content)
	switch words[0] {
	case "/channel", "/c":
		if len(words) != 2 {
			addMsg(&m, m.channel, systemMsg{
				content: "Command Syntax: /c|channel {channel}",
			})
			break
		}
		if len(m.channels[m.channel]) == 0 {
			delete(m.channels, m.channel)
		}
		m.channel = words[1]
		m.textarea.Placeholder = fmt.Sprintf("Message #%s...", m.channel)
		if _, ok := m.channels[m.channel]; !ok {
			m.channels[m.channel] = make([]message, 0)
		}
		redraw(&m)
	case "/clear":
		m.channels.Clear(m.channel)
		redraw(&m)
	case "/help", "/?":
		addMsg(&m, m.channel, systemMsg{
			content: "Available Commands: channel, clear, help, me, nick, quit",
		})
	case "/me", "/act":
		if len(words) < 2 {
			addMsg(&m, m.channel, systemMsg{
				content: "Command Syntax: /me {action...}",
			})
			break
		}
		sendOut <- makeMessage(m.channel, strings.Join(words[1:], " "), true)
	case "/nick":
		if len(words) != 2 {
			addMsg(&m, m.channel, systemMsg{
				content: "Command Syntax: /nick {nick}",
			})
			break
		}
		m.nick = words[1]
		sendOut <- makeIdent(m.nick)
	case "/quit", "/q":
		return m, tea.Quit
	default:
		addMsg(&m, m.channel, systemMsg{
			content: fmt.Sprintf("Unknown Command: %s", words[0]),
		})
	}

	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("\n%s\n\n%s", m.viewport.View(), m.textarea.View())
}
