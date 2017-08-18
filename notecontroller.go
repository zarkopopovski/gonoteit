package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"

	"io/ioutil"

	"git.cerebralab.com/george/logo"
)

type NoteController struct {
	dbConnection *DBConnection
}

func (nController *NoteController) index(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	index, err := ioutil.ReadFile("./web/index.html")

	logo.RuntimeError(err)

	if err != nil {
		logo.RuntimeError(err)
		return
	}

	fmt.Fprintf(w, string(index))
}

func (nController *NoteController) ShowFavIcon(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	http.ServeFile(w, r, "./web/favicon.png")
}

func (nController *NoteController) saveUserNote(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userID := r.FormValue("user_id")
	noteID := r.FormValue("note_id")
	title := r.FormValue("title")
	body := r.FormValue("body")
	//dateCreated := r.FormValue("date_created")

	query := ""

	if noteID == "" {
		sha1Hash := sha1.New()
		sha1Hash.Write([]byte(time.Now().String() + title + body + userID))
		sha1HashString := sha1Hash.Sum(nil)

		noteID = fmt.Sprintf("%x", sha1HashString)

		query = "INSERT INTO notes(id, user_id, title, body, date_created) VALUES($1, $2, $3, $4, datetime('now'))"
		_, err := nController.dbConnection.db.Exec(query, noteID, userID, title, body)

		if err == nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)

			if err := json.NewEncoder(w).Encode(map[string]string{"status": "success", "error_code": "0", "note_id": noteID}); err != nil {
				panic(err)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	} else {
		query = "UPDATE notes SET title=$1, body=$2 WHERE id=$3"
		_, err := nController.dbConnection.db.Exec(query, title, body, noteID)

		if err == nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)

			if err := json.NewEncoder(w).Encode(map[string]string{"status": "success", "error_code": "0", "note_id": noteID}); err != nil {
				panic(err)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

}

func (nController *NoteController) deleteUserNote(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	noteID := r.FormValue("note_id")

	query := "DELETE FROM notes WHERE id=$1"

	_, err := nController.dbConnection.db.Exec(query, noteID)

	if err == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(map[string]string{"status": "success", "error_code": "0"}); err != nil {
			panic(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(err); err != nil {
		panic(err)
	}
}

func (nController *NoteController) listAllUserNotes(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userID := r.FormValue("user_id")

	query := "SELECT id, title, body, date_created FROM notes WHERE user_id=$1 ORDER BY date_created DESC"

	rows, err := nController.dbConnection.db.Query(query, userID)

	if err == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		notes := make([]*NoteModel, 0)

		for rows.Next() {

			newNote := new(NoteModel)

			_ = rows.Scan(
				&newNote.Id,
				&newNote.Title,
				&newNote.Body,
				&newNote.Date)

			notes = append(notes, newNote)

		}

		if err := json.NewEncoder(w).Encode(notes); err != nil {
			panic(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(err); err != nil {
		panic(err)
	}

}
