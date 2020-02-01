package main

// curl -d status=504 -d rampup=5

import (
	"flag"
	"fmt"
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

	pubsub := redisClient.Subscribe(hangID)
	ch := pubsub.Channel()

	for msg := range ch {
		maxRampUp, _ := strconv.Atoi(msg.Payload)
		rampUp := time.Duration(rand.Intn(maxRampUp))
		time.Sleep(rampUp * time.Second)
		fmt.Fprint(w, "done")
		return
	}
}

func cont(w http.ResponseWriter, req *http.Request) {
	rampUpTime := req.URL.Query()["rampup"][0]
	if rampUpTime == "" {
		rampUpTime = "5"
	}

	hangID := req.URL.Path[len("/continue/"):]
	if hangID == "" {
		fmt.Fprintln(w, "Provide an id: /continue/<number>")
	}

	err := redisClient.Publish(hangID, rampUpTime).Err()
	if err != nil {
		panic(err)
	}
}
