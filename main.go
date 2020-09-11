package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/shengbojia/gorouter"
	"io/ioutil"
	"net/http"
)

var db *sql.DB

type Note struct {
	Id string `json:"id"`
	Title string `json:"title"`
	Body string `json:"body"`
	WrittenOn string `json:"written_on"`
}

func getNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var notes []Note

	result, err := db.Query("SELECT * FROM notes")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "database query error"}`))
		return
	}

	defer result.Close()

	for result.Next() {
		var note Note
		err := result.Scan(&note.Id, &note.Title, &note.Body, &note.WrittenOn)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "database data error"}`))
			return
		}
		notes = append(notes, note)
	}

	json.NewEncoder(w).Encode(notes)
}

func getNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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

func createNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	stmt, err := db.Prepare("INSERT INTO notes VALUES(?, ?, ?, ?)")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "error preparing database query"}`))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "error reading request body"}`))
		return
	}

	keyVal := make(map[string]string)
	err := json.Unmarshal(body, &keyVal)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "error unmarshalling"}`))
		return
	}
	id := keyVal["id"]
	title := keyVal["title"]
	noteBody := keyVal["body"]
	writtenOn := keyVal["written_on"]

	result, err := stmt.Exec(id, title, noteBody, writtenOn)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "error executing database query"}`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"success": "created"}"`))
}

func main() {
	var err error

	db, err = sql.Open("mysql", "root:1234abcd@tcp(127.0.01:3306)/notes_app")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	router := gorouter.New()
	http.ListenAndServe(":8000", router)
}


