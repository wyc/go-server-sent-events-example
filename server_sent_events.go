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
			message := Message{
				ID:   i,
				Text: texts[rand.Intn(len(texts))],
				Tags: tags[rand.Intn(len(tags))],
			}

			select {
			case messages <- message:
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

	// Type assertions to check ResponseWriter for necessary features
	closeNotifier, ok := w.(http.CloseNotifier)
	if !ok {
		log.Printf("Error: ResponseWriter does not support CloseNotify()")
		return
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Println("Error: ResponseWriter does not support Flush()")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Cache-Control", "no-cache")
	for {
		// Check to see if the connection has closed
		select {
		case <-closeNotifier.CloseNotify():
			log.Printf("%s disconnected", r.RemoteAddr)
			return
		default:
		}

		// Marshal a Message into JSON
		message := <-messages
		bs, err := json.Marshal(message)
		if err != nil {
			log.Fatal("Marhsal: ", err)
		}

		// Put into the SSE format:
		// https://html.spec.whatwg.org/multipage/comms.html#the-eventsource-interface
		bs = append([]byte("data: "), bs...)
		bs = append(bs, []byte("\n\n")...)

		// Write and flush
		_, err = w.Write(bs)
		if err != nil {
			log.Printf("Write to %s failed", r.RemoteAddr)
			return
		}
		flusher.Flush()
		log.Printf("Sent message #%d to %s", message.ID, r.RemoteAddr)
	}
}

func main() {
	http.HandleFunc("/messages", MessageServer)
	log.Println("Listening on http://0.0.0.0:8080/messages")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
