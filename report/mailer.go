package report

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"time"
)

const MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

type Mailer struct {
	inCh chan *StatusReport
}

func NewMailer(inCh chan *StatusReport) *Mailer {
	mailer := new(Mailer)
	mailer.inCh = inCh
	return mailer
}

func (mailer *Mailer) Start() {
	go func() {
		mailer.Mail()
		time.Sleep(10 * time.Minute)
	}()
}

func NewRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}
func (r *Request) SendEmail() (bool, error) {
	mime := MIME
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(subject + mime + "\n" + r.body)
	addr := "smtp.gmail.com:587"

	if err := smtp.SendMail(addr, smtp.PlainAuth("", "testbedmonitorreports@gmail.com", "hworxdpgshqfozdv", "smtp.gmail.com"), "testbedmonitorreports@gmail.com", r.to, msg); err != nil {
		fmt.Printf("Error sending email: %s", err)
		return false, err
	}
	return true, nil
}

func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		fmt.Printf("Error at ParseFiles %s", err)
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		fmt.Printf("Error at Execute %s", err)
		return err
	}
	r.body = buf.String()
	return nil
}

func (mailer *Mailer) Mail() {
	templateData := struct {
		Name string
		URL  string
	}{
		Name: "CB",
		URL:  "URL",
	}
	r := NewRequest([]string{"christiannebarry9@gmail.com"}, "Hello World!", "Hello, World!")
	if err := r.ParseTemplate("template.html", templateData); err == nil {
		ok, _ := r.SendEmail()
		fmt.Printf("Email sent %s", ok)
	}

}
