package email

import "errors"

var ErrCouldNotSendOTP = errors.New("could not send OTP")
var ErrCouldNotSendConfirmation = errors.New("could not send confirmation link")
