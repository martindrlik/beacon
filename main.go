package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	go start()
	http.HandleFunc("/listen", decor(listen))
	http.HandleFunc("/put", decor(put))
	err := http.ListenAndServe("127.0.0.1:8080", nil)
	log.Fatal(err)
}

func start() {
	for {
		select {
		case r := <-recvch:
			ch, ok := recvrs[r.key]
			if !ok {
				ch = make(chan string)
				recvrs[r.key] = ch
			}
			r.ch <- ch
		case m := <-msgch:
			ch, ok := recvrs[m.key]
			if !ok {
				continue
			}
			ch <- m.text
		}
	}
}

func listen(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	recv := receiver{
		key: r.FormValue("key"),
		ch:  make(chan chan string),
	}
	select {
	case recvch <- recv:
	case <-ctx.Done():
		ise(ctx.Err())
	}
	var ch chan string
	select {
	case ch = <-recv.ch:
	case <-ctx.Done():
		ise(ctx.Err())
	}
	select {
	case m := <-ch:
		fmt.Fprint(w, m)
	case <-ctx.Done():
		ise(ctx.Err())
	}
}

func put(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	m := message{
		key:  r.FormValue("key"),
		text: r.FormValue("text"),
	}
	select {
	case msgch <- m:
		fmt.Fprint(w, "sent ok")
	case <-ctx.Done():
		ise(ctx.Err())
	}
}

type receiver struct {
	key string
	ch  chan chan string
}

type message struct {
	key  string
	text string
}

var (
	recvch = make(chan receiver)
	recvrs = make(map[string]chan string)

	msgch = make(chan message)
)

func decor(fn func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, time.Minute)
		defer cancel()

		path := r.URL.Path
		defer un(trace(path))
		defer func() {
			switch x := recover().(type) {
			case nil:
			case error:
				http.Error(
					w,
					x.Error(),
					http.StatusInternalServerError)
			}
		}()
		fn(w, r)
	}
}

func un(s string, start time.Time) {
	t := time.Now()
	elapsed := t.Sub(start)
	log.Printf("leaving %s (took %v)", s, elapsed)
}

func trace(s string) (string, time.Time) {
	start := time.Now()
	log.Printf("entering %s", s)
	return s, start
}

func ise(err error) {
	log.Print(err)
	panic(err)
}
