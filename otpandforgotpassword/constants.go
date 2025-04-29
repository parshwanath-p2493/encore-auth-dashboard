package otpandforgotpassword

import "time"

const (
	OtpExp     = time.Minute * 10
	otpCharSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

	OTPlength    = 5
	OtpKeyPrefix = "password-reset"
)
