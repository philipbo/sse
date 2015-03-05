package main

import (
	"fmt"
	"github.com/martini-contrib/render"
	"log"
	"net/http"
)

type Broker struct {
	clients        map[chan string]bool
	newClients     chan chan string
	defunctClients chan chan string
	messages       chan string
}

func NewBroker() *Broker {
	b := &Broker{
		make(map[chan string]bool),
		make(chan (chan string)),
		make(chan (chan string)),
		make(chan string),
	}

	b.Start()

	return b
}

func (b *Broker) Start() {
	go func() {
		for {
			select {
			case s := <-b.newClients:
				b.clients[s] = true
				log.Println("Added new client")
			case s := <-b.defunctClients:
				delete(b.clients, s)
				log.Println("Removed client")
			case msg := <-b.messages:
				for s, _ := range b.clients {
					s <- msg
				}
				log.Printf("Broadcast message to %d clients", len(b.clients))
			}
		}
	}()
}

func indexHandler(rd render.Render, r *http.Request) {
	log.Println("Finished HTTP request at ", r.URL.Path)
	rd.HTML(200, "index", "HTML5 服务器发送事件")
}

func sseHandler(w http.ResponseWriter, r *http.Request, b *Broker) {
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	messageChan := make(chan string)

	b.newClients <- messageChan

	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		b.defunctClients <- messageChan
		log.Println("HTTP connection just closed.")
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		msg := <-messageChan
		fmt.Fprintf(w, "data: Message %s\n\n", msg)
		f.Flush()
	}

	log.Println("Finished HTTP request at ", r.URL.Path)
}
