package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Status struct {
    Water int `json:"water"`
    Wind  int `json:"wind"`
}

type Response struct {
    Status Status `json:"status"`
}

var (
    mutex   sync.Mutex
    // counter int
)

func main() {
    go updateJSONPeriodically()

    http.HandleFunc("/", getStatus)
    http.ListenAndServe(":8080", nil)
}

func updateJSONPeriodically() {
    for {
        time.Sleep(15 * time.Second)
        updateJSON()
    }
}

func updateJSON() {
    mutex.Lock()
    defer mutex.Unlock()

    status := Status{
        Water: rand.Intn(100) + 1,
        Wind:  rand.Intn(100) + 1,
    }

    response := Response{
        Status: status,
    }

    data, err := json.MarshalIndent(response, "", "    ")
    if err != nil {
        panic(err)
    }

    err = ioutil.WriteFile("status.json", data, 0644)
    if err != nil {
        panic(err)
    }
}

func getStatus(w http.ResponseWriter, r *http.Request) {
    mutex.Lock()
    defer mutex.Unlock()

    data, err := ioutil.ReadFile("status.json")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    var response Response
    err = json.Unmarshal(data, &response)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(response)
}
