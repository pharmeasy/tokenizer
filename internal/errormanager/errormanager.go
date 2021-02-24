package errormanager

import (
	"errors"
	"net/http"

	"bitbucket.org/pharmaeasyteam/goframework/logging"
	"go.uber.org/zap"

	"bitbucket.org/pharmaeasyteam/goframework/render"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/badresponse"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/encryption"
)

var genericBadRequestError error
var genericInternalServerError error
var genericForbiddenError error

// SetGenericErrors sets generic errors
func SetGenericErrors() {
	genericBadRequestError = errors.New("Invalid request parameters passed")
	genericForbiddenError = errors.New("You are not allowed to perform this action")
	genericInternalServerError = errors.New("Something went wrong. Please try again later")
}

func getGenericErrorByStatus(status uint) error {
	switch status {
	case http.StatusForbidden:
		return genericForbiddenError
	case http.StatusBadRequest:
		return genericBadRequestError
	case http.StatusInternalServerError:
		return genericInternalServerError
	}

	return nil
}

// SetEncryptionError logs non sensitive encryption data and returns a generic error
func SetEncryptionError(requestParams *encryption.EncryptRequest, err error, status uint) error {
	genericError := getGenericErrorByStatus(status)
	if requestParams != nil {
		for i := 0; i < len(requestParams.EncryptRequestData); i++ {
			requestParams.EncryptRequestData[i].Content = ""
			requestParams.EncryptRequestData[i].Salt = ""
		}
		logging.GetLogger().Error(genericError.Error(), zap.Error(err), zap.Any("encryptionRequest", &requestParams))
		return genericError
	}
	logging.GetLogger().Error(genericError.Error(), zap.Error(err))

	return genericError
}

// RenderEncryptionErrorResponse renders encryption error response
func RenderEncryptionErrorResponse(w http.ResponseWriter, req *http.Request, status uint, err error) {
	render.JSONWithStatus(w, req, int(status), badresponse.ExceptionResponse(status, err.Error()))
}

// SetValidationEmptyError sets an empty error
func SetValidationEmptyError(value string) error {
	return errors.New(value + " is blank")
}

// SetValidationDecodeError sets errors in decoding
func SetValidationDecodeError(requestType string, err error) error {
	return errors.New("Unable to decode " + requestType + " request params." + err.Error())
}

// SetError Sets error based on error context
func SetError(errorContext string, err error) error {
	return errors.New(errorContext + err.Error())
}
