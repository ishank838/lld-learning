package main

import "log"

type NOTIFICATION_TYPE int

const (
	EMAIL_NOTIFICATION NOTIFICATION_TYPE = iota
	PHONE_NOTIFICATION
)

type Notification interface {
	Send(message string)
}

type EmailNotification struct {
	//internal Deps
}

func (e *EmailNotification) Send(message string) {
	log.Println("sending email for ", message)
}

type PhoneNotification struct {
	//internal Deps
}

func (e *PhoneNotification) Send(message string) {
	log.Println("sending phone message for ", message)
}

func NewNotification(kind NOTIFICATION_TYPE) Notification {

	switch kind {
	case EMAIL_NOTIFICATION:
		return &EmailNotification{}
	case PHONE_NOTIFICATION:
		return &PhoneNotification{}
	}

	return &PhoneNotification{}
}

func main() {
	email := NewNotification(EMAIL_NOTIFICATION)

	email.Send("test message")

	phone := NewNotification(PHONE_NOTIFICATION)

	phone.Send("test message")
}
