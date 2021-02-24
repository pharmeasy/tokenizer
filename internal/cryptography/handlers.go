package cryptography

import (
	"net/http"
	"time"

	"bitbucket.org/pharmaeasyteam/goframework/logging"
	"bitbucket.org/pharmaeasyteam/goframework/render"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/database"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/identity"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/keysetmanager"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/badresponse"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/db"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/decryption"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/encryption"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/metadata"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/tokenmanager"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/validator"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/keyset"
	"go.uber.org/zap"
)

func (c *ModuleCrypto) status(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("ok"))
}

func (c *ModuleCrypto) encrypt(w http.ResponseWriter, req *http.Request) {

	// get parsed data
	requestParams, err := validator.ValidateEncryptionRequest(req)
	if err != nil {
		render.JSONWithStatus(w, req, http.StatusBadRequest, badresponse.ExceptionResponse(http.StatusBadRequest, err.Error()))
		return
	}

	// validate identifier
	isAuthorized := identity.AuthenticateRequest(requestParams.Identifier)
	if !isAuthorized {
		render.JSONWithStatus(w, req, http.StatusForbidden, badresponse.ExceptionResponse(http.StatusForbidden, "You are forbidden to perform this action"))
		return
	}

	// authorize encryption levels
	isAuthorized = identity.AuthorizeLevelForEncryption(requestParams)
	if !isAuthorized {
		render.JSONWithStatus(w, req, http.StatusForbidden, badresponse.ExceptionResponse(http.StatusForbidden, "You are forbidden to perform this action"))
		return
	}

	// encrypt data
	encryptedData, err := encryptTokenData(requestParams)
	if err != nil {
		render.JSONWithStatus(w, req, http.StatusInternalServerError, badresponse.ExceptionResponse(http.StatusInternalServerError, "Error encountered while encrypting token data."))
		return
	}

	render.JSON(w, req, encryptedData)
}

func (c *ModuleCrypto) decrypt(w http.ResponseWriter, req *http.Request) {

	// validate request params
	requestParams, err := validator.ValidateDecryptionRequest(req)
	if err != nil {
		render.JSONWithStatus(w, req, http.StatusBadRequest, badresponse.ExceptionResponse(http.StatusBadRequest, err.Error()))
		return
	}

	// validate identifier
	isAuthorized := identity.AuthenticateRequest(requestParams.Identifier)
	if !isAuthorized {
		render.JSONWithStatus(w, req, http.StatusForbidden, badresponse.ExceptionResponse(http.StatusForbidden, "You are forbidden to perform this action"))
		return
	}

	// fetch records
	tokenData, err := getTokenData(requestParams)
	if err != nil {
		render.JSONWithStatus(w, req, http.StatusInternalServerError, badresponse.ExceptionResponse(http.StatusInternalServerError, "Error encountered while fetching token data."))
		return
	}

	// authorize token access
	isAuthorized = identity.AuthorizeTokenAccess(tokenData, requestParams.Identifier)
	if !isAuthorized {
		render.JSONWithStatus(w, req, http.StatusForbidden, badresponse.ExceptionResponse(http.StatusForbidden, "You are forbidden to perform this action"))
		return
	}

	// decrypt data
	decryptedData, err := decryptTokenData(tokenData, requestParams)
	if err != nil {
		render.JSONWithStatus(w, req, http.StatusInternalServerError, badresponse.ExceptionResponse(http.StatusInternalServerError, "Error encountered while decrypting token data."))
		return
	}

	render.JSON(w, req, decryptedData)

}

func (c *ModuleCrypto) getMetaData(w http.ResponseWriter, req *http.Request) {

	// validate request params
	requestParams, err := validator.ValidateMetadataRequest(req)
	if err != nil {
		render.JSONWithStatus(w, req, http.StatusBadRequest, badresponse.ExceptionResponse(http.StatusBadRequest, err.Error()))
		return
	}

	// validate access
	isAuthorized := identity.AuthenticateRequest(requestParams.Identifier)
	if !isAuthorized {
		render.JSONWithStatus(w, req, http.StatusForbidden, badresponse.ExceptionResponse(http.StatusForbidden, "You are forbidden to perform this action"))
		return
	}

	// fetch records
	tokenData, err := database.GetItemsByToken(requestParams.Tokens)
	if err != nil {
		render.JSONWithStatus(w, req, http.StatusForbidden, badresponse.ExceptionResponse(http.StatusForbidden, "Error encountered while fetching token data."))
		return
	}

	// authorize token access
	isAuthorized = identity.AuthorizeTokenAccess(&tokenData, requestParams.Identifier)
	if !isAuthorized {
		render.JSONWithStatus(w, req, http.StatusForbidden, badresponse.ExceptionResponse(http.StatusForbidden, "You are forbidden to perform this action"))
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
		render.JSONWithStatus(w, req, http.StatusBadRequest, badresponse.ExceptionResponse(http.StatusBadRequest, err.Error()))
		return
	}

	// validate access
	isAuthorized := identity.AuthenticateRequest(requestParams.Identifier)
	if !isAuthorized {
		render.JSONWithStatus(w, req, http.StatusForbidden, badresponse.ExceptionResponse(http.StatusForbidden, "You are forbidden to perform this action"))
		return
	}

	// fetch records
	payloadSize := len(requestParams.UpdateParams)
	tokenIDs := make([]string, payloadSize)

	for i := 0; i < payloadSize; i++ {
		tokenIDs[i] = requestParams.UpdateParams[i].Token
	}

	tokenData, err := database.GetItemsByToken(tokenIDs)
	if err != nil {
		render.JSONWithStatus(w, req, http.StatusForbidden, badresponse.ExceptionResponse(http.StatusForbidden, "Error encountered while fetching token data."))
		return
	}

	// authorize token access
	isAuthorized = identity.AuthorizeTokenAccess(&tokenData, requestParams.Identifier)
	if !isAuthorized {
		render.JSONWithStatus(w, req, http.StatusForbidden, badresponse.ExceptionResponse(http.StatusForbidden, "You are forbidden to perform this action"))
		return
	}

	// update metadata
	err = updateMetaItems(requestParams)
	if err != nil {
		render.JSONWithStatus(w, req, http.StatusInternalServerError, badresponse.ExceptionResponse(http.StatusInternalServerError, "Error encountered while updating metadata."+err.Error()))
		return
	}

	render.JSON(w, req, "Metadata updated successfully.")
}

func getTokenData(requestParams *decryption.DecryptRequest) (*map[string]db.TokenData, error) {

	payloadSize := len(requestParams.DecryptRequestData)
	tokenIDs := make([]string, payloadSize)

	for i := 0; i < payloadSize; i++ {
		tokenIDs[i] = requestParams.DecryptRequestData[i].Token
	}

	tokenData, err := database.GetItemsByToken(tokenIDs)
	if err != nil {
		return nil, err
	}

	return &tokenData, nil
}

func encryptTokenData(requestParams *encryption.EncryptRequest) (*encryption.EncryptResponse, error) {
	encryptionResponse := encryption.EncryptResponse{}
	reqParamsData := requestParams.EncryptRequestData

	// get keyset handler
	keyName, keysetHandle, err := keysetmanager.GetKeysetHandlerForEncryption()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(reqParamsData); i++ {
		// encrypt text
		ciphertext, err := dataEncrypt(reqParamsData[i].Content, reqParamsData[i].Salt, keysetHandle)
		if err != nil {
			return nil, err
		}

		dbTokenData := db.TokenData{
			Level:     reqParamsData[i].Level,
			Content:   ciphertext,
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
			Key:       *keyName,
			Metadata:  reqParamsData[i].Metadata,
		}

		token, err := storeEncryptedData(dbTokenData, 1)
		if err != nil {
			return nil, err
		}

		encryptionResponse.ResponseData = append(encryptionResponse.ResponseData,
			encryption.ResponseData{
				ID:    reqParamsData[i].ID,
				Token: tokenmanager.FormatToken(*token),
			})

	}

	return &encryptionResponse, nil
}

func dataEncrypt(data string, salt string, kh *keyset.Handle) ([]byte, error) {
	a, err := aead.New(kh)
	if err != nil {
		logging.GetLogger().Error("Problem in AEAD wrapper generation", zap.Error(err))
		return nil, err
	}

	ct, err := a.Encrypt([]byte(data), []byte(salt))
	if err != nil {
		logging.GetLogger().Error("Problem in Encryption", zap.Error(err))
		return nil, err
	}

	return ct, nil
}
func dataDecrypt(cipherText []byte, salt string, kh *keyset.Handle) (*string, error) {
	a, err := aead.New(kh)
	if err != nil {
		logging.GetLogger().Error("Error encountered while initializing aead handler", zap.Error(err))
		return nil, err
	}

	pt, err := a.Decrypt(cipherText, []byte(salt))
	if err != nil {
		logging.GetLogger().Error("Error encountered while decrypting data", zap.Error(err))
		return nil, err
	}

	plainText := string(pt)

	return &plainText, nil
}

func storeEncryptedData(dbTokenData db.TokenData, attempt int) (*string, error) {
	dbTokenData.TokenID = tokenmanager.Uniquetoken()
	err := database.PutItem(dbTokenData)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				logging.GetLogger().Error("Token ID clash detected.", zap.Error(err))
				if attempt > 3 {
					logging.GetLogger().Error("Token ID clash exceeded attempt threshold.", zap.Error(err))
					return nil, err
				}
				attempt++
				storeEncryptedData(dbTokenData, attempt)
			default:
				return nil, err
			}
		}
	}

	return &dbTokenData.TokenID, nil
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
		decryptedText, err := dataDecrypt(dbTokenData.Content, reqParamsData[i].Salt, kh)
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

func updateMetaItems(requestParams *metadata.MetaUpdateRequest) error {

	payloadSize := len(requestParams.UpdateParams)
	for i := 0; i < payloadSize; i++ {

		err := database.UpdateMetadataByToken(requestParams.UpdateParams[i].Token, requestParams.UpdateParams[i].Metadata)
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
