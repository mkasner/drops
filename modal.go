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
}
type NotificationType int

const (
	Success NotificationType = iota
)

type Notification struct {
	Type NotificationType
	Text string
}
