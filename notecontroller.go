package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"io/ioutil"
)

type NoteController struct {
	dbConnection *DBConnection
	config       *Config
}

func (nController *NoteController) index(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	index, err := ioutil.ReadFile("./web/index.html")

	//panic(err)

	if err != nil {
		panic(err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Fprintf(w, string(index))
}

func (nController *NoteController) ShowFavIcon(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "image/png")

	http.ServeFile(w, r, "./web/favicon.png")
}

func (nController *NoteController) saveUserNote(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userID := r.FormValue("user_id")
	noteID := r.FormValue("note_id")
	title := r.FormValue("title")
	body := r.FormValue("body")

	config := nController.config

	userEmail := ""
	userQuery := "SELECT email FROM users WHERE id=$1"
	e := nController.dbConnection.db.QueryRow(userQuery, userID).Scan(&userEmail)

	if e != nil {
	}

	query := ""

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

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

			if userEmail != "" {
				notifier := &Notifier{
					config:    config,
					userMail:  userEmail,
					noteTitle: title,
					noteBody:  body,
				}
				notifier.sendNotification("create")
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

			if userEmail != "" {
				notifier := &Notifier{
					config:    config,
					userMail:  userEmail,
					noteTitle: title,
				}
				notifier.sendNotification("edit")
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

	noteTitle := ""
	noteBody := ""
	dateCreated := ""

	userEmail := ""

	noteQuery := "SELECT n.title, n.body, n.date_created, u.email FROM notes n INNER JOIN users u ON n.id=$1 AND n.user_id=u.id"
	_ = nController.dbConnection.db.QueryRow(noteQuery, noteID).Scan(&noteTitle, &noteBody, &dateCreated, &userEmail)

	config := nController.config

	query := "DELETE FROM notes WHERE id=$1"

	_, err := nController.dbConnection.db.Exec(query, noteID)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if err == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(map[string]string{"status": "success", "error_code": "0"}); err != nil {
			panic(err)
		}

		if userEmail != "" {
			notifier := &Notifier{
				config:      config,
				userMail:    userEmail,
				noteTitle:   noteTitle,
				noteBody:    noteBody,
				dateCreated: dateCreated,
			}
			notifier.sendNotification("delete")
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

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

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
