package helpers

import (
	"log"
	"net/http"
)

func LogNewRequest(r *http.Request) {
	log.Printf("%sREQUEST:%s Addr: %s    url: %s    method: %s",
		Green, Reset,
		r.RemoteAddr,
		r.URL,
		r.Method,
	)
}

func LogError(e string) {
	log.Printf("%sERROR:%s %s",
		Red, Reset,
		e,
	)
}

// ANSI color codes
const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
)
