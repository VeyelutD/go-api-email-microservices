package auth

import "errors"

var (
	ErrWrongOTP                      = errors.New("wrong OTP provided")
	ErrUserOTPNotFound               = errors.New("user OTP not found in DB")
	ErrUserConfirmationTokenNotFound = errors.New("user confirmation token not found in DB")
)
