package cryptography

import (
	"net/http"

	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/keyset"

	"bitbucket.org/pharmaeasyteam/goframework/logging"
	"bitbucket.org/pharmaeasyteam/goframework/models"
	"bitbucket.org/pharmaeasyteam/goframework/render"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/jsonparser"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/request/encryption"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/response"
	encryption2 "bitbucket.org/pharmaeasyteam/tokenizer/internal/models/response/encryption"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/uuidmodule"

	"go.uber.org/zap"
)

//DataEncrypt ...
func DataEncrypt(data string, key string) string {

	// AEAD primitive
	kh, err := keyset.NewHandle(aead.AES256GCMKeyTemplate())
	if err != nil {
		logging.GetLogger().Error("Problem in keyset generation", zap.Error(err))
	}

	a, err := aead.New(kh)
	if err != nil {
		logging.GetLogger().Error("Problem in AEAD wrapper generation", zap.Error(err))
	}

	ct, err := a.Encrypt([]byte(data), []byte(key))
	if err != nil {
		logging.GetLogger().Error("Problem in Encryption", zap.Error(err))
	}

	return string(ct)
}

// DataEncryptWrapper is the main encryption function which gives us the response
func DataEncryptWrapper(data []encryption.Data) encryption2.Response {
	var response encryption2.Response
	temp := []encryption2.Data{}
	for i := 0; i < len(data); i++ {
		uniqueID := uuidmodule.Uniquetoken()
		temp = append(temp, encryption2.Data{
			Token:   uniqueID,
			Content: DataEncrypt(data[i].Content, "key"),
		})
	}
	response.Data = temp
	return response
}

// getTokens ...
func (c *ModuleCrypto) getTokens(w http.ResponseWriter, req *http.Request) {

	//get parsed data
	content := jsonparser.ParseData(req)

	// content nil means that validation check logic from parse.go
	if content == nil {
		render.JSONWithStatus(w, req, http.StatusBadRequest, response.ExceptionResponse(http.StatusBadRequest, "Your request is malformed"))
		return
	}

	// output the response
	response := DataEncryptWrapper(content)
	render.JSON(w, req, response)

	//datastore object
	// dataStore := models2.DataStore{}
	// dataStore.TokenID = uniqueID
	// dataStore.EncryptedData = data
	// dataStore.Source = source
	// dataStore.EncryptionMode = 0 // need to figure out the logic
	// dataStore.Severity = 1       // need to figure out the logic

	/*
		dataStore object will be stored in dynamoDB. need to figure out the logic

	*/
}

// getData ...
func (c *ModuleCrypto) getData(w http.ResponseWriter, req *http.Request) {
	render.JSON(w, req, models.Response{Msg: "Wait for Implementation"})
}
