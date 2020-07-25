package client

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"simulate-client/packet"
	"strings"
)

type Channel struct {
	conn    *websocket.Conn
	send    chan string
	receive chan []byte
	close   chan struct{}
}

func GetClientChannel(serverUrl string) *Channel {
	var schema, host, path string
	index := strings.Index(serverUrl, ":")

	if index > -1 {
		schema = (serverUrl)[0:index]
		host = (serverUrl)[index+3:]
		pathIndex := strings.Index(host, "/")
		if pathIndex > -1 {
			path = host[pathIndex+1:]
			host = host[:pathIndex]
		}
	}
	u := url.URL{Scheme: schema, Host: host, Path: path}
	log.Println(u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial : ", err)
	}
	closeChannel := make(chan struct{})
	channel := &Channel{
		conn:    conn,
		receive: make(chan []byte),
		close:   closeChannel,
	}

	return channel
}

func (channel *Channel) Read() {
	channel.RegisterHandler(func() {
		msgType, message, err := channel.conn.ReadMessage()
		log.Println("msg type", msgType)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("read error: ", err)
			}
			channel.Close()
			return
		}
		channel.receive <- message
	}, nil)
}

func (channel *Channel) Close() {
	if channel.conn == nil {
		return
	}
	log.Println("close channel ")
	close(channel.close)
	close(channel.receive)
	channel.conn.Close()
	channel.conn = nil
}

func (channel *Channel) GetMsg() string {
	s := <-channel.receive
	fmt.Println("receive text ", s)
	return string(packet.Decode(s))
}

func (channel *Channel) GetCloseChannel() chan struct{} {
	return channel.close
}

func (channel *Channel) SendMsg(msg string) {
	if msg == "" {
		return
	}
	req := packet.Encode(msg)
	err := channel.conn.WriteMessage(2, req)
	if err != nil {
		log.Println("Send error:", err)
	}
}

func (channel *Channel) IsClose() bool {
	return channel.conn == nil
}

func (channel *Channel) RegisterHandler(defaultHandler func(), closingHandler func()) {
	go func() {
		for {
			select {
			case <-channel.GetCloseChannel():
				if closingHandler != nil {
					closingHandler()
				}
				return
			default:

				if defaultHandler != nil {
					defaultHandler()
				}
			}
		}
	}()
}
