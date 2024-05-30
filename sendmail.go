package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

type message struct {
	Subject   string
	Body      []byte
	Recipient string
	Sender    string
	ReplyTo   string
}

func buildAndSend(cfg *Config, sub FormSubmission) error {
	message, err := buildEmailMessage(sub)
	if err != nil {
		fmt.Println("Failed to build email message:", err)
		return err
	}

	err = sendEmail(cfg, message)
	if err != nil {
		fmt.Println("Failed to send email:", err)
		return err
	}

	fmt.Println("Email sent successfully!")
	return nil
}

func buildEmailMessage(sub FormSubmission) (message, error) {
	tmpl := loadTemplate(sub.Id)

	var body bytes.Buffer
	if err := tmpl.ExecuteTemplate(&body, sub.Id+".html", sub.Body); err != nil {
		return message{}, err
	}

	headers := "From: <" + sub.FormCfg.Mail.Sender + ">\r\n" +
		"To: <" + sub.FormCfg.Mail.Recipient + ">\r\n" +
		"Subject: " + sub.FormCfg.Mail.Subject + " - " + sub.Id + "\r\n" +
		"Reply-To: <" + sub.Body[sub.FormCfg.Fields.Email] + ">\r\n" +
		"MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"

	return message{
		Subject:   sub.FormCfg.Mail.Subject + " - " + sub.Id,
		Body:      []byte(headers + body.String()),
		Recipient: "<" + sub.FormCfg.Mail.Recipient + ">",
		Sender:    "<" + sub.FormCfg.Mail.Sender + ">",
		ReplyTo:   sub.Body[sub.FormCfg.Fields.Email],
	}, nil
}

func sendEmail(cfg *Config, message message) error {
	auth := smtp.PlainAuth("", cfg.Smtp.User, cfg.Smtp.Password, cfg.Smtp.Host)
	err := smtp.SendMail(fmt.Sprintf("%s:%d", cfg.Smtp.Host, cfg.Smtp.Port), auth, message.Sender, []string{message.Recipient}, message.Body)
	if err != nil {
		return err
	}

	return nil
}

func loadTemplate(id string) *template.Template {
	defaultTemplate, err := template.New("default").ParseFiles("forms/default.html")
	if err != nil {
		fmt.Println("Failed to parse default template:", err)
		return nil
	}

	template, err := template.ParseFiles(fmt.Sprintf("forms/%s.html", id))
	if err != nil {
		fmt.Println("Failed to parse template:", err)
		return defaultTemplate
	}

	return template
}
