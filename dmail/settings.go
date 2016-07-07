package dmail

type MailSettings struct {
	Server   string
	Port     int
	Username string
	Password string
}

func NewMailSettings(server string, port int, username, password string) MailSettings {
	return MailSettings{
		Server:   server,
		Port:     port,
		Username: username,
		Password: password,
	}
}
