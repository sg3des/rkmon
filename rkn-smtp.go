package main

import (
	"bytes"
	"html/template"
	"net"
	"net/smtp"
	"strings"
	"time"
)

type SMTP struct {
	Host string
	User string
	Pass string
}

func NewSMTP(host, user, pass string) *SMTP {
	return &SMTP{
		Host: host,
		User: user,
		Pass: pass,
	}
}

func (s SMTP) SendMail(from, to string, msg []byte) error {
	host := strings.Split(s.Host, ":")[0]
	auth := smtp.PlainAuth(s.User, s.User, s.Pass, host)

	return smtp.SendMail(s.Host, auth, from, []string{to}, msg)
}

var AlertTemplate = `From: {{.From}}
To: {{.To}}
Subject: WARNING! Your IP address is blocked!

At {{.Time}} some of your IP addresses have been blocked:
{{range .Blocked}}	{{.}}
{{end}}
`

type alertTemplateValues struct {
	From    string
	To      string
	Time    string
	Blocked []net.IP
}

func (s SMTP) SendAlertTemplate(to string, blocked []net.IP) error {
	var buf bytes.Buffer

	values := alertTemplateValues{
		From:    s.User,
		To:      to,
		Time:    time.Now().Format("15:04:05 _2 Jan 2006"),
		Blocked: blocked,
	}

	t := template.Must(template.New("email").Parse(AlertTemplate))
	t.Execute(&buf, values)

	return s.SendMail(s.User, to, buf.Bytes())
}
