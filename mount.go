package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"log"
	"net/http"
	"time"
)

var (
	broker *Broker
)

func mount(war string) {
	log.Println("mount")

	m := martini.Classic()
	m.Handlers(martini.Recovery())
	m.Use(martini.Static(war, martini.StaticOptions{SkipLogging: true}))
	m.Use(render.Renderer(render.Options{
		Extensions: []string{".html", ".shtml"},
	}))
	broker = NewBroker()
	m.Map(broker)

	m.Get("/", indexHandler)
	m.Get("/events/", sseHandler)
	http.Handle("/", m)

	go func() {
		for i := 0; ; i++ {

			// Create a little message to send to clients,
			// including the current time.
			broker.messages <- fmt.Sprintf(" %d - 哈哈，Oh, Yeah!!! - the time is %v", i, time.Now())

			// Print a nice log message and sleep for 5s.
			log.Printf("Sent message %d ", i)
			time.Sleep(5 * 1e9)

		}
	}()
}
