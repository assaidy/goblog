package log

import (
	"log"
	"net/http"
)

func NewRequest(r *http.Request) {
	log.Printf("%sREQUEST:%s Addr: %s    url: %s    method: %s",
		Green, Reset,
		r.RemoteAddr,
		r.URL,
		r.Method,
	)
}

func Error(e string, s http.ConnState) {
	log.Printf("%sERROR:%s %s    sent status code: %d",
		Red, Reset,
		e,
		s,
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
