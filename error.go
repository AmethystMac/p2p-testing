package main

import "log"

func handleError(e error) {
	log.Panicln(e)
}