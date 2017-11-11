package main

var (
	users = make(map[string]*User)
)

type User struct {
	Undelivered map[string][]Message
	Delivered   map[string][]Message
}
