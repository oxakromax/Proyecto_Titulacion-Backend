package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"strings"
)

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown from server")
		}
	}
	return nil, nil
}

func SendMail(body, subject string, to []string) error {

	// Sender data.
	from := "monitordeprocesos@outlook.com"
	password := "Monitor123!"

	// smtp server configuration.
	smtpHost := "smtp-mail.outlook.com"
	smtpPort := "587"

	conn, err := net.Dial("tcp", "smtp-mail.outlook.com:587")
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return err
	}

	tlsconfig := &tls.Config{
		ServerName: smtpHost,
	}

	if err = c.StartTLS(tlsconfig); err != nil {
		return err
	}

	auth := LoginAuth(from, password)

	if err = c.Auth(auth); err != nil {
		return err
	}

	msg := []byte("To: " + strings.Join(to, ";") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\n")

	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("Email Sent!")
	return nil
}
