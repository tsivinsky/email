package email

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
)

type transport struct {
	auth   *auth
	sender string
}

func (t *transport) createMessage(recipient, subject string, body []byte, contentType MessageType) *message {
	from := mail.Address{Name: "", Address: t.sender}
	to := mail.Address{Name: "", Address: recipient}

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject
	headers["Content-Type"] = string(contentType)

	text := ""
	for k, v := range headers {
		text += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	text += "\r\n" + string(body)

	return &message{
		text:        text,
		from:        t.sender,
		to:          recipient,
		contentType: contentType,
	}
}

func (t *transport) createClient(conn net.Conn) (*smtp.Client, error) {
	client, err := smtp.NewClient(conn, t.auth.hostname)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (t *transport) sendMessage(message *message) error {
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         t.auth.getHost(),
	}

	conn, err := tls.Dial("tcp", t.auth.getHost(), tlsconfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := t.createClient(conn)
	defer client.Quit()

	if err = client.Auth(t.auth.auth); err != nil {
		return err
	}

	if err = client.Mail(message.from); err != nil {
		return err
	}

	if err = client.Rcpt(message.to); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message.text))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}

func (t *transport) SendEmail(recipient, subject string, body []byte, contentType MessageType) error {
	msg := t.createMessage(recipient, subject, body, contentType)
	return t.sendMessage(msg)
}

func NewTransport(auth *auth, sender string) *transport {
	return &transport{
		auth:   auth,
		sender: sender,
	}
}
