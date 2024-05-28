package main

import (
	"fmt"
	"net/smtp"
)

type message struct {
	Subject   string
	Body      string
	Recipient string
	Sender    string
	ReplyTo   string
}

func buildAndSend(cfg *Config, sub FormSubmission) error {
	// Build the email message using a template
	message, err := buildEmailMessage(sub)
	if err != nil {
		fmt.Println("Failed to build email message:", err)
		return err
	}

	// Send the email using the configured SMTP server
	err = sendEmail(cfg, message)
	if err != nil {
		fmt.Println("Failed to send email:", err)
		return err
	}

	fmt.Println("Email sent successfully!")
	return nil
}

func buildEmailMessage(sub FormSubmission) (message, error) {
	body := fmt.Sprintf("From: %s\nReply-To: %s\n\n%s", sub.Fields[sub.FormCfg.Field.Email], sub.Fields[sub.FormCfg.Field.Email], sub.Fields[sub.FormCfg.Field.Message])

	return message{
		Subject:   sub.FormCfg.Mail.Subject + " - " + sub.Id,
		Body:      body,
		Recipient: sub.FormCfg.Mail.Recipient,
		Sender:    sub.FormCfg.Mail.Sender,
		ReplyTo:   sub.Fields[sub.FormCfg.Field.Email],
	}, nil
}

func sendEmail(cfg *Config, message message) error {
	auth := smtp.PlainAuth("", cfg.Smtp.User, cfg.Smtp.Password, cfg.Smtp.Host)
	err := smtp.SendMail(fmt.Sprintf("%s:%d", cfg.Smtp.Host, cfg.Smtp.Port), auth, message.Sender, []string{message.Recipient}, []byte(message.Body))
	if err != nil {
		return err
	}

	return nil
}
