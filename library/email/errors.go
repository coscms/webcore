package email

import "errors"

var (
	ErrSMTPNoSet          = errors.New(`SMTP is not set`)
	ErrSenderNoSet        = errors.New(`The sender address is not set`)
	ErrRecipientNoSet     = errors.New(`The recipient address is not set`)
	ErrSendChannelTimeout = errors.New(`SendMail: The sending channel timed out`)
)
