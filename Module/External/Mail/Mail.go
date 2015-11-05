package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/griesbacher/Yarbes/Config"
	"github.com/griesbacher/Yarbes/Config/ConfigLayouts"
	"github.com/griesbacher/Yarbes/Tools/Strings"
	"log"
	"net"
	"net/mail"
	"net/smtp"
)

//TODO: Client log bib schreiben
func main() {
	var addressField string
	var address string
	var event string
	addresses := []mail.Address{}
	flag.Usage = func() {
		fmt.Println(`Yarbes-Mail by Philip Griesbacher @ 2015`)
	}
	flag.StringVar(&addressField, "addressField", "address", "references the filed in which the emailaddress can be found")
	flag.StringVar(&address, "address", "root@example.net", "address to send mail to")
	flag.StringVar(&event, "event", "", "the event")
	flag.Parse()

	if addressField != "" {
		jsonMap := Strings.UnmarshalJSONEvent(event)
		addresses = append(addresses, mail.Address{Address: jsonMap[addressField]})
	}
	if address != "" {
		addresses = append(addresses, mail.Address{Address: addressField})
	}

	Config.InitMailConfig("Module/External/Mail/mail.gcfg")
	con := Config.GetMailConfig()
	Sendmail("mail by Yarbes", event, con, true, addresses)

}

//Sendmail sends a email to the given address, useTLS can be used for encryption
func Sendmail(subj, body string, config *ConfigLayouts.Mail, useTLS bool, to ...mail.Address) {
	from := mail.Address{Name: config.Mail.FromName, Address: config.Mail.FromAddress}
	headers := make(map[string]string)
	headers["From"] = from.String()
	for _, address := range to {
		headers["To"] += address.String()
	}
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
		for _, address := range to {
			if err = c.Rcpt(address.Address); err != nil {
				log.Panic(err)
			}
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
		sendTo := []string{}
		for _, address := range to {
			sendTo = append(sendTo, address.Address)
		}
		err := smtp.SendMail(
			config.Mail.Server,
			auth,
			from.Address,
			sendTo,
			[]byte(message),
		)
		if err != nil {
			log.Fatal(err)
		}
	}
}
