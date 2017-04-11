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

func (t NotificationType) String() string {
	switch t {
	case Success:
		return "success"
	case Warning:
		return "warning"
	case Info:
		return "information"
	default:
		return "success"
	}
}

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
