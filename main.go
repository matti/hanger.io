package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var waitingRequests = make(map[int]bool)

func hang(w http.ResponseWriter, req *http.Request) {
	hangID, err := strconv.Atoi(req.URL.Path[len("/hang/"):])
	if err != nil {
		fmt.Fprintln(w, "Provide an id: /hang/<number>")
	}

	waitingRequests[hangID] = true
	for waitingRequests[hangID] {
		time.Sleep(1 * time.Second)
		fmt.Println(waitingRequests[hangID])
	}
	delete(waitingRequests, hangID)
}

func unhang(w http.ResponseWriter, req *http.Request) {
	hangID, err := strconv.Atoi(req.URL.Path[len("/unhang/"):])
	if err != nil {
		fmt.Fprintln(w, "Provide an id: /unhang/<number>")
	}

	waitingRequests[hangID] = false
}

func main() {
	http.HandleFunc("/hang/", hang)
	http.HandleFunc("/unhang/", unhang)
	http.ListenAndServe("localhost:8080", nil)
}
