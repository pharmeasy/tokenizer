package cryptography

import (
	"net/http"

	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/keyset"
	"github.com/google/tink/go/tink"

	"bitbucket.org/pharmaeasyteam/goframework/logging"
	"bitbucket.org/pharmaeasyteam/goframework/models"
	"bitbucket.org/pharmaeasyteam/goframework/render"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/jsonparser"

	//"bitbucket.org/pharmaeasyteam/tokenizer/internal/kms/aws"
	models2 "bitbucket.org/pharmaeasyteam/tokenizer/internal/models"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/uuidmodule"

	"go.uber.org/zap"
)

//DataEncrypt ...
func DataEncrypt(data string, key string) ([]byte, tink.AEAD) {

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

	return ct, a
}

// getTokens ...
func (c *ModuleCrypto) getTokens(w http.ResponseWriter, req *http.Request) {
	//UUID token
	uniqueID := uuidmodule.Uniquetoken()

	// err := jsonparser.CheckAllParams(req)
	// log.Fatal(err)
	//jsonparser.CheckAllParams(req)

	//get parsed data
	content, source := jsonparser.ParseData(req)

	// Encryption task
	key := "private_key"
	data, _ := DataEncrypt(content, key)

	//datastore object
	dataStore := models2.DataStore{}
	dataStore.TokenID = uniqueID
	dataStore.EncryptedData = data
	dataStore.Source = source
	dataStore.EncryptionMode = 0 // need to figure out the logic
	dataStore.Severity = 1       // need to figure out the logic

	/*
		dataStore object will be stored in dynamoDB. need to figure out the logic

	*/

	//Return the token id to the client
	tokenString := "\"token\" : " + uniqueID.String()
	w.Write([]byte(tokenString))
	//w.Write(data)
}

// getData ...
func (c *ModuleCrypto) getData(w http.ResponseWriter, req *http.Request) {
	render.JSON(w, req, models.Response{Msg: "Wait for Implementation"})
}
