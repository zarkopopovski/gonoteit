package main

import (
	//"flag"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"strconv"

	"github.com/knq/ini"
)

type GNoteAPI struct {
	dbConnection *DBConnection
	uController  *UserController
	nController  *NoteController
	config       *Config
}

type Config struct {
	mailServer   string
	mailPort     int
	mailUsername string
	mailPassword string
}

func CreateGNoteAPI(config *Config) *GNoteAPI {
	apiHandlers := &GNoteAPI{
		dbConnection: CreateDBConnection(),
		uController:  &UserController{},
		nController:  &NoteController{},
		config:       config,
	}

	apiHandlers.uController.dbConnection = apiHandlers.dbConnection
	apiHandlers.uController.config = apiHandlers.config

	apiHandlers.nController.dbConnection = apiHandlers.dbConnection
	apiHandlers.nController.config = apiHandlers.config

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

	fileCfg, err := ini.LoadFile("config.cfg")
	if err != nil {
		log.Fatal("Error with service configuration %s", err)
	}

	port := fileCfg.GetKey("service-1.port")

	if port == "" {
		log.Fatal("Error with port number configuration")
	}

	serverPort := fileCfg.GetKey("service-1.mailport")
	serverPortI, _ := strconv.Atoi(serverPort)

	config := &Config{
		mailServer:   fileCfg.GetKey("service-1.mailserver"),
		mailPort:     serverPortI,
		mailUsername: fileCfg.GetKey("service-1.mailusername"),
		mailPassword: fileCfg.GetKey("service-1.mailpassword"),
	}

	gnoteAPI := CreateGNoteAPI(config)
	router := CreateNewRouter(gnoteAPI)

	log.Fatal(http.ListenAndServe(":"+port, router))
}
