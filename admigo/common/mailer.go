package common

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type eMessage struct {
	From    string
	Subject string
	To      string
	Email   string
	Key     string
	Link    string
}

func SendEmailUserRegistration(email string) {
	e := Env()
	message := &eMessage{
		From:    e.Mail.From,
		To:      email,
		Subject: "User Registration on Admigo",
		Email:   email,
		Key:     Encrypt(email),
		Link:    fmt.Sprintf("%s:%d", e.Mail.GotoUrl, e.Port),
	}
	var body bytes.Buffer
	GenerateMail(&body, message, "layout", "reguser")

	auth := smtp.PlainAuth("",
		e.Mail.Username,
		e.Mail.Password,
		e.Mail.Host,
	)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         e.Mail.Host,
	}

	servername := fmt.Sprintf("%s:%d", e.Mail.Host, e.Mail.Port)

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		Danger("mailer.go tls.Dial", err)
	}

	c, err := smtp.NewClient(conn, e.Mail.Host)
	if err != nil {
		Danger("mailer.go smtp.NewClient", err)
	}

	if err = c.Auth(auth); err != nil {
		Danger("mailer.go c.Auth", err)
	}

	// To && From
	if err = c.Mail(message.From); err != nil {
		Danger("mailer.go c.Mail", err)
	}

	if err = c.Rcpt(message.To); err != nil {
		Danger("mailer.go c.Rcpt", err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		Danger("mailer.go c.Data()", err)
	}

	_, err = w.Write(body.Bytes())
	if err != nil {
		Danger("mailer.go w.Write", err)
	}

	err = w.Close()
	if err != nil {
		Danger("mailer.go w.Close", err)
	}

	c.Quit()
}
