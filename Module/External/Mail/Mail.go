package main

import (
	"crypto/tls"
	"fmt"
	"github.com/griesbacher/SystemX/Config"
	"github.com/griesbacher/SystemX/Config/ConfigLayouts"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		Config.InitMailConfig("Module/External/Mail/mail.gcfg")
		con := Config.GetMailConfig()
		Sendmail(mail.Address{"philip", "griesbacher@consol.de"}, "mail by SystemX", os.Args[1], con, true)
	}
}

func Sendmail(to mail.Address, subj, body string, config *ConfigLayouts.Mail, useTLS bool) {
	from := mail.Address{config.Mail.FromName, config.Mail.FromAddress}
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	host, _, _ := net.SplitHostPort(config.Mail.Server)
	auth := smtp.PlainAuth("", config.Mail.Username, config.Mail.Password, host)

	if useTLS {
		conn, err := tls.Dial("tcp", config.Mail.Server, &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         host,
		})
		if err != nil {
			log.Panic(err)
		}

		c, err := smtp.NewClient(conn, host)
		if err != nil {
			log.Panic(err)
		}
		defer c.Quit()

		if err = c.Auth(auth); err != nil {
			log.Panic(err)
		}

		if err = c.Mail(from.Address); err != nil {
			log.Panic(err)
		}

		if err = c.Rcpt(to.Address); err != nil {
			log.Panic(err)
		}

		w, err := c.Data()
		if err != nil {
			log.Panic(err)
		}
		defer w.Close()

		_, err = w.Write([]byte(message))
		if err != nil {
			log.Panic(err)
		}
	} else {
		err := smtp.SendMail(
			config.Mail.Server,
			auth,
			from.Address,
			[]string{to.Address},
			[]byte(message),
		)
		if err != nil {
			log.Fatal(err)
		}
	}
}
