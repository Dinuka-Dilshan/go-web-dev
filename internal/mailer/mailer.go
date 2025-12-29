package mailer

import "embed"

const (
	FromName            = "GoperSocial"
	MaxRetries          = 3
	UserWelcomeTemplate = "user-invitation.tmpl"
)

//go:embed templates
var FS embed.FS

type Client interface {
	Send(templateFile, username, email string, data any, isSandBox bool) error
}
