package errormanager

import (
	"errors"
	"net/http"

	"bitbucket.org/pharmaeasyteam/goframework/logging"
	"go.uber.org/zap"

	"bitbucket.org/pharmaeasyteam/goframework/render"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/decryption"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/encryption"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/errorresponse"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/metadata"
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

// SetValidationEmptyError sets an empty error
func SetValidationEmptyError(value string) error {
	return errors.New(value + " is blank or not in range")
}

// SetValidationDecodeError sets errors in decoding
func SetValidationDecodeError(requestType string, err error) error {
	return errors.New("Unable to decode " + requestType + " request params." + err.Error())
}

// SetError Sets error based on error context
func SetError(errorContext string, err error) error {
	if err == nil {
		return errors.New(errorContext)
	}
	return errors.New(errorContext + err.Error())
}

// RenderEncryptionErrorResponse renders encryption error response
func RenderEncryptionErrorResponse(w http.ResponseWriter, req *http.Request, status uint, err error) {
	render.JSONWithStatus(w, req, int(status), errorresponse.ExceptionResponse(status, err.Error()))
}

// RenderDecryptionErrorResponse renders decryption error response
func RenderDecryptionErrorResponse(w http.ResponseWriter, req *http.Request, status uint, err error) {
	render.JSONWithStatus(w, req, int(status), errorresponse.ExceptionResponse(status, err.Error()))
}

// RenderGetMetadataErrorResponse renders get metadata by tokens error response
func RenderGetMetadataErrorResponse(w http.ResponseWriter, req *http.Request, status uint, err error) {
	render.JSONWithStatus(w, req, int(status), errorresponse.ExceptionResponse(status, err.Error()))
}

// RenderUpdateMetadataErrorResponse renders updatemetadata error response
func RenderUpdateMetadataErrorResponse(w http.ResponseWriter, req *http.Request, status uint, err error) {
	render.JSONWithStatus(w, req, int(status), errorresponse.ExceptionResponse(status, err.Error()))
}

//SetDecryptionError logs non sensitive encryption data and returns a generic error
func SetDecryptionError(requestParams *decryption.DecryptRequest, err error, status uint) error {
	genericError := getGenericErrorByStatus(status)
	if requestParams != nil {
		for i := 0; i < len(requestParams.DecryptRequestData); i++ {
			requestParams.DecryptRequestData[i].Salt = ""
			requestParams.DecryptRequestData[i].Token = ""
		}
		logging.GetLogger().Error(genericError.Error(), zap.Error(err), zap.Any("decryptionRequest", &requestParams))
		return genericError
	}
	logging.GetLogger().Error(genericError.Error(), zap.Error(err))
	return genericError
}

// SetMetadataError logs non sensitive metadata related information and returns a generic error
func SetMetadataError(requestParams *metadata.MetaRequest, err error, status uint) error {
	genericError := getGenericErrorByStatus(status)
	if requestParams != nil {
		for i := 0; i < len(requestParams.Tokens); i++ {
			requestParams.Tokens[i] = ""
		}
		logging.GetLogger().Error(genericError.Error(), zap.Error(err), zap.Any("getmetadataRequest", &requestParams))
		return genericError
	}
	logging.GetLogger().Error(genericError.Error(), zap.Error(err))
	return genericError
}

// SetUpdateMetadataError logs non sensitive updatemetadata related information and returns a generic error
func SetUpdateMetadataError(requestParams *metadata.MetaUpdateRequest, err error, status uint) error {
	genericError := getGenericErrorByStatus(status)
	if requestParams != nil {
		for i := 0; i < len(requestParams.UpdateParams); i++ {
			requestParams.UpdateParams[i].Metadata = ""
			requestParams.UpdateParams[i].Token = ""
		}

		logging.GetLogger().Error(genericError.Error(), zap.Error(err), zap.Any("updatemetaRequest", &requestParams))
		return genericError

	}
	logging.GetLogger().Error(genericError.Error(), zap.Error(err))
	return genericError
}
