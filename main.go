package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/get", get)
	http.HandleFunc("/put", put)
	err := http.ListenAndServe("127.0.0.1:8080", nil)
	log.Fatal(err)
}

func get(w http.ResponseWriter, r *http.Request) {
	defer un(trace("get"))

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	select {
	case rm := <-messageCh:
		fmt.Fprintln(w, rm.message.Text)
	case <-ctx.Done():
		err := ctx.Err()
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func put(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer un(trace("put"))

	sender := r.FormValue("sender")
	recv := r.FormValue("receiver")
	m := Message{
		Text: r.FormValue("text"),
		Time: time.Now(),
	}

	rm := rawMessage{
		sender:   sender,
		receiver: recv,
		message:  m,
		response: make(chan string),
	}
	messageCh <- rm
	select {
	case res := <-rm.response:
		fmt.Fprintln(w, res)
	case <-ctx.Done():
		err := ctx.Err()
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
