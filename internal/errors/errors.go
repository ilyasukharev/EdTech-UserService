package errors

import "errors"

var (
	NothingToUpdateErr             = errors.New("nothing to update")
	BodyInvalidFormatErr           = errors.New("request body has invalid format")
	BodyInvalidContentErr          = errors.New("body content is not valid")
	DuplicateValueErr              = errors.New("value already exists")
	RegistrationIDNotFoundErr      = errors.New("registration id not found")
	RegistrationEmailMismatchErr   = errors.New("registration email mismatch")
	UserNotFoundErr                = errors.New("user not found")
	ChildNotFoundErr               = errors.New("child not found")
	ReferralNotFoundErr            = errors.New("referral not found")
	OTPCodeMismatchErr             = errors.New("OTP code mismatch")
	PathArgumentValueIncorrectErr  = errors.New("path argument value is incorrect")
	QueryArgumentValueIncorrectErr = errors.New("query argument value is incorrect")
)
