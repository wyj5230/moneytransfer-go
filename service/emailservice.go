package service

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"os"
)

type MailboxConf struct {
	Title         string
	Body          string
	RecipientList []string
	Sender        string
	SPassword     string
	SMTPAddr      string
	SMTPPort      int
}

func sendEmail(subject string, content string, recipient string) {
	var mailConf MailboxConf
	mailConf.Title = subject
	mailConf.Body = content
	mailConf.RecipientList = []string{recipient}
	mailConf.Sender = "327113606@qq.com"
	mailConf.SPassword = os.Getenv("MONEYTRANSFER_EMAIL_SECRET")
	mailConf.SMTPAddr = `smtp.qq.com`
	mailConf.SMTPPort = 25

	m := gomail.NewMessage()
	m.SetHeader(`From`, mailConf.Sender)
	m.SetHeader(`To`, mailConf.RecipientList...)
	m.SetHeader(`Subject`, mailConf.Title)
	m.SetBody(`text/html`, mailConf.Body)
	err := gomail.NewDialer(mailConf.SMTPAddr, mailConf.SMTPPort, mailConf.Sender, mailConf.SPassword).DialAndSend(m)
	if err != nil {
		log.Fatalf("Send Email Failed, %s", err.Error())
		return
	}
	fmt.Printf("Send Email Successfully to: %s\n\n", recipient)
}
