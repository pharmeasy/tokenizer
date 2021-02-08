package cryptography

import (
	"encoding/json"
	"net/http"

	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/keyset"

	"bitbucket.org/pharmaeasyteam/goframework/logging"
	"bitbucket.org/pharmaeasyteam/goframework/models"
	"bitbucket.org/pharmaeasyteam/goframework/render"
	kms "bitbucket.org/pharmaeasyteam/tokenizer/internal/kms/aws"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/badresponse"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/request/encryption"
	encryption2 "bitbucket.org/pharmaeasyteam/tokenizer/internal/models/response/encryption"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/uuidmodule"

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
	source := test.Source
	level := test.Level

	if requestId < 1 || source == "" || level < 1 {
		logging.GetLogger().Error("Problem in input params", zap.Error(err))
		return nil, "", -1
	}

	for i := 0; i < len(test.Data); i++ {
		if test.Data[i].Content == "" || test.Data[i].Salt == "" {
			return nil, "", -1
		}
	}
	return test.Data, test.Source, test.Level

}

// severity mapper function
func validateMapper(source string, level int) int {
	var list = map[string]int{
		"iron":  1,
		"oms":   2,
		"alloy": 3,
	}

	src := list[source]
	if src < level {
		return http.StatusForbidden
	}
	return http.StatusOK
}

// getTokens ...
func (c *ModuleCrypto) getTokens(w http.ResponseWriter, req *http.Request) {

	fileName := "/Users/riddhimanparasar/tokenizer/internal/cryptography/keyset1.json"
	kh := kms.GetKeysets(fileName)

	//get parsed data
	content, source, level := validateAndParseRequest(req)
	// // check severity mapper
	responseCode := validateMapper(source, level)

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
	response := DataEncryptWrapper(content, kh)
	render.JSON(w, req, response)
}

// getData ...
func (c *ModuleCrypto) getData(w http.ResponseWriter, req *http.Request) {
	render.JSON(w, req, models.Response{Msg: "Wait for Implementation"})
}

func (c *ModuleCrypto) updateMeta(w http.ResponseWriter, req *http.Request) {

}
