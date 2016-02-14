package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Location struct {
	Device    string `json:"device"`
	Date      string `json:"date"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

var m map[string]*Location = make(map[string]*Location, 0)

func main() {
	http.HandleFunc("/", locationHandler)

	p := os.Getenv("PORT")
	http.ListenAndServe(":"+p, nil)
}

func locationHandler(w http.ResponseWriter, r *http.Request) {
	// check URL parameter
	d := r.URL.Path[1:]

	// return location for device name from URL
	if d != "" {
		returnLocation(d, w, r)
		return
	}

	// store new location
	receiveLocation(w, r)
	return
}

func returnLocation(d string, w http.ResponseWriter, r *http.Request) {
	// find location in map
	l := m[d]
	if l == nil {
		log.Println("location not found for", d)
		http.NotFound(w, r)
		return
	}

	// convert location to json
	js, err := json.Marshal(l)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return location as json
	log.Println("returning location", string(js))
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func receiveLocation(w http.ResponseWriter, r *http.Request) {
	l := &Location{
		Device:    r.FormValue("username"),
		Date:      r.FormValue("date"),
		Latitude:  r.FormValue("latitude"),
		Longitude: r.FormValue("longitude"),
	}

	if l.Date == "" || l.Device == "" || l.Latitude == "" || l.Longitude == "" {
		msg := "invalid parameters given"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// store location in map
	log.Println("setting location to", *l)
	m[l.Device] = l

}
