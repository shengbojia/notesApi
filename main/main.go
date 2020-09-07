package main

import (
	"encoding/json"
	"github.com/shengbojia/gorouter"
	"net/http"
)

type Note struct {
	Id string `json:"id"`
	Title string `json:"title"`
	Body string `json:"body"`
}

var notes []Note

func getNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func getNote(w http.ResponseWriter, r *http.Request) {
	w.Header ().Set("Content-Type", "application/json")
	if idParam, ok := gorouter.GetParam(r, "id"); ok {
		for _, note := range notes {
			if note.Id == idParam {
				json.NewEncoder(w).Encode(note)
				break
			}
		}
	}
	json.NewEncoder(w).Encode(&Note{})
}

func main() {
	router := gorouter.New()

	http.ListenAndServe(":8000", router)
}
