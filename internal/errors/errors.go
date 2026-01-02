package errors

import "errors"

var (
	NothingToUpdateErr           = errors.New("nothing to update")
	BodyInvalidFormatErr         = errors.New("request body has invalid format")
	BodyInvalidContentErr        = errors.New("body content is not valid")
	DuplicateValueErr            = errors.New("value already exists")
	RegistrationIDNotFoundErr    = errors.New("registration id not found")
	RegistrationEmailMismatchErr = errors.New("registration email mismatch")
	UserNotFoundErr              = errors.New("user not found")
)
