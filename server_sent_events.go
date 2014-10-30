package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Message struct {
	ID   int      `json:"id"`
	Text string   `json:"text"`
	Tags []string `json:"tags"`
}

// Regularly generate random messages
func messageSource(done <-chan struct{}) <-chan Message {
	messages := make(chan Message)
	texts := []string{"Hello, World!", "Literals are great.", "JSON is nifty."}
	tags := [][]string{
		{"good", "fun"},
		{"evil", "boring"},
		{"meh", "yawn", "whatever"},
	}

	go func() {
		defer close(messages)
		log.Println("messageSource started")
		defer log.Println("messageSource stopped")
		for i := 0; ; i++ {
			select {
			case messages <- Message{
				ID:   i,
				Text: texts[rand.Intn(len(texts))],
				Tags: tags[rand.Intn(len(tags))],
			}:
				<-time.After(1 * time.Second)

			case <-done:
				// Prevents goroutine from blocking forever if consumer dies
				return
			}
		}
	}()
	return messages
}

// Open a persistent connection and send generated messages
func MessageServer(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection from " + r.RemoteAddr)
	done := make(chan struct{})
	defer close(done)
	messages := messageSource(done)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Cache-Control", "no-cache")
	for {
		// Check to see if the connection has closed
		cn, ok := w.(http.CloseNotifier)
		if !ok {
			log.Printf("Error: ResponseWriter does not support CloseNotify()")
			return
		}
		select {
		case <-cn.CloseNotify():
			log.Printf("%s disconnected", r.RemoteAddr)
			return
		default:
		}

		// Marshal a Message into JSON
		message := <-messages
		bs, err := json.MarshalIndent(message, "", "  ")
		if err != nil {
			log.Fatal("MarhsalIndent: ", err)
		}

		// Put into the SSE format:
		// https://html.spec.whatwg.org/multipage/comms.html#the-eventsource-interface
		bs = append([]byte("data: "), bs...)
		bs = append(bs, []byte("\n\n")...)

		// Write and flush
		f, ok := w.(http.Flusher)
		if !ok {
			log.Println("Error: ResponseWriter does not support Flush()")
			return
		}
		w.Write(bs)
		f.Flush()
		log.Printf("Sent message #%d to %s", message.ID, r.RemoteAddr)
	}
}

func main() {
	http.HandleFunc("/messages", MessageServer)
	log.Println("Listening on http://127.0.0.1:8000/messages")
	err := http.ListenAndServe("127.0.0.1:8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
