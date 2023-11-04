package email

type MessageType string

const (
	MessageHTML      = "text/html"
	MessagePlainText = "text/plain"
)

type message struct {
	text        string
	from        string
	to          string
	contentType MessageType
}
