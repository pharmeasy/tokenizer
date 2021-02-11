package cryptography

import (
	"encoding/json"
	"fmt"
	"errors"
	"net/http"

	"bitbucket.org/pharmaeasyteam/tokenizer/internal/database"
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

var keysetMap = map[int]string{
	0: "keyset1.json",
	1: "keyset2.json",
	2: "keyset3.json",
	3: "keyset4.json",
}

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

// DataEncryptWrapper is the main encryption function which gives us the response
func DataEncryptWrapper(data []encryption.Data, kh *keyset.Handle) encryption2.Response {
	var response encryption2.Response
	temp := []encryption2.Data{}
	for i := 0; i < len(data); i++ {
		uniqueID := uuidmodule.Uniquetoken()
		cipherText := DataEncrypt(data[i].Content, data[i].Salt, kh)
		temp = append(temp, encryption2.Data{
			ID:     data[i].ID,
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
	response := DataEncryptWrapper(content, keysetHandler)
	render.JSON(w, req, response)
	w.Write([]byte(keysetName))
}

func (c *ModuleCrypto) decrypt(w http.ResponseWriter, req *http.Request) {

	// validate request params
	dataParams, tokenIDs, err := validateDecryptionRequest(req)
	if err != nil {
		render.JSONWithStatus(w, req, http.StatusBadRequest, badresponse.ExceptionResponse(http.StatusBadRequest, err.Error()))
		return
	}

	// fetch records
	tokenData, err := database.GetItemsByToken(tokenIDs)
	if err != nil {
		render.JSONWithStatus(w, req, http.StatusForbidden, badresponse.ExceptionResponse(http.StatusInternalServerError, "Error encounteed while fetching token data."))
		return
	}

	// validate access
	isAuthorized := authorizeRequest(dataParams.Identifier, dataParams.Level)
	if !isAuthorized {
		render.JSONWithStatus(w, req, http.StatusForbidden, badresponse.ExceptionResponse(http.StatusForbidden, "You are forbidden to perform this action"))
		return
	}

	// decrypt records (Tink + Salt)

	// decorate response


	render.JSON(w, req, tokenData)

}

func validateDecryptionRequest(req *http.Request) (decryption.Request, []string, error) {

	decoder := json.NewDecoder(req.Body)
	params := decryption.Request{}
	err := decoder.Decode(&params)
	if err != nil {
		logging.GetLogger().Error("Unable to decode decryption request params.", zap.Error(err))
		return params, nil, err
	}

	genericError := errors.New("Invalid request parameters passed.")

	if params.Level < 1 {
		logging.GetLogger().Error("Invalid Level.", zap.Error(err))
		return params, nil, genericError
	}

	if params.Identifier == "" {
		logging.GetLogger().Error("Identifier is empty.", zap.Error(err))
		return params, nil, genericError
	}

	payloadSize := len(params.Data)
	tokenIDs := make([]string, payloadSize)

	for i := 0; i < payloadSize; i++ {
		tokenIDs[i] = params.Data[i].Token
		if params.Data[i].Token == "" {
			logging.GetLogger().Error("Empty token passed.", zap.Error(err))
			return params, nil, genericError
		}

		if params.Data[i].Salt == "" {
			logging.GetLogger().Error("Empty salt passed.", zap.Error(err))
			return params, nil, genericError
		}
	}

	return params, tokenIDs, nil
}

func authorizeRequest(accessToken string, level int) bool {
	return true
}

func (c *ModuleCrypto) updateMeta(w http.ResponseWriter, req *http.Request) {

}
