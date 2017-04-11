package drops

import "html/template"

type DialogType int

const (
	Error DialogType = iota
)

type Dialog struct {
	Type    DialogType
	Title   string
	Message string
	Html    template.HTML
	Actions []MenuItem
	Form    Form
}
type NotificationType int

const (
	Success NotificationType = iota
	Warning
	Info
)

type Notification struct {
	Type NotificationType
	Text string
}

type Form struct {
	Action string
	Class  string
	Data   map[string]string
	Submit MenuItem
}
