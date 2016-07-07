package dmail

import (
	"bytes"
	"encoding/base64"
	"fmt"
)

const (
	MARKER = "DMAILEMAILBOUNDARY"

	lineMaxLength = 500 // length for attachment
)

type Email struct {
	From        string
	To          string
	Cc, Bcc     string
	Subject     string
	Body        []byte
	Attachments []Attachment
}

type Attachment struct {
	ContentType string
	Name        string
	Data        []byte
}

func NewMessage(email Email) []byte {
	buf := bytes.NewBuffer(nil)
	if email.Attachments != nil && len(email.Attachments) > 0 {
		// with attachements
		buf = newMessageMultipart(buf, email)
	} else {
		// no attachements
		buf.WriteString(fmt.Sprintf(emailTpl(), email.From, email.To, email.Cc, email.Subject)) // headers
		buf.Write(email.Body)                                                                   // body
	}
	return buf.Bytes()
}

func newMessageMultipart(buf *bytes.Buffer, email Email) *bytes.Buffer {
	buf.WriteString(fmt.Sprintf("From: %s\r\n To: %s\r\n Cc: %s\r\n Subject: %s\r\n MIME-Version: 1.0\r\n Content-Type: multipart/mixed; boundary=%s\r\n--%s", email.From, email.To, email.Cc, email.Subject, MARKER, MARKER))
	buf.WriteString(fmt.Sprintf("\r\n Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s\r\n--%s", string(email.Body), MARKER))
	buf = processAttachments(buf, email)
	return buf

}

func emailTpl() string {
	return "From: %s \r\nTo: %s \r\nCc: %s \r\nSubject: %s \r\nMIME-version: 1.0 \r\nContent-Type: text/html; charset=\"UTF-8\"\r\n"
}

func processAttachments(buf *bytes.Buffer, email Email) *bytes.Buffer {
	for _, att := range email.Attachments {
		buf.WriteString(fmt.Sprintf("\r\nContent-Type: \"%s\"; name=\"%s\"\r\nContent-Transfer-Encoding:base64\r\nContent-Disposition: attachment; filename=\"%s\"\r\n\r\n", att.ContentType, att.Name, att.Name))

		encoded := base64.StdEncoding.EncodeToString(att.Data)

		//split the encoded file in lines (doesn't matter, but low enough not to hit a max limit)
		nbrLines := len(encoded) / lineMaxLength

		//append lines to buffer
		for i := 0; i < nbrLines; i++ {
			buf.WriteString(encoded[i*lineMaxLength:(i+1)*lineMaxLength] + "\n")

		} //for

		//append last line in buffer
		buf.WriteString(encoded[nbrLines*lineMaxLength:])

		//part 3 will be the attachment
		buf.WriteString(fmt.Sprintf("\r\n--%s--", MARKER))
	}
	return buf
}
