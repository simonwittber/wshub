// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"hash/fnv"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")
var static = flag.String("files", "static", "path for static files")

var hubs = make(map[uint32]*Hub)
var noClients = make(chan uint32)

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func cleanup() {
	for {
		deadKey := <-noClients
		delete(hubs, deadKey)
	}
}

func main() {
	flag.Parse()
	fs := http.FileServer(http.Dir(*static))
	go cleanup()
	http.Handle("/", fs)
	http.HandleFunc("/hub/", func(w http.ResponseWriter, r *http.Request) {
		key := hash(r.URL.Path)
		hub, ok := hubs[key]
		isServer := false
		if !ok {
			hub = newHub(key, noClients)
			hubs[key] = hub
			isServer = true
			go hub.run()
		}
		serveWs(hub, w, r, isServer)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
