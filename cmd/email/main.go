package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/tsivinsky/email"
)

var (
	user     = flag.String("u", "", "auth user")
	password = flag.String("p", "", "auth password")
	hostname = flag.String("H", "", "smtp hostname")
	port     = flag.String("P", "465", "smtp port")
	from     = flag.String("f", "", "email address to send from")
	isHTML   = flag.Bool("html", false, "send email in html format")
)

func main() {
	flag.Parse()

	host := net.JoinHostPort(*hostname, *port)

	auth, err := email.NewAuth(*user, *password, host)
	if err != nil {
		panic(err)
	}

	t := email.NewTransport(auth, *from)

	to := mustPrompt("Where to send email: ")
	subject := mustPrompt("Subject: ")

	body, err := getEmailBody()
	if err != nil {
		panic(err)
	}

	var contentType email.MessageType = email.MessagePlainText
	if *isHTML {
		contentType = email.MessageHTML
	}

	err = t.SendEmail(to, subject, body, contentType)
	if err != nil {
		panic(err)
	}

	fmt.Println("email should be on the way")
}
