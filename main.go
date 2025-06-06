package main

/*
$ curl http://localhost:9999/_gocache/scores/Tom
630

$ curl http://localhost:9999/_gocache/scores/kkk
kkk not exist
*/

import (
	"fmt"
	"gocache"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	gocache.NewGroup("scores", 2<<10, gocache.Getterfunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	peers := gocache.NewHTTPPool(addr)
	log.Println("gocache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
