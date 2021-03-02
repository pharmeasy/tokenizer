package cryptography

import (
	"net/http"
	"time"
    "fmt"
	"bitbucket.org/pharmaeasyteam/goframework/render"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/database"
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
	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/keyset"
)

func (c *ModuleCrypto) status(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("ok"))
}

func (c *ModuleCrypto) encrypt(w http.ResponseWriter, req *http.Request) {

	// get parsed data
	requestParams, err := validator.ValidateEncryptionRequest(req)
	if err != nil {
		errormanager.RenderEncryptionErrorResponse(w, req, http.StatusBadRequest,
			errormanager.SetEncryptionError(requestParams, err, http.StatusBadRequest))
		return
	}

	// validate identifier
	isAuthenticated := identity.AuthenticateRequest(requestParams.Identifier)
	if !isAuthenticated {
		errormanager.RenderEncryptionErrorResponse(w, req, http.StatusForbidden,
			errormanager.SetEncryptionError(requestParams, nil, http.StatusForbidden))
		return
	}

	// authorize encryption levels
	isAuthorized := identity.AuthorizeLevelForEncryption(requestParams)
	if !isAuthorized {
		errormanager.RenderEncryptionErrorResponse(w, req, http.StatusForbidden,
			errormanager.SetEncryptionError(requestParams, nil, http.StatusForbidden))
		return
	}

	// encrypt data
	encryptedData, err := encryptTokenData(requestParams)
	if err != nil {
		errormanager.RenderEncryptionErrorResponse(w, req, http.StatusInternalServerError,
			errormanager.SetEncryptionError(requestParams, err, http.StatusInternalServerError))
		return
	}

	render.JSON(w, req, encryptedData)
}

func (c *ModuleCrypto) decrypt(w http.ResponseWriter, req *http.Request) {

	// validate request params
	requestParams, err := validator.ValidateDecryptionRequest(req)
	if err != nil {
		errormanager.RenderDecryptionErrorResponse(w, req, http.StatusBadRequest,
			errormanager.SetDecryptionError(requestParams, err, http.StatusBadRequest))
		return
	}

	// validate identifier
	isAuthenticated := identity.AuthenticateRequest(requestParams.Identifier)
	if !isAuthenticated {
		errormanager.RenderDecryptionErrorResponse(w, req, http.StatusForbidden,
			errormanager.SetDecryptionError(requestParams, nil, http.StatusForbidden))
		return
	}
    fmt.Println("before fetching")
	// fetch records
	tokenData, err := getTokenData(requestParams)
	if err != nil {
		errormanager.RenderDecryptionErrorResponse(w, req, http.StatusInternalServerError,
			errormanager.SetDecryptionError(requestParams, err, http.StatusInternalServerError))
		return
	}
    fmt.Println(tokenData)

	// authorize token access
	isAuthorized := identity.AuthorizeTokenAccess(tokenData, requestParams.Identifier)
	if !isAuthorized {
		errormanager.RenderDecryptionErrorResponse(w, req, http.StatusForbidden,
			errormanager.SetDecryptionError(requestParams, nil, http.StatusForbidden))
		return
	}

	// decrypt data
	decryptedData, err := decryptTokenData(tokenData, requestParams)
	if err != nil {
		errormanager.RenderDecryptionErrorResponse(w, req, http.StatusInternalServerError,
			errormanager.SetDecryptionError(requestParams, err, http.StatusInternalServerError))
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
	tokenData, err := database.GetItemsByToken(requestParams.Tokens)
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

	tokenData, err := database.GetItemsByToken(tokenIDs)
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

	// update metadata
	err = updateMetaItems(requestParams)
	if err != nil {
		errormanager.RenderUpdateMetadataErrorResponse(w, req, http.StatusInternalServerError,
			errormanager.SetUpdateMetadataError(requestParams, err, http.StatusInternalServerError))
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

	//tokenData, err := database.GetItemsByToken(tokenIDs)
	tokenData, err := database.GetItemsByTokenInBatch(tokenIDs)
	
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

		token, err := storeEncryptedData(dbTokenData)
		if err != nil {
			return nil, err
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

func storeEncryptedData(dbTokenData db.TokenData) (string, error) {
	attempts := 0
	var err error

	for attempts < 3 {
		attempts++
		dbTokenData.TokenID = tokenmanager.Uniquetoken()
		err = database.PutItem(dbTokenData)
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

func integrityChecker(cipherText []byte, plainText string, salt string, kh *keyset.Handle) bool {
	pt, _ := dataDecryptAEAD(cipherText, salt, kh)
	if *pt == plainText {
		return true
	}
	return false
}
