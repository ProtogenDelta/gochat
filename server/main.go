package main

import (
	"fmt"
	"log"
	"net"
	"sync"

	gochat "github.com/protogendelta/gochat/lib/gochat/v1"

	"google.golang.org/protobuf/proto"
)

var (
	conns = make(map[net.Addr]connection)
	mx    sync.RWMutex
)

const ServerPort = 1054

func main() {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", ServerPort))
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Printf("Listening on 0.0.0.0:%d", ServerPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	sub(conn)
	defer unsub(conn)

	for {
		data := make([]byte, 4096)
		length, err := conn.Read(data)
		if err != nil {
			return
		}

		var payload gochat.C2SMessage
		err = proto.Unmarshal(data[:length], &payload)
		if err != nil {
			log.Println(err)
			continue
		}

		handleC2S(conn, &payload)
	}
}

func handleC2S(conn net.Conn, payload *gochat.C2SMessage) {
	switch msg := payload.Content.(type) {
	case *gochat.C2SMessage_Ident_:
		conns[conn.RemoteAddr()] = conns[conn.RemoteAddr()].setNick(msg.Ident.Name)
	case *gochat.C2SMessage_ChatMessage_:
		nick := conns[conn.RemoteAddr()].nick
		if nick != "" {
			publish(makeMessage(msg, nick))
		}
	}
}

func sub(conn net.Conn) {
	mx.Lock()
	defer mx.Unlock()
	conns[conn.RemoteAddr()] = connection{
		sock: &conn,
	}
}

func unsub(conn net.Conn) {
	if nick := conns[conn.RemoteAddr()].nick; nick != "" {
		publish(makePart(nick))
	}
	mx.Lock()
	defer mx.Unlock()
	delete(conns, conn.RemoteAddr())
}

func publish(msg *gochat.S2CMessage) {
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Printf("Failed to encode message: %s", err)
		return
	}
	publishRaw(data)
}

func publishRaw(data []byte) {
	mx.RLock()
	defer mx.RUnlock()

	for _, conn := range conns {
		(*conn.sock).Write(data)
	}
}
