package report

import (
	"bytes"
	"html/template"
	"log"
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

type NotificationTemplate struct {
	TowerIP   string
	Field     string
	Value     string
	Timestamp string
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
	defer func() {
		if rec := recover(); rec != nil {
			log.Println("Unable to send email. Recovered.")
		}
	}()
	if err := smtp.SendMail(host+":"+mailPort, smtp.PlainAuth("", fromEmail, password, host), fromEmail, r.to, msg); err != nil {
		log.Panicf("Error sending email: %s", err)
		return false, err
	}
	return true, nil
}

func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	defer func() {
		if rec := recover(); rec != nil {
			log.Println("Unable to parse template. Recovered.")
		}
	}()
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		log.Panicf("Error parsing email template %s: %s\n", templateFileName, err)
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		log.Panicf("Error exectuing email template: %s\n", err)
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
		log.Printf("Email sent %s\n", subject)
	}
}

func MailNotification(subject string, emailData NotificationTemplate) {
	toEmails := os.Getenv("DESTINATION")
	destination := strings.Split(toEmails, ",")
	r := NewRequest(destination, subject, "body")
	if err := r.ParseTemplate("notification_template.html", emailData); err == nil {
		r.SendEmail()
		log.Printf("Email sent %s\n", subject)
	}
}
