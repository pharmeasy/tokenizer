package cryptography

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/pharmaeasyteam/goframework/logging"
	"bitbucket.org/pharmaeasyteam/goframework/render"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/errormanager"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/identity"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/keysetmanager"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/db"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/decryption"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/encryption"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/metadata"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/tokenmanager"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/validator"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/getsentry/sentry-go"
	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/keyset"
	"go.uber.org/zap"
)

func (c *ModuleCrypto) encrypt(w http.ResponseWriter, req *http.Request) {
	sentry.CaptureMessage("Sentry integrated in tokenizer")
	requestParams, err := validator.ValidateEncryptionRequest(req)
	if err != nil {
		err = errormanager.SetEncryptionError(requestParams, err, http.StatusBadRequest)

		errormanager.RenderEncryptionErrorResponse(w, req, http.StatusBadRequest, err)
		return
	}

	// validate identifier
	isAuthenticated := identity.AuthenticateRequest(requestParams.Identifier)
	if !isAuthenticated {
		err = errormanager.SetEncryptionError(requestParams, nil, http.StatusForbidden)

		errormanager.RenderEncryptionErrorResponse(w, req, http.StatusForbidden, err)
		return
	}

	// authorize encryption levels
	isAuthorized := identity.AuthorizeLevelForEncryption(requestParams)
	if !isAuthorized {
		err = errormanager.SetEncryptionError(requestParams, nil, http.StatusForbidden)

		errormanager.RenderEncryptionErrorResponse(w, req, http.StatusForbidden, err)
		return
	}
	// encrypt data
	encryptedData, err := encryptTokenData(requestParams, c, req.Context())
	if err != nil {
		err = errormanager.SetEncryptionError(requestParams, err, http.StatusInternalServerError)

		errormanager.RenderEncryptionErrorResponse(w, req, http.StatusInternalServerError, err)
		return
	}
	render.JSON(w, req, encryptedData)
}

func (c *ModuleCrypto) decrypt(w http.ResponseWriter, req *http.Request) {

	requestParams, err := validator.ValidateDecryptionRequest(req)
	if err != nil {
		err = errormanager.SetDecryptionError(requestParams, err, http.StatusBadRequest)

		errormanager.RenderDecryptionErrorResponse(w, req, http.StatusBadRequest, err)
		return
	}

	// validate identifier
	isAuthenticated := identity.AuthenticateRequest(requestParams.Identifier)
	if !isAuthenticated {
		err = errormanager.SetDecryptionError(requestParams, nil, http.StatusForbidden)

		errormanager.RenderDecryptionErrorResponse(w, req, http.StatusForbidden, err)
		return
	}

	// fetch records
	tokenData, err := getTokenData(requestParams, c)
	if err != nil {
		err = errormanager.SetDecryptionError(requestParams, err, http.StatusInternalServerError)

		errormanager.RenderDecryptionErrorResponse(w, req, http.StatusInternalServerError, err)
		return
	}

	// authorize token access
	isAuthorized := identity.AuthorizeTokenAccess(tokenData, requestParams.Identifier)
	if !isAuthorized {
		err = errormanager.SetDecryptionError(requestParams, nil, http.StatusForbidden)

		errormanager.RenderDecryptionErrorResponse(w, req, http.StatusForbidden, err)

		return
	}
	decryptedData, err := decryptTokenData(tokenData, requestParams)
	if err != nil {
		err = errormanager.SetDecryptionError(requestParams, err, http.StatusInternalServerError)

		errormanager.RenderDecryptionErrorResponse(w, req, http.StatusInternalServerError, err)
		return
	}

	render.JSON(w, req, decryptedData)

}

func (c *ModuleCrypto) getMetaData(w http.ResponseWriter, req *http.Request) {

	// validate request params
	requestParams, err := validator.ValidateMetadataRequest(req)
	if err != nil {
		errormanager.RenderGetMetadataErrorResponse(w, req, http.StatusBadRequest,
			errormanager.SetMetadataError(requestParams, err, http.StatusBadRequest))
		return
	}

	// validate access
	isAuthenticated := identity.AuthenticateRequest(requestParams.Identifier)
	if !isAuthenticated {
		errormanager.RenderGetMetadataErrorResponse(w, req, http.StatusForbidden,
			errormanager.SetMetadataError(requestParams, nil, http.StatusForbidden))
		return
	}

	// fetch records
	dbInterface := c.database
	tokenData, err := dbInterface.GetItemsByTokenInBatch(requestParams.Tokens)
	if err != nil {
		errormanager.RenderGetMetadataErrorResponse(w, req, http.StatusInternalServerError,
			errormanager.SetMetadataError(requestParams, err, http.StatusInternalServerError))
		return
	}

	// authorize token access
	isAuthorized := identity.AuthorizeTokenAccess(&tokenData, requestParams.Identifier)
	if !isAuthorized {
		errormanager.RenderGetMetadataErrorResponse(w, req, http.StatusForbidden,
			errormanager.SetMetadataError(requestParams, nil, http.StatusForbidden))
		return
	}

	// return metadata

	metadataResponse := getMetaItems(tokenData)

	render.JSON(w, req, metadataResponse)
}

func (c *ModuleCrypto) updateMetadata(w http.ResponseWriter, req *http.Request) {

	// validate request params
	requestParams, err := validator.ValidateMetadataUpdateRequest(req)
	if err != nil {
		errormanager.RenderUpdateMetadataErrorResponse(w, req, http.StatusBadRequest,
			errormanager.SetUpdateMetadataError(requestParams, err, http.StatusBadRequest))
		return
	}

	// validate access
	isAuthenticated := identity.AuthenticateRequest(requestParams.Identifier)
	if !isAuthenticated {
		errormanager.RenderUpdateMetadataErrorResponse(w, req, http.StatusForbidden,
			errormanager.SetUpdateMetadataError(requestParams, nil, http.StatusForbidden))
		return
	}

	// fetch records
	payloadSize := len(requestParams.UpdateParams)
	tokenIDs := make([]string, payloadSize)

	for i := 0; i < payloadSize; i++ {
		tokenIDs[i] = requestParams.UpdateParams[i].Token
	}

	dbInterface := c.database

	tokenData, err := dbInterface.GetItemsByTokenInBatch(tokenIDs)
	if err != nil {
		errormanager.RenderUpdateMetadataErrorResponse(w, req, http.StatusInternalServerError,
			errormanager.SetUpdateMetadataError(requestParams, err, http.StatusInternalServerError))
		return
	}

	// authorize token access
	isAuthorized := identity.AuthorizeTokenAccess(&tokenData, requestParams.Identifier)
	if !isAuthorized {
		errormanager.RenderUpdateMetadataErrorResponse(w, req, http.StatusForbidden,
			errormanager.SetUpdateMetadataError(requestParams, nil, http.StatusForbidden))
		return
	}

	//update metadata
	err = updateMetaItems(requestParams, tokenData, c)
	if err != nil {
		errormanager.RenderUpdateMetadataErrorResponse(w, req, http.StatusInternalServerError,
			errormanager.SetUpdateMetadataError(requestParams, err, http.StatusInternalServerError))
		return
	}

	render.JSON(w, req, "Metadata updated successfully.")
}

func getTokenData(requestParams *decryption.DecryptRequest, c *ModuleCrypto) (*map[string]db.TokenData, error) {

	payloadSize := len(requestParams.DecryptRequestData)
	tokenIDs := make([]string, payloadSize)

	for i := 0; i < payloadSize; i++ {
		tokenIDs[i] = requestParams.DecryptRequestData[i].Token
	}

	dbInterface := c.database
	tokenData, err := dbInterface.GetItemsByTokenInBatch(tokenIDs)

	if err != nil {
		return nil, err
	}

	return &tokenData, nil
}

func encryptTokenData(requestParams *encryption.EncryptRequest, c *ModuleCrypto, ctx context.Context) (*encryption.EncryptResponse, error) {
	encryptionResponse := encryption.EncryptResponse{}
	reqParamsData := requestParams.EncryptRequestData
	dbInterface := c.database
	// get keyset handler

	keyName, keysetHandle, err := keysetmanager.GetKeysetHandlerForEncryption()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(reqParamsData); i++ {
		// encrypt text
		ciphertext, err := dataEncryptAEAD(reqParamsData[i].Content, reqParamsData[i].Salt, keysetHandle)
		if err != nil {
			return nil, err
		}

		integrityChecker := integrityChecker(ciphertext, reqParamsData[i].Content, reqParamsData[i].Salt, keysetHandle)
		if !integrityChecker {
			return nil, errormanager.SetError("data integrity compromised", nil)
		}

		dbTokenData := db.TokenData{
			Level:     reqParamsData[i].Level,
			Content:   ciphertext,
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
			Key:       *keyName,
			Metadata:  reqParamsData[i].Metadata,
		}

		token, err := storeEncryptedData(dbTokenData, c)
		if err != nil {
			return nil, err
		}

		integrityCheckerAdvancedFlag := false
		maxAttempts := 5
		currentAttempts := 0

		for !integrityCheckerAdvancedFlag && currentAttempts < maxAttempts {
			integrityCheckerAdvancedFlag = integrityCheckerAdvanced(token, reqParamsData[i].Content, reqParamsData[i].Salt, keysetHandle, c, requestParams.RequestID)
			currentAttempts = currentAttempts + 1

			if !integrityCheckerAdvancedFlag {
				logging.GetLogger().Info(fmt.Sprintf("database integrity failed for attempt no. %d for request_id : %s", currentAttempts, requestParams.RequestID))
			}

		}
		if !integrityCheckerAdvancedFlag {
			err := dbInterface.DeleteItemByToken(token)
			if err != nil {
				return nil, err
			}
			return nil, errormanager.SetError("database integrity compromised after max attempts for request_id : %s", nil)
		}

		encryptionResponse.ResponseData = append(encryptionResponse.ResponseData,
			encryption.ResponseData{
				ID:    reqParamsData[i].ID,
				Token: tokenmanager.FormatToken(token),
			})

	}

	return &encryptionResponse, nil
}

func dataEncryptAEAD(data string, salt string, kh *keyset.Handle) ([]byte, error) {
	a, err := aead.New(kh)
	if err != nil {
		return nil, errormanager.SetError("Error encountered during AEAD wrapper generation.", err)
	}

	ct, err := a.Encrypt([]byte(data), []byte(salt))
	if err != nil {
		return nil, errormanager.SetError("Error encountered during AEAD generation", err)
	}

	return ct, nil
}
func dataDecryptAEAD(cipherText []byte, salt string, kh *keyset.Handle) (*string, error) {
	a, err := aead.New(kh)
	if err != nil {
		return nil, errormanager.SetError("Error encountered while initializing AEAD handler.", err)
	}

	pt, err := a.Decrypt(cipherText, []byte(salt))
	if err != nil {
		return nil, errormanager.SetError("Error encountered while decrypting data with AEAD.", err)
	}

	plainText := string(pt)

	return &plainText, nil
}

func storeEncryptedData(dbTokenData db.TokenData, c *ModuleCrypto) (string, error) {
	attempts := 0
	var err error
	dbInterface := c.database
	for attempts < 3 {
		attempts++
		dbTokenData.TokenID = tokenmanager.Uniquetoken()
		err = dbInterface.PutItem(dbTokenData)
		// handle token clashes for 3 attempts
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeConditionalCheckFailedException:
					continue
				}
			}
		}
		break
	}

	return dbTokenData.TokenID, err
}

func decryptTokenData(tokenData *map[string]db.TokenData, requestParams *decryption.DecryptRequest) (*decryption.DecryptResponse, error) {
	decryptionResponse := decryption.DecryptResponse{}
	reqParamsData := requestParams.DecryptRequestData
	for i := 0; i < len(reqParamsData); i++ {
		token := reqParamsData[i].Token
		dbTokenData := (*tokenData)[token]

		// select keyset
		kh, err := keysetmanager.GetKeysetHandlerForDecryption(dbTokenData.Key)
		if err != nil {
			return nil, err
		}

		// decrypt with salt
		decryptedText, err := dataDecryptAEAD(dbTokenData.Content, reqParamsData[i].Salt, kh)
		if err != nil {
			return nil, err
		}

		decryptionResponse.DecryptionResponseData = append(decryptionResponse.DecryptionResponseData,
			decryption.DecryptResponseData{
				Token:    tokenmanager.FormatToken(token),
				Content:  *decryptedText,
				Metadata: dbTokenData.Metadata,
			})
	}

	return &decryptionResponse, nil
}

func updateMetaItems(requestParams *metadata.MetaUpdateRequest, tokenData map[string]db.TokenData, c *ModuleCrypto) error {

	for _, v := range requestParams.UpdateParams {
		meta := tokenData[v.Token]
		meta.Metadata = v.Metadata
		meta.UpdatedAt = time.Now().Format(time.RFC3339)
		tokenData[v.Token] = meta
	}
	dbInterface := c.database
	for k, v := range tokenData {
		err := dbInterface.UpdateMetadataByToken(k, v.Metadata, v.UpdatedAt)
		if err != nil {
			return err
		}
	}

	return nil
}

func getMetaItems(tokenData map[string]db.TokenData) *metadata.MetaResponse {
	metaResponse := metadata.MetaResponse{}

	for _, dbTokenData := range tokenData {
		metaResponse.MetaParams = append(metaResponse.MetaParams,
			metadata.MetaParams{
				Token:    tokenmanager.FormatToken(dbTokenData.TokenID),
				Metadata: dbTokenData.Metadata,
			})
	}

	return &metaResponse
}

// integrity checker
func integrityChecker(cipherText []byte, plainText string, salt string, kh *keyset.Handle) bool {
	pt, _ := dataDecryptAEAD(cipherText, salt, kh)
	if *pt == plainText {
		return true
	}
	return false
}

func integrityCheckerAdvanced(token string, plainText string, salt string, kh *keyset.Handle, c *ModuleCrypto, requestID string) bool {
	dbInterface := c.database
	tokenData, err := dbInterface.GetItemsByToken([]string{token})
	if err != nil {
		logging.GetLogger().Error(fmt.Sprintf("Failed to read from database for Request id %s: ", requestID), zap.Error(err))
		return false
	}

	cipherText := tokenData[token].Content
	plainTextFromDB, err := dataDecryptAEAD(cipherText, salt, kh)
	if err != nil {
		logging.GetLogger().Error(fmt.Sprintf("Failed to decrypt data for Request id %s: ", requestID), zap.Error(err))
		return false
	}

	if *plainTextFromDB == plainText {
		return true
	}
	logging.GetLogger().Error(fmt.Sprintf("data not matched for Request id %s: ", requestID))

	return false
}
