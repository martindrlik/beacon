package main

import "log"

func un(s string) {
	log.Printf("%s ended", s)
}

func trace(s string) string {
	log.Printf("%s started", s)
	return s
}
