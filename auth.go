package email

import (
	"net"
	"net/smtp"
)

type auth struct {
	auth     smtp.Auth
	hostname string
	port     string
}

func (a *auth) getHost() string {
	return net.JoinHostPort(a.hostname, a.port)
}

func NewAuth(user, password, host string) (*auth, error) {
	hostname, port, err := net.SplitHostPort(host)
	if err != nil {
		return nil, err
	}

	return &auth{
		auth:     smtp.PlainAuth("", user, password, hostname),
		hostname: hostname,
		port:     port,
	}, nil
}
