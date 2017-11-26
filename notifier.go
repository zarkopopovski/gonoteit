package main

import gomail "gopkg.in/gomail.v2"

type Notifier struct {
	config      *Config
	userMail    string
	noteTitle   string
	noteBody    string
	dateCreated string
}

func (notifier *Notifier) sendNotification(notificationType string) (err error) {
	config := notifier.config

	go func() {
		m := gomail.NewMessage()

		m.SetHeader("From", config.mailUsername)
		m.SetHeader("To", notifier.userMail)
		m.SetHeader("Subject", "goNoteIT: System notification")

		if notificationType == "create" {
			m.SetBody("text/html", "Hi, you have just created new note <br/>"+notifier.noteTitle+"<br/>"+notifier.noteBody)
		} else if notificationType == "edit" {
			m.SetBody("text/html", "Hi, your note <b>"+notifier.noteTitle+"</b> is updated.")
		} else if notificationType == "delete" {
			m.SetBody("text/html", "Hi, your note <b>"+notifier.noteTitle+"</b> created on "+notifier.dateCreated+" is deleted. Bellow you can see the note<br/>"+notifier.noteBody)
		}

		d := gomail.NewPlainDialer(config.mailServer, config.mailPort, config.mailUsername, config.mailPassword)

		if err := d.DialAndSend(m); err != nil {
			panic(err)
		}
		return
	}()

	return
}
