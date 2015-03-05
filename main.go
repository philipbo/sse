package main

import (
	"flag"
	"github.com/go-martini/martini"
	"log"
	"net/http"
)

func main() {
	log.Println("start", martini.Env)
	addr := flag.String("p", ":3000", "address where the server listen on")
	war := flag.String("war", "./war", "directory of war files")

	mount(*war)

	log.Printf("start web server on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
