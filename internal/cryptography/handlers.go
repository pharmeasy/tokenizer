package cryptography

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/pharmaeasyteam/goframework/logging"
	"bitbucket.org/pharmaeasyteam/goframework/render"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/database"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/keysetmanager"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/badresponse"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/datadecryption"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/db"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/encryption"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/metadata"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/request/decryption"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/tokenmanager"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/keyset"
	"go.uber.org/zap"
)

func DataDecrypt(cipherText []byte, salt string, kh *keyset.Handle) (*string, error) {
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

func validateEncryptionRequest(req *http.Request) (*encryption.EncryptRequest, error) {
	decoder := json.NewDecoder(req.Body)
	test := encryption.EncryptRequest{}
	err := decoder.Decode(&test)
	if err != nil {
		logging.GetLogger().Error("Problem in input params", zap.Error(err))
		return nil, err
	}
	if test.RequestID == "" {
		logging.GetLogger().Error("Problem in input params", zap.Error(err))
		return nil, err
	}
	if test.Identifier == "" {
		logging.GetLogger().Error("Problem in input params", zap.Error(err))
		return nil, err
	}
	if test.Level < 1 {
		logging.GetLogger().Error("Problem in input params", zap.Error(err))
		return nil, err
	}
	for _, v := range test.RequestData {
		if v.Salt == "" || v.Content == "" {
			return nil, err
		}
	}
	return &test, nil
}

func authorizeTokenAccessForEncryption(identifier string, level int) bool {
	var list = map[string]int{
		"iron":  1,
		"oms":   2,
		"alloy": 3,
	}

	src := list[identifier]
	if src < level {
		return false
	}

	return true
}

func (c *ModuleCrypto) status(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("ok"))
}

func (c *ModuleCrypto) encrypt(w http.ResponseWriter, req *http.Request) {

	//get parsed data
	requestParams, err := validateEncryptionRequest(req)
	if err != nil {
		render.JSONWithStatus(w, req, http.StatusBadRequest, badresponse.ExceptionResponse(http.StatusBadRequest, err.Error()))
		return
	}

	// validate identifier
	isAuthorized := authorizeRequest(requestParams.Identifier)
	if !isAuthorized {
		render.JSONWithStatus(w, req, http.StatusForbidden, badresponse.ExceptionResponse(http.StatusForbidden, "You are forbidden to perform this action"))
		return
	}

	// validate level
	isAuthorized = authorizeTokenAccessForEncryption(requestParams.Identifier, requestParams.Level)
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
	requestParams, err := validateDecryptionRequest(req)
	if err != nil {
		render.JSONWithStatus(w, req, http.StatusBadRequest, badresponse.ExceptionResponse(http.StatusBadRequest, err.Error()))
		return
	}

	// validate access
	isAuthorized := authorizeRequest(requestParams.Identifier)
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
	isAuthorized = authorizeTokenAccess(tokenData, requestParams.Level)
	if !isAuthorized {
		render.JSONWithStatus(w, req, http.StatusForbidden, badresponse.ExceptionResponse(http.StatusForbidden, "You are forbidden to perform this action"))
		return
	}

	// decrypt data
	decryptedData, err := decryptTokenData(tokenData, requestParams)
	if err != nil {
		render.JSONWithStatus(w, req, http.StatusForbidden, badresponse.ExceptionResponse(http.StatusInternalServerError, "Error encountered while decrypting token data."))
		return
	}

	render.JSON(w, req, decryptedData)

}

func validateDecryptionRequest(req *http.Request) (*decryption.Request, error) {

	decoder := json.NewDecoder(req.Body)
	params := decryption.Request{}
	err := decoder.Decode(&params)
	if err != nil {
		logging.GetLogger().Error("Unable to decode decryption request params.", zap.Error(err))
		return nil, err
	}

	genericError := errors.New("Invalid request parameters passed.")

	if params.Level < 1 {
		logging.GetLogger().Error("Invalid Level.", zap.Error(err))
		return nil, genericError
	}

	if params.Identifier == "" {
		logging.GetLogger().Error("Identifier is empty.", zap.Error(err))
		return nil, genericError
	}

	for i := 0; i < len(params.Data); i++ {
		if params.Data[i].Token == "" {
			logging.GetLogger().Error("Empty token passed.", zap.Error(err))
			return nil, genericError
		}

		if params.Data[i].Salt == "" {
			logging.GetLogger().Error("Empty salt passed.", zap.Error(err))
			return nil, genericError
		}
	}

	return &params, nil
}

func getTokenData(requestParams *decryption.Request) (*map[string]db.TokenData, error) {

	payloadSize := len(requestParams.Data)
	tokenIDs := make([]string, payloadSize)

	for i := 0; i < payloadSize; i++ {
		tokenIDs[i] = requestParams.Data[i].Token
	}

	tokenData, err := database.GetItemsByToken(tokenIDs)
	if err != nil {
		return nil, err
	}

	return &tokenData, nil
}

func authorizeRequest(accessToken string) bool {

	accessTokenArr := []string{"iron", "oms", "alloy"}
	for _, v := range accessTokenArr {
		if v == accessToken {
			return true
		}
	}
	return false
}

func authorizeTokenAccess(tokenData *map[string]db.TokenData, level int) bool {

	for _, v := range *tokenData {
		if level < v.Level {
			return false
		}
	}
	return true
}

func encryptTokenData(requestParams *encryption.EncryptRequest) (*encryption.EncryptResponse, error) {
	encryptionResponse := encryption.EncryptResponse{}
	reqParamsData := requestParams.RequestData

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

		fmt.Println(string(ciphertext))

		dbTokenData := db.TokenData{
			Level:     requestParams.Level,
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

//DataEncrypt returns the cipher text
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

func decryptTokenData(tokenData *map[string]db.TokenData, requestParams *decryption.Request) (*datadecryption.DecryptionResponse, error) {
	decryptionResponse := datadecryption.DecryptionResponse{}
	reqParamsData := requestParams.Data
	for i := 0; i < len(reqParamsData); i++ {
		token := reqParamsData[i].Token
		dbTokenData := (*tokenData)[token]

		// select keyset
		kh, err := keysetmanager.GetKeysetHandlerForDecryption(dbTokenData.Key)
		if err != nil {
			return nil, err
		}

		// decrypt with salt
		decryptedText, err := DataDecrypt(dbTokenData.Content, reqParamsData[i].Salt, kh)
		if err != nil {
			return nil, err
		}

		decryptionResponse.DecryptionResponseData = append(decryptionResponse.DecryptionResponseData,
			datadecryption.DecryptionResponseData{
				Token:    tokenmanager.FormatToken(token),
				Content:  *decryptedText,
				Metadata: dbTokenData.Metadata,
			})
	}

	return &decryptionResponse, nil
}

func (c *ModuleCrypto) updateMetadata(w http.ResponseWriter, req *http.Request) {

	// validate request params
	requestParams, err := validateMetadataUpdateRequest(req)
	if err != nil {
		render.JSONWithStatus(w, req, http.StatusBadRequest, badresponse.ExceptionResponse(http.StatusBadRequest, err.Error()))
		return
	}

	// validate access
	isAuthorized := authorizeRequest(requestParams.Identifier)
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
	isAuthorized = authorizeTokenAccess(&tokenData, requestParams.Level)
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

func validateMetadataUpdateRequest(req *http.Request) (*metadata.MetaUpdateRequest, error) {

	decoder := json.NewDecoder(req.Body)
	params := metadata.MetaUpdateRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		logging.GetLogger().Error("Unable to decode metadata request params.", zap.Error(err))
		return nil, err
	}

	genericError := errors.New("Invalid request parameters passed.")

	if params.Level < 1 {
		logging.GetLogger().Error("Invalid Level.", zap.Error(err))
		return nil, genericError
	}

	if params.Identifier == "" {
		logging.GetLogger().Error("Identifier is empty.", zap.Error(err))
		return nil, genericError
	}

	for i := 0; i < len(params.UpdateParams); i++ {
		if params.UpdateParams[i].Token == "" {
			logging.GetLogger().Error("Empty token passed.", zap.Error(err))
			return nil, genericError
		}

		if params.UpdateParams[i].Metadata == "" {
			logging.GetLogger().Error("Empty metadata passed.", zap.Error(err))
			return nil, genericError
		}
	}

	return &params, nil
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

func (c *ModuleCrypto) getMetaData(w http.ResponseWriter, req *http.Request) {

	// validate request params
	requestParams, err := validateMetadataRequest(req)
	if err != nil {
		render.JSONWithStatus(w, req, http.StatusBadRequest, badresponse.ExceptionResponse(http.StatusBadRequest, err.Error()))
		return
	}

	// validate access
	isAuthorized := authorizeRequest(requestParams.Identifier)
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
	isAuthorized = authorizeTokenAccess(&tokenData, requestParams.Level)
	if !isAuthorized {
		render.JSONWithStatus(w, req, http.StatusForbidden, badresponse.ExceptionResponse(http.StatusForbidden, "You are forbidden to perform this action"))
		return
	}

	// return metadata
	metadataResponse := getMetaItems(tokenData)
	render.JSON(w, req, metadataResponse)
}

func validateMetadataRequest(req *http.Request) (*metadata.MetaRequest, error) {

	decoder := json.NewDecoder(req.Body)
	params := metadata.MetaRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		logging.GetLogger().Error("Unable to decode metadata request params.", zap.Error(err))
		return nil, err
	}

	genericError := errors.New("Invalid request parameters passed.")

	if params.Level < 1 {
		logging.GetLogger().Error("Invalid Level.", zap.Error(err))
		return nil, genericError
	}

	if params.Identifier == "" {
		logging.GetLogger().Error("Identifier is empty.", zap.Error(err))
		return nil, genericError
	}

	for i := 0; i < len(params.Tokens); i++ {
		if params.Tokens[i] == "" {
			logging.GetLogger().Error("Empty token passed.", zap.Error(err))
			return nil, genericError
		}
	}

	return &params, nil
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
