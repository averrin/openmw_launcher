package main

import (
	"fmt"
	"os"
		"log"
	//	"time"
	//	"github.com/jmcvetta/randutil"
	"github.com/go-martini/martini"
	"github.com/gorilla/websocket"
	"github.com/martini-contrib/render"
	"net/http"
	wsh "ws_helpers"
	"encoding/json"
	"code.google.com/p/go-uuid/uuid"
	"time"
	"errors"
	"github.com/miketheprogrammer/go-thrust/dispatcher"
	"github.com/miketheprogrammer/go-thrust/session"
	"github.com/miketheprogrammer/go-thrust/spawn"
	"github.com/miketheprogrammer/go-thrust/window"
)

type Message struct {
	Id 		string `json:"Id"`
	Message string `json:"Message"`
	Type 	string `json:"Type"`
}

func WSHandler(w http.ResponseWriter, r *http.Request) {

	if len(wsh.ActiveClients) > 0 {
		return
	}
	
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}
	ip := ws.RemoteAddr()
	message := new(Message)
	client := wsh.ClientConn{ws, ip, uuid.New(), 0, message}

	wsh.AddClient(client)
	for {
		log.Println(len(wsh.ActiveClients), wsh.ActiveClients)
		messageType, msg, err := ws.ReadMessage()
		client.MessageType = messageType
		if err != nil {
			log.Println("Disconnected", client.Id)
			log.Println(err)
			return
		}
		err = json.Unmarshal(msg, &message)
		if err != nil {
			fmt.Println(err, string(msg))
			client.SendError(string(msg))
			continue
		}
		client.LastMessage = message
		
		switch message.Type {
		case "connect":
			client.SendMessage(wsh.NewMessage("connect", client.Id))
		case "reply":
			fmt.Println(message)
			wsh.Messages[message.Id].Reply = message.Message
			wsh.Messages[message.Id].Channel <- message.Message 
			defer wsh.DeleteMessage(wsh.Messages[message.Id])
			
		default:
			client.SendError("Unknown command")
		}
	}
}

func main() {
	options := NewOptions()
	if options.LauncherConfig["General"]["firstrun"] == "true" {
		fmt.Println("Its a first run of OpenMW, please run official omwlauncher for setting Morrowind path and initial settings")
		os.Exit(1)
	}

	for _, f := range options.GetAvailableContentFiles() {
		if Pos(f, options.ContentFiles.List) != -1 {
			fmt.Print(" [x] ")
		} else {
			fmt.Print(" [ ] ")
		}
		fmt.Println(f)
	}

	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
	Directory:  "templates",
	Extensions: []string{".tmpl", ".html"},
}))
	m.Use(martini.Static("static", martini.StaticOptions{Prefix: "static"}))

	m.Get("/sock", WSHandler)
	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", options)
	})

	go func() {
		for {
			time.Sleep(1000 * time.Millisecond)
			SetContent("#title", uuid.New())
		}
	}()

	spawn.SetBaseDirectory("./")
	spawn.Run()

	mysession := session.NewSession(false, false, "cache")

	thrustWindow := window.NewWindow("http://localhost:3000/", mysession)
	thrustWindow.Show()
	thrustWindow.Maximize()
	thrustWindow.Focus()

	// NonBLOCKING - note in other examples this was blocking.
	go dispatcher.RunLoop()

	m.Run()

}

func Get(msgtype string, data string) (string, error) {
	if len(wsh.ActiveClients) == 0 {
		return "", errors.New("No clients")
	}
	m := new(wsh.Message)
	err := errors.New("")
	for client, _ := range wsh.ActiveClients {
		m, err = client.SendMessage(wsh.NewMessage(msgtype, data))
		if err != nil {
			fmt.Println(err)
			return "", err
		}
	}
	var value string
	for s := range m.Channel {
		value = s
		break
	}
	return value, nil
}

func GetValue(selector string) (string, error){
	return Get("get_value", selector)
}
func SetValue(selector string, value string){
	data := struct{
			Selector string
			Content string
		}{selector, value}
	for client, _ := range wsh.ActiveClients {
		client.SendMessage(wsh.NewMessage("set_value", data))
	}
}
func SetContent(selector string, content string){
	data := struct{
		Selector string
		Content string
	}{selector, content}
	for client, _ := range wsh.ActiveClients {
		client.SendMessage(wsh.NewMessage("set_content", data))
	}
}
func GetContent(selector string) (string, error){
	return Get("get_content", selector)
}
