package main

import (
	"net"
	"net/http"

	"github.com/rs/zerolog"
)

func addLogData(r *http.Request, e *zerolog.Event, msg string) {
	headersMap := make(map[string][]string)
	for name, values := range r.Header {
		headersMap[name] = values // Taking the first value for simplicity
	}

	e.Str("method", r.Method)
	e.Str("path", r.URL.Path)
	e.IPAddr("IP", net.IP(r.RemoteAddr))
	e.Fields(headersMap)

	e.Msg(msg)
}
