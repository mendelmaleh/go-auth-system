package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	redis "github.com/alphazero/Go-Redis"
)

var spec *redis.ConnectionSpec
var client redis.Client
var e error

var pinDuration int64 = 60 // How many time the pin should be valid (in seconds)

func main() {

	spec = redis.DefaultSpec().Db(2)
	client, e = redis.NewSynchClientWithSpec(spec)
	if e != nil {
		fmt.Println("Failed to create the client", e)
		return
	}

	http.HandleFunc("/checkPin", checkPin)
	http.HandleFunc("/getPin", getPin)

	http.ListenAndServe(":80", nil)
}

func rangeIn(low, hi int) int {
	return low + rand.Intn(hi-low)
}

func checkPin(w http.ResponseWriter, req *http.Request) {

	pin, ok := req.URL.Query()["pin"]
	if !ok || len(pin[0]) < 1 {
		fmt.Fprintf(w, "Error: pin is missing.")
		return
	}

	ex, e := client.Exists(pin[0])
	if e != nil {
		fmt.Fprintf(w, "Error: failed getting pin from DB.")
		return
	}

	if ex {
		fmt.Fprintf(w, "true")
	} else {
		fmt.Fprintf(w, "false")
	}
}

func getPin(w http.ResponseWriter, req *http.Request) {
	pin := strconv.Itoa(rangeIn(100000, 999999))
	e := client.Set(pin, []byte("active"))
	client.Expire(pin, pinDuration)
	if e != nil {
		fmt.Fprintf(w, "Error: failed setting pin in the DB.")
		return
	}
	fmt.Fprintf(w, pin)
}
