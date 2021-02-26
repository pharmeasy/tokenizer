package validator

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bitbucket.org/pharmaeasyteam/tokenizer/internal/errormanager"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/decryption"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/encryption"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/metadata"
)

// ValidateEncryptionRequest provides validation logic for the incoming encryption request
func ValidateEncryptionRequest(req *http.Request) (*encryption.EncryptRequest, error) {
	set := make(map[string]bool)
	decoder := json.NewDecoder(req.Body)
	params := encryption.EncryptRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		return nil, errormanager.SetValidationDecodeError("encryption", err)
	}

	if params.RequestID == "" {
		return &params, errormanager.SetValidationEmptyError("Request ID")
	}

	if params.Identifier == "" {
		return &params, errormanager.SetValidationEmptyError("Identifier")
	}

	for _, v := range params.EncryptRequestData {
		level, _ := strconv.Atoi(v.Level)
		if v.Level == "" || level < 1 || level > 7 {
			return &params, errormanager.SetValidationEmptyError("Level")
		}

		if v.Content == "" {
			return &params, errormanager.SetValidationEmptyError("Content")
		}

		if v.ID == "" {
			return &params, errormanager.SetValidationEmptyError("ID")
		}

		if !set[v.ID] {
			set[v.ID] = true
		} else {
			return &params, errormanager.SetValidationEmptyError("ID may be repeated or")
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
		return nil, errormanager.SetValidationDecodeError("decryption", err)
	}

	if params.RequestID == "" {
		return &params, errormanager.SetValidationEmptyError("Request ID")
	}

	if params.Identifier == "" {
		return &params, errormanager.SetValidationEmptyError("Identifier")
	}

	for i := 0; i < len(params.DecryptRequestData); i++ {
		if params.DecryptRequestData[i].Token == "" {
			return &params, errormanager.SetValidationEmptyError("Token")
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
		return nil, errormanager.SetValidationDecodeError("metadata", err)
	}

	if params.RequestID == "" {
		return &params, errormanager.SetValidationEmptyError("Request ID")
	}

	if params.Identifier == "" {
		return &params, errormanager.SetValidationEmptyError("Identifier")
	}

	for i := 0; i < len(params.Tokens); i++ {
		if params.Tokens[i] == "" {
			return &params, errormanager.SetValidationEmptyError("Token")
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
		return nil, errormanager.SetValidationDecodeError("metadata update", err)
	}

	if params.RequestID == "" {
		return &params, errormanager.SetValidationEmptyError("Request ID")
	}

	if params.Identifier == "" {
		return &params, errormanager.SetValidationEmptyError("Identifier")
	}

	for i := 0; i < len(params.UpdateParams); i++ {
		if params.UpdateParams[i].Token == "" {
			return &params, errormanager.SetValidationEmptyError("Token")
		}

		if params.UpdateParams[i].Metadata == "" {
			return &params, errormanager.SetValidationEmptyError("Metadata")
		}
	}

	return &params, nil
}
