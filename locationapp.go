package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

type Location struct {
	Speed          string `json:"speed"`
	Eventtype      string `json:"eventtype"`
	Locationmethod string `json:"locationmethod"`
	Username       string `json:"username"`
	Date           string `json:"date"`
	Distance       string `json:"distance"`
	Phonenumber    string `json:"phonenumber"`
	Sessionid      string `json:"sessionid"`
	Direction      string `json:"direction"`
	Latitude       string `json:"latitude"`
	Longitude      string `json:"longitude"`
	Accuracy       string `json:"accuracy"`
	Extrainfo      string `json:"extrainfo"`
}

var lm map[string]*Location = make(map[string]*Location, 0)

func main() {
	http.HandleFunc("/", locationHandler)

	p := os.Getenv("PORT")
	http.ListenAndServe(":"+p, nil)
}

func locationHandler(w http.ResponseWriter, r *http.Request) {
	device := r.URL.Path[1:]
	if device != "" {
		returnLocation(w, r, device)
		return
	}

	receiveLocation(w, r)
	return
}

func returnLocation(w http.ResponseWriter, r *http.Request, device string) {
	location := lm[device]
	if location == nil {
		log.Println("No location saved for device", device)
		http.NotFound(w, r)
		return
	}

	js, err := json.Marshal(location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Returning location", string(js))
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func receiveLocation(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	// map[speed:[0] eventtype:[android] locationmethod:[fused] username:[hb] date:[2016-02-11+23%3A26%3A29] distance:[0.0] phonenumber:[51ce74b7-0780-43b7-85b4-9994e33cffe2] sessionid:[32a40c06-6be9-4cb8-b27f-19c07995deff] direction:[0] latitude:[51.3691482] accuracy:[71] extrainfo:[0] longitude:[12.7449269]]

	l := &Location{
		Speed:          strings.Join(r.Form["speed"], ""),
		Eventtype:      strings.Join(r.Form["eventtype"], ""),
		Locationmethod: strings.Join(r.Form["locationmethod"], ""),
		Username:       strings.Join(r.Form["username"], ""),
		Date:           strings.Join(r.Form["date"], ""),
		Distance:       strings.Join(r.Form["distance"], ""),
		Phonenumber:    strings.Join(r.Form["phonenumber"], ""),
		Sessionid:      strings.Join(r.Form["sessionid"], ""),
		Direction:      strings.Join(r.Form["direction"], ""),
		Latitude:       strings.Join(r.Form["latitude"], ""),
		Longitude:      strings.Join(r.Form["longitude"], ""),
		Accuracy:       strings.Join(r.Form["accuracy"], ""),
		Extrainfo:      strings.Join(r.Form["extrainfo"], ""),
	}

	lm[l.Username] = l

	log.Println("Set location to", l)
}
