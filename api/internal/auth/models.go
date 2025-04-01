package auth

type VerifyOTPBody struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=4"`
}
type SendOTPBody struct {
	Email string `json:"email" validate:"required,email"`
}
