package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"fmt"
	"github.com/gorilla/pat"
	"html/template"
	"net/http"
)

var addr = flag.String("port", ":8080", "http service address")
var homeTempl = template.Must(template.ParseFiles("templates/index.html"))
var h = Hub{
	broadcast:   make(chan string),
	register:    make(chan *Connection),
	unregister:  make(chan *Connection),
	connections: make(map[*Connection]bool),
}

func checkError(err error, responseWriter http.ResponseWriter) bool {
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusNotFound)
		return true
	}
	return false
}

func Home(responseWriter http.ResponseWriter, request *http.Request) {
	homeTempl.Execute(responseWriter, request.Host)
}

func websocketHandler(ws *websocket.Conn) {
	c := &Connection{send: make(chan string, 256), ws: ws}
	h.register <- c
	defer func() { h.unregister <- c }()
	go c.writer()
	c.reader()
}

func main() {
	flag.Parse()
	go h.run()
	r := pat.New()
	r.Get("/", Home)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.Handle("/", r)
	http.Handle("/ws", websocket.Handler(websocketHandler))
	fmt.Printf("Server running on port %s!", *addr)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", *addr), nil); err != nil {
		fmt.Println("ListenAndServe:", err)
	}

}
