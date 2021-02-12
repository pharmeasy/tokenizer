package cryptography

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/metadata"

	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/db"

	"bitbucket.org/pharmaeasyteam/tokenizer/internal/database"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/datadecryption"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/db"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/request/decryption"

	"bitbucket.org/pharmaeasyteam/goframework/logging"
	"bitbucket.org/pharmaeasyteam/goframework/render"

	//"bitbucket.org/pharmaeasyteam/tokenizer/internal/database"
	kms "bitbucket.org/pharmaeasyteam/tokenizer/internal/kms/aws"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/badresponse"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/request/encryption"
	encryption2 "bitbucket.org/pharmaeasyteam/tokenizer/internal/models/response/encryption"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/uuidmodule"
	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/keyset"

	"go.uber.org/zap"
)

//DataEncrypt returns the cipher text
func DataEncrypt(data string, salt string, kh *keyset.Handle) []byte {
	a, err := aead.New(kh)
	if err != nil {
		logging.GetLogger().Error("Problem in AEAD wrapper generation", zap.Error(err))
	}

	ct, err := a.Encrypt([]byte(data), []byte(salt))
	if err != nil {
		logging.GetLogger().Error("Problem in Encryption", zap.Error(err))
	}

	return ct
}

func DataDecrypt(cipherText string, salt string, kh *keyset.Handle) (*string, error) {
	a, err := aead.New(kh)
	if err != nil {
		logging.GetLogger().Error("Error encountered while initializing aead handler", zap.Error(err))
		return nil, err
	}

	pt, err := a.Decrypt([]byte(cipherText), []byte(salt))
	if err != nil {
		logging.GetLogger().Error("Error encountered while decrypting data", zap.Error(err))
		return nil, err
	}

	plainText := string(pt)

	return &plainText, nil
}

// DataEncryptWrapper is the main encryption function which gives us the response
func DataEncryptWrapper(data []encryption.Data, kh *keyset.Handle) encryption2.Response {
	var response encryption2.Response
	temp := []encryption2.Data{}
	// for i := 0; i < len(data); i++ {
	// 	uniqueID := uuidmodule.Uniquetoken()
	// 	cipherText := DataEncrypt(data[i].Content, data[i].Salt, kh)
	// 	temp = append(temp, encryption2.Data{
	// 		ID:     data[i].ID,
	// 		Token:  uniqueID,
	// 		Cipher: cipherText,
	// 	})
	// }

	for _, v := range data {
		uniqueID := uuidmodule.Uniquetoken()
		cipherText := DataEncrypt(v.Content, v.Salt, kh)
		temp = append(temp, encryption2.Data{
			ID:     v.ID,
			Token:  uniqueID,
			Cipher: cipherText,
		})
	}
	response.Data = temp
	return response
}

func validateAndParseRequest(req *http.Request) ([]encryption.Data, string, int) {

	decoder := json.NewDecoder(req.Body)
	test := encryption.Request{}
	err := decoder.Decode(&test)
	if err != nil {
		logging.GetLogger().Error("Problem in input params", zap.Error(err))
	}

	// input validations
	requestId := test.RequestID
	identifier := test.Identifier
	level := test.Level

	if requestId == "" || identifier == "" || level < 1 {
		logging.GetLogger().Error("Problem in input params", zap.Error(err))
		return nil, "", -1
	}

	for i := 0; i < len(test.Data); i++ {
		if test.Data[i].Content == "" || test.Data[i].Salt == "" {
			return nil, "", -1
		}
	}
	return test.Data, test.Identifier, test.Level

}

// severity mapper function
func validateMapper(identifier string, level int) int {
	var list = map[string]int{
		"iron":  1,
		"oms":   2,
		"alloy": 3,
	}

	src := list[identifier]
	if src < level {
		return http.StatusForbidden
	}
	return http.StatusOK
}

// getTokens ...
func (c *ModuleCrypto) getTokens(w http.ResponseWriter, req *http.Request) {
	// if len(kms.DecryptedKeysetMap) != 0 && len(kms.KeysetArr) != 0 {
	// 	kms.KeysetArr = kms.KeysetName()
	// 	kms.DecryptKeyset()
	// }
	kms.DecryptKeyset()
	KeysetArr := kms.KeysetName(kms.DecryptedKeysetMap)
	fmt.Println(len(KeysetArr))
	fmt.Println(len(kms.DecryptedKeysetMap))
	keysetName := kms.SelectKeyset(KeysetArr)
	keysetHandler := kms.DecryptedKeysetMap[keysetName]
	//get parsed data
	content, identifier, level := validateAndParseRequest(req)
	// // check severity mapper
	responseCode := validateMapper(identifier, level)

	// malformed request
	if content == nil {
		render.JSONWithStatus(w, req, http.StatusBadRequest, badresponse.ExceptionResponse(http.StatusBadRequest, "Your request is malformed"))
		return
	}

	// forbidden request
	if responseCode == 403 {
		render.JSONWithStatus(w, req, http.StatusForbidden, badresponse.ExceptionResponse(http.StatusForbidden, "You are forbidden to perform this action"))
		return
	}

	//Encryption
	fmt.Println(len(content))
	response := DataEncryptWrapper(content, keysetHandler)

	// store record
	//item := make([]db.TokenData{} , len(response.Data))
	item := make([]db.TokenData, len(response.Data))
	//itenm4 := make([]db.TokenData)

	// for i, v := range response.Data {
	// 	item[i].Content = string(v.Cipher)
	// 	item[i].Key = v.Token
	// 	item[i].Level = level
	// 	item[i].Meta = content[i].MetaData
	// 	item[i].TokenID = strconv.Itoa(content[i].ID)
	// 	item[i].CreatedAt = time.Now().String()
	// 	item[i].UpdatedAt = time.Now().String()
	// 	database.PutItem(item[i])
	// }

	for _, v := range response.Data {
		item = append(item, db.TokenData{
			Content:   string(v.Cipher),
			Key:       v.Token,
			Level:     level,
			Meta:      "",
			TokenID:   "1",
			CreatedAt: time.Now().String(),
			UpdatedAt: time.Now().String(),
		})
	}

	render.JSON(w, req, item)
	w.Write([]byte(keysetName))
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
		render.JSONWithStatus(w, req, http.StatusForbidden, badresponse.ExceptionResponse(http.StatusInternalServerError, "Error encounteed while fetching token data."))
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
		render.JSONWithStatus(w, req, http.StatusForbidden, badresponse.ExceptionResponse(http.StatusInternalServerError, "Error encounteed while decrypting token data."))
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
	return true
}

func authorizeTokenAccess(*map[string]db.TokenData, int) bool {
	return true
}

func decryptTokenData(tokenData *map[string]db.TokenData, requestParams *decryption.Request) (*datadecryption.DecryptionResponse, error) {
	decryptionResponse := datadecryption.DecryptionResponse{}
	reqParamsData := requestParams.Data
	for i := 0; i < len(reqParamsData); i++ {
		token := reqParamsData[i].Token
		dbTokenData := (*tokenData)[token]

		// select keyset
		kh, err := getKeysetHandle(dbTokenData.Key)
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
				Token:    token,
				Content:  *decryptedText,
				Metadata: dbTokenData.Metadata,
			})
	}

	return &decryptionResponse, nil
}

func getKeysetHandle(handleName string) (*keyset.Handle, error) {
	if kh, ok := kms.DecryptedKeysetMap[handleName]; ok {
		return kh, nil
	}
	err := errors.New("Something went wrong while processing your request")
	logging.GetLogger().Error("Valid keyset not found for handle name."+handleName, zap.Error(err))

	return nil, err
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
		render.JSONWithStatus(w, req, http.StatusForbidden, badresponse.ExceptionResponse(http.StatusForbidden, "Error encounteed while fetching token data."))
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
		render.JSONWithStatus(w, req, http.StatusForbidden, badresponse.ExceptionResponse(http.StatusForbidden, "Error encounteed while fetching token data."))
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
				Token:    dbTokenData.TokenID,
				Metadata: dbTokenData.Metadata,
			})
	}

	return &metaResponse
}
