package main

import (
	//"flag"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type GNoteAPI struct {
	dbConnection *DBConnection
	uController  *UserController
	nController  *NoteController
}

func CreateGNoteAPI() *GNoteAPI {
	apiHandlers := &GNoteAPI{
		dbConnection: CreateDBConnection(),
		uController:  &UserController{},
		nController:  &NoteController{},
	}

	apiHandlers.uController.dbConnection = apiHandlers.dbConnection
	apiHandlers.nController.dbConnection = apiHandlers.dbConnection

	return apiHandlers
}

func CreateNewRouter(handlers *GNoteAPI) *httprouter.Router {
	router := httprouter.New()

	router.GET("/", handlers.nController.index)
	router.POST("/save_note", handlers.nController.saveUserNote)
	router.POST("/delete_note", handlers.nController.deleteUserNote)
	router.POST("/list_notes", handlers.nController.listAllUserNotes)
	router.POST("/register_user", handlers.uController.registerNewUser)
	router.POST("/login_user", handlers.uController.checkUserCredentials)
	router.POST("/reset_password", handlers.uController.sendMailWithResetLink)
	router.POST("/update_user", handlers.uController.updateUserProfile)

	return router
}

func main() {
	gnoteAPI := CreateGNoteAPI()
	router := CreateNewRouter(gnoteAPI)

	log.Fatal(http.ListenAndServe(":8080", router))
}
