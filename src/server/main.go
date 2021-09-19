package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type counters struct {
	sync.Mutex
	view  int
	click int
}

type counters_key struct {
	sync.Mutex 
	contents string 
	time time.Time
}

var (
	c = counters{}

	content = []string{"sports", "entertainment", "business", "education"}
	
	//Store to keep counter after regular interval, this can be accessed by content type and time 
	counter_store map[counters_key]*counters
	
	last_refresh_time = time.Now().Unix()
	time_window = 10
	rate_limit = 100
)

//Create counterMap to store view and clicks corresponding to each content type
func createCounterMaps() {
	// Data struncture to support event tracking (views and clicks) by content selection
	counter_map = map[string]counters{
		"sports":        counters{},
		"entertainment": counters{},
		"business":      counters{},
		"education":     counters{},
	}
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to EQ Works 😎")
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	data := content[rand.Intn(len(content))]
	c = &contentMap[data]

	c.Lock()
	c.view++
	c.Unlock()

	err := processRequest(r)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	// simulate random click call
	if rand.Intn(100) < 50 {
		processClick(data)
	}
}

func processRequest(r *http.Request) error {
	time.Sleep(time.Duration(rand.Int31n(50)) * time.Millisecond)
	return nil
}

func processClick(data string) error {
	c.Lock()
	c.click++
	c.Unlock()

	return nil
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	if !isAllowed() {
		w.WriteHeader(429)
		return
	}
}

func isAllowed() bool {
	// apply rate limit on basis of time_window
	if refresh_time < (time.Now().Unix() - time_window) {
		refresh_time = time.Now().Unix()
		statusRequests = 1
		return true
	}
	//apply rate limit on basis of number of requests
	if statusRequests < rate_limit {
		statusRequests++
		return true
	}
	return false
}

func uploadCounters() error {
	for true {
		time.Sleep(5000 * time.Millisecond)
		var i = 0
		for i < 4 {
			key := counters_key{contents: content[i], time: time.Now()}
			counterStore[key] = counter_map[content[i]]
			i++
		}
		createCounterMaps()
	}
	return nil
}

func main() {
	go uploadCounters()
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/stats/", statsHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
