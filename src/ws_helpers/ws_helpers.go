package ws_helpers

import (
	"github.com/gorilla/websocket"
	"net"
	"sync"
	"encoding/json"
	"log"
	"code.google.com/p/go-uuid/uuid"
	"fmt"
)

var ActiveClients = make(map[ClientConn]int)
var ActiveClientsRWMutex sync.RWMutex
var Messages = make(map[string]*Message)
var MessagesRWMutex sync.RWMutex

type Message struct {
	Type    	 string
	Message 	 interface {}
	Id			 string
	Reply   	 string `json:"-"`
	Channel chan string	`json:"-"`
}

func AddMessage(m *Message) {
	MessagesRWMutex.Lock()
	m.Channel = make(chan string)
	Messages[m.Id] = m
	MessagesRWMutex.Unlock()
}

func DeleteMessage(m *Message) {
	MessagesRWMutex.Lock()
	delete(Messages, m.Id)
	MessagesRWMutex.Unlock()
}

func NewMessage(msgtype string, msg interface {}) *Message {
	m := new(Message)
	m.Type = msgtype
	m.Message = msg
	return m
}

type ClientConn struct {
	Websocket *websocket.Conn
	ClientIP  net.Addr
	Id        string
	MessageType int
	LastMessage interface {}
}

func AddClient(cc ClientConn) {
	ActiveClientsRWMutex.Lock()
	ActiveClients[cc] = 0
	ActiveClientsRWMutex.Unlock()
}

func DeleteClient(cc ClientConn) {
	ActiveClientsRWMutex.Lock()
	delete(ActiveClients, cc)
	ActiveClientsRWMutex.Unlock()
}

func BroadcastMessage(msg interface {}) {
	ActiveClientsRWMutex.RLock()
	defer ActiveClientsRWMutex.RUnlock()
	ret, _ := json.Marshal(msg)

	for client, _ := range ActiveClients {
		if err := client.Websocket.WriteMessage(1, ret); err != nil {
			log.Fatalln(client.Id, err)
			return
		}
	}
}

func (client *ClientConn) SendMessage(msg *Message) (*Message, error){
	msg.Id = uuid.New()
	ret, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
	}
	if err := client.Websocket.WriteMessage(1, ret); err != nil {
		fmt.Println(err)
		return nil, err
	}
	AddMessage(msg)
	return msg, nil
}

func (client *ClientConn) SendError(error string) {
	client.SendMessage(NewMessage("error", error))
}
