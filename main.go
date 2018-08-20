// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")
var static = flag.String("files", "static", "path for static files")


var hubs = make(map[string]*Hub)

func main() {
	flag.Parse()
    fs := http.FileServer(http.Dir(*static))
    http.Handle("/", fs)

	http.HandleFunc("/hub", func(w http.ResponseWriter, r *http.Request) {
        hub, ok := hubs[r.URL.Path]
        if !ok {
            hub = newHub()
            hubs[r.URL.Path] = hub
            go hub.run()
        }
		serveWs(hub, w, r)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
