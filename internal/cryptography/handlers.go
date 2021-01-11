package cryptography

import (
	"log"
	"net/http"

	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/keyset"
	"github.com/google/tink/go/tink"

	"bitbucket.org/pharmaeasyteam/goframework/render"
	"bitbucket.org/pharmaeasyteam/goframework/models"	
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/uuidmodule"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/jsonparser"
	models2 "bitbucket.org/pharmaeasyteam/tokenizer/internal/models"

)

//Encryption function. Key need to be provided from KMS 
func DataEncrypt(data string , key string) ([]byte , []byte , tink.AEAD) {

	// AEAD primitive
	kh, err := keyset.NewHandle(aead.AES256GCMKeyTemplate())
	if err != nil {
		log.Fatal(err)
	}

	a, err := aead.New(kh)
	if err != nil {
		log.Fatal(err)
	}

	ct, err := a.Encrypt([]byte (data) , []byte (key))
	if err != nil {
		log.Fatal(err)
	}

	return ct, []byte (key) , a
} 


func (c *ModuleCrypto) getTokens(w http.ResponseWriter , req *http.Request) {
	//UUID token
	uniqueId :=  uuidmodule.Uniquetoken()
	
	//get parsed data
	content , source := jsonparser.ParseData(req)
	
	// Encryption task
	key := "private_key"
	data , _ , _  := DataEncrypt(content , key)
	
	//datastore object
	dataStore := models2.DataStore{}
	dataStore.TokenID = uniqueId
	dataStore.EncryptedData = data
	dataStore.Source = source
	dataStore.EncryptionMode = 0 // need to figure out the logic
	dataStore.Severity = 1 		 // need to figure out the logic

	/*
		dataStore object will be stored in dynamoDB. need to figure out the logic

	*/
	
	//Return the token id to the client 
	tokenString := "\"token\" : " + uniqueId.String()
	w.Write([]byte(tokenString))
}


// After decryption 
func (c *ModuleCrypto) getData(w http.ResponseWriter, req *http.Request) {
	render.JSON(w, req, models.Response{Msg: "Wait for Implementation"})
}