package report

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"strings"
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
	host := os.Getenv("HOST")
	mailPort := os.Getenv("MAIL_PORT")
	fromEmail := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")

	if err := smtp.SendMail(host+":"+mailPort, smtp.PlainAuth("", fromEmail, password, host), fromEmail, r.to, msg); err != nil {
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

func MailReport(subject string, emailData reportTemplate) {
	toEmails := os.Getenv("DESTINATION")
	destination := strings.Split(toEmails, ",")
	r := NewRequest(destination, subject, "body")
	if err := r.ParseTemplate("report_template.html", emailData); err == nil {
		r.SendEmail()
		fmt.Printf("Email sent %s\n", subject)
	}
}

func MailNotification(subject string, emailData NotificationTemplate) {
	toEmails := os.Getenv("DESTINATION")
	destination := strings.Split(toEmails, ",")
	r := NewRequest(destination, subject, "body")
	if err := r.ParseTemplate("notification_template.html", emailData); err == nil {
		r.SendEmail()
		fmt.Printf("Email sent %s\n", subject)
	}
}
