package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"

	redis "github.com/alphazero/Go-Redis"
)

const pinDuration int64 = 60 // how long the pin should be valid (in seconds)

func main() {
	spec := redis.DefaultSpec().Db(2)
	client, err := redis.NewSynchClientWithSpec(spec)
	if err != nil {
		log.Println("failed to create the client", err)
		return
	}

	http.HandleFunc("/getPin/", func(w http.ResponseWriter, r *http.Request) {
		out, err := getPin(client)
		if err != nil {
			log.Println(err)
		}

		io.WriteString(w, out)
	})

	http.HandleFunc("/checkPin/", func(w http.ResponseWriter, r *http.Request) {
		out, err := checkPin(client, r.URL.Query())
		if err != nil {
			log.Println(err)
		}

		io.WriteString(w, out)
	})

	http.ListenAndServe(":80", nil)
}

func getPin(client redis.Client) (string, error) {
	pin := strconv.Itoa(rangeIn(100000, 999999))

	if err := client.Set(pin, []byte("active")); err != nil {
		return "error: failed setting pin in the database", err
	}

	if _, err := client.Expire(pin, pinDuration); err != nil {
		return "error: failed setting expiratoin in the database", err
	}

	return pin, nil
}

func checkPin(client redis.Client, values url.Values) (string, error) {
	pin, ok := values["pin"]
	if !ok || len(pin[0]) < 1 {
		return "error: pin is missing", nil
	}

	exists, err := client.Exists(pin[0])
	if err != nil {
		return "error: failed getting pin from database", err
	}

	return strconv.FormatBool(exists), nil
}

func rangeIn(low, hi int) int {
	return low + rand.Intn(hi-low)
}
