package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

type UserController struct {
	dbConnection *DBConnection
}

func (uController *UserController) registerNewUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	sha1Hash := sha1.New()
	sha1Hash.Write([]byte(password))
	sha1HashString := sha1Hash.Sum(nil)

	passwordEnc := fmt.Sprintf("%x", sha1HashString)

	sha1Hash = sha1.New()
	sha1Hash.Write([]byte(time.Now().String() + username + password + email))
	sha1HashString = sha1Hash.Sum(nil)

	userID := fmt.Sprintf("%x", sha1HashString)

	query := "INSERT INTO users(id, username, email, password) VALUES('" + userID + "','" + username + "','" + email + "','" + passwordEnc + "')"

	_, err := uController.dbConnection.db.Exec(query)

	if err == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		user := &UserModel{
			Id:       userID,
			UserName: username,
			Email:    email,
		}

		if err := json.NewEncoder(w).Encode(user); err != nil {
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

func (uController *UserController) checkUserCredentials(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	sha1Hash := sha1.New()
	sha1Hash.Write([]byte(password))
	sha1HashString := sha1Hash.Sum(nil)

	passwordEnc := fmt.Sprintf("%x", sha1HashString)

	query := "SELECT id, username, email FROM users WHERE username='" + username + "' AND password = '" + passwordEnc + "'"

	newUser := new(UserModel)

	err := uController.dbConnection.db.QueryRow(query).Scan(
		&newUser.Id,
		&newUser.UserName,
		&newUser.Email)

	if err == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(newUser); err != nil {
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

func (uController *UserController) sendMailWithResetLink(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

}

func (uController *UserController) updateUserProfile(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

}
