package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/rs/zerolog"
)

func addLogData(r *http.Request, e *zerolog.Event, msg string) {
	for name, values := range r.Header {
		for i, v := range values {
			oname := fmt.Sprintf("%s-%d", name, i)
			e.Str(oname, v)
		}
	}

	e.Str("method", r.Method)
	e.Str("path", r.URL.Path)
	e.IPAddr("IP", net.IP(r.RemoteAddr))

	e.Msg(msg)
}
