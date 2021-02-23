package validator

import (
	"encoding/json"
	"errors"
	"net/http"

	"bitbucket.org/pharmaeasyteam/goframework/logging"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/decryption"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/encryption"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/metadata"
	"go.uber.org/zap"
)

var errorGeneric error

// ValidateEncryptionRequest provides validation logic for the incoming encryption request
func ValidateEncryptionRequest(req *http.Request) (*encryption.EncryptRequest, error) {
	decoder := json.NewDecoder(req.Body)
	params := encryption.EncryptRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		logDecodeError("encryption", err)
		return nil, errorGeneric
	}

	if params.RequestID == "" {
		logEmptyError("Request ID", params.RequestID)
		return nil, errorGeneric
	}

	if params.Identifier == "" {
		logEmptyError("Identifier", params.RequestID)
		return nil, errorGeneric
	}

	if params.Level == "" {
		logEmptyError("Level", params.RequestID)
		return nil, errorGeneric
	}

	for _, v := range params.RequestData {
		if v.Content == "" {
			logEmptyError("Content", params.RequestID)
			return nil, errorGeneric
		}
	}

	return &params, nil
}

// ValidateDecryptionRequest provides validation logic for the incoming decryption request
func ValidateDecryptionRequest(req *http.Request) (*decryption.DecryptRequest, error) {
	decoder := json.NewDecoder(req.Body)
	params := decryption.DecryptRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		logDecodeError("decryption", err)
		return nil, errorGeneric
	}

	if params.RequestID == "" {
		logEmptyError("Request ID", params.RequestID)
		return nil, errorGeneric
	}

	if params.Level == "" {
		logEmptyError("Level", params.RequestID)
		return nil, errorGeneric
	}

	if params.Identifier == "" {
		logEmptyError("Identifier", params.RequestID)
		return nil, errorGeneric
	}

	for i := 0; i < len(params.DecryptRequestData); i++ {
		if params.DecryptRequestData[i].Token == "" {
			logEmptyError("Token", params.RequestID)
			return nil, errorGeneric
		}
	}

	return &params, nil
}

// ValidateMetadataRequest provides validation logic for the incoming metadata request
func ValidateMetadataRequest(req *http.Request) (*metadata.MetaRequest, error) {
	decoder := json.NewDecoder(req.Body)
	params := metadata.MetaRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		logDecodeError("metadata", err)
		return nil, errorGeneric
	}

	if params.RequestID == "" {
		logEmptyError("Request ID", params.RequestID)
		return nil, errorGeneric
	}

	if params.Level == "" {
		logEmptyError("Level", params.RequestID)
		return nil, errorGeneric
	}

	if params.Identifier == "" {
		logEmptyError("Identifier", params.RequestID)
		return nil, errorGeneric
	}

	for i := 0; i < len(params.Tokens); i++ {
		if params.Tokens[i] == "" {
			logEmptyError("Token", params.RequestID)
			return nil, errorGeneric
		}
	}

	return &params, nil
}

// ValidateMetadataUpdateRequest provides validation logic for the incoming metadata update request
func ValidateMetadataUpdateRequest(req *http.Request) (*metadata.MetaUpdateRequest, error) {
	decoder := json.NewDecoder(req.Body)
	params := metadata.MetaUpdateRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		logDecodeError("metadata update", err)
		return nil, errorGeneric
	}

	if params.RequestID == "" {
		logEmptyError("Request ID", params.RequestID)
		return nil, errorGeneric
	}

	if params.Level == "" {
		logEmptyError("Level", params.RequestID)
		return nil, errorGeneric
	}

	if params.Identifier == "" {
		logEmptyError("Identifier", params.RequestID)
		return nil, errorGeneric
	}

	for i := 0; i < len(params.UpdateParams); i++ {
		if params.UpdateParams[i].Token == "" {
			logEmptyError("Token", params.RequestID)
			return nil, errorGeneric
		}

		if params.UpdateParams[i].Metadata == "" {
			logEmptyError("Metadata", params.RequestID)
			return nil, errorGeneric
		}
	}

	return &params, nil
}

// logEmptyError logs an empty error
func logEmptyError(value string, requestID string) {
	logging.GetLogger().Error(value+" is blank", zap.Any("requestId", requestID))
}

// logDecodeError logs errors in decoding
func logDecodeError(value string, err error) {
	logging.GetLogger().Error("Unable to decode "+value+" request params", zap.Error(err))
}

// SetGenericError sets generic error upon app load
func SetGenericError() {
	errorGeneric = errors.New("invalid request parameters passed")
}
