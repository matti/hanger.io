package main

import (
	"flag"
	"fmt"
	"github.com/dustin/go-broadcast"
	"github.com/go-redis/redis"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var (
	port          = flag.String("p", "8080", "HTTP Server port number")
	redisHost     = flag.String("h", "127.0.0.1", "Redis host")
	redisPort     = flag.Int("rp", 6379, "Redis port number")
	redisPassword = flag.String("pass", "", "Redis password")
	redisDBIndex  = flag.Int("db", 0, "Redis DB index")

	redisClient *redis.Client
	hangers     = make(map[string]broadcast.Broadcaster)
)

func main() {
	flag.Parse()
	redisClient = redis.NewClient(&redis.Options{
		Addr:     *redisHost + ":" + strconv.Itoa(*redisPort),
		Password: *redisPassword,
		DB:       *redisDBIndex,
	})

	fmt.Println("Starting HTTP server...")
	http.HandleFunc("/pause/", pause)
	http.HandleFunc("/continue/", cont)
	err := http.ListenAndServe("localhost:"+*port, nil)
	if err != nil {
		panic(err)
	}
}

func pause(w http.ResponseWriter, req *http.Request) {
	hangID := req.URL.Path[len("/pause/"):]
	if hangID == "" {
		fmt.Fprintln(w, "Provide an id: /pause/<number>")
	}

	if broadcaster, exists := hangers[hangID]; exists {
		ch := make(chan interface{})
		broadcaster.Register(ch)
		defer broadcaster.Unregister(ch)

		maxRampUp := <-ch
		sleepAndRespond(w, maxRampUp.(int), "done")
	} else {
		hangers[hangID] = broadcast.NewBroadcaster(5)
		broadcaster := hangers[hangID]

		pubsub := redisClient.Subscribe(hangID)
		ch := pubsub.Channel()

		for msg := range ch {
			maxRampUp, _ := strconv.Atoi(msg.Payload)
			broadcaster.Submit(maxRampUp)
			sleepAndRespond(w, maxRampUp, "done")
			delete(hangers, hangID)
			return
		}
	}
}

func cont(w http.ResponseWriter, req *http.Request) {
	hangID := req.URL.Path[len("/continue/"):]
	if hangID == "" {
		fmt.Fprintln(w, "Provide an id: /continue/<number>")
	}

	var rampUpTime string
	params := req.URL.Query()["rampup"]

	if len(params) == 0 {
		rampUpTime = "5"
	} else {
		rampUpTime = req.URL.Query()["rampup"][0]
	}

	err := redisClient.Publish(hangID, rampUpTime).Err()
	if err != nil {
		panic(err)
	}
}

func sleepAndRespond(w http.ResponseWriter, maxRampUp int, message string) {
	rampUp := time.Duration(rand.Intn(maxRampUp))
	time.Sleep(rampUp * time.Second)
	fmt.Fprint(w, "done")
}
