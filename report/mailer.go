package report

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

const MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

type Request struct {
	from    string
	to      []string
	subject string
	body    string
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
	subject := "Subject: " + r.subject + "\n"
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

func Mail(subject string, emailData []*TemplateData) {
	r := NewRequest([]string{"christiannebarry9@gmail.com"}, subject, "body")
	if err := r.ParseTemplate("template.html", emailData); err == nil {
		r.SendEmail()
		fmt.Printf("Email sent %s\n", subject)
	}

}
