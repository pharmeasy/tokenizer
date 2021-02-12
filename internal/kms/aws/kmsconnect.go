package kms

import (
	//"bytes"
	//"fmt"
	//"io/ioutil"

	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/pharmaeasyteam/goframework/logging"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/keysetmodel"
	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/core/registry"

	//"github.com/google/tink/go/insecurecleartextkeyset"
	"github.com/google/tink/go/integration/awskms"
	"github.com/google/tink/go/keyset"

	"go.uber.org/zap"
)

var DecryptedKeysetMap = make(map[string]*keyset.Handle)



func destringify(str string) *strings.Reader {
	var keyset keysetmodel.EncryptedKeyset
	s, err := strconv.Unquote(str)
	if err != nil {
		logging.GetLogger().Error("Problem in destrinify", zap.Error(err))
	}
	err = json.Unmarshal([]byte(s), &keyset)
	if err != nil {
		logging.GetLogger().Error("Problem in Unmarshal", zap.Error(err))
	}
	JsonKeyset, _ := json.Marshal(keyset)
	myReader := strings.NewReader(string(JsonKeyset))
	return myReader
}

func LoadKeyset() map[string]*strings.Reader {
	fileName, _ := ioutil.ReadFile("/Users/riddhimanparasar/tokenizer/keysetmap.json")
	//fmt.Println(string(fileName))
	m := make(map[string]string)
	err := json.Unmarshal(fileName, &m)
	if err != nil {
		logging.GetLogger().Error("Problem in unmarshal", zap.Error(err))
	}
	keysetmap := make(map[string]*strings.Reader)

	for k, v := range m {
		keysetmap[k] = destringify(strconv.Quote(v))
	}
	return keysetmap
}

func DecryptKeyset() {
	keysetmap := make(map[string]*strings.Reader)
	keysetmap = LoadKeyset()

	awsclient, err := awskms.NewClient("aws-kms://arn:aws:kms:ap-south-1:127603365779:key/8d853831-94e6-4ac7-a0c7-3e2795e9715b")
	if err != nil {
		logging.GetLogger().Error("Problem in AEAD backend generation", zap.Error(err))
	}
	registry.RegisterKMSClient(awsclient)
	backend, err := awsclient.GetAEAD("aws-kms://arn:aws:kms:ap-south-1:127603365779:key/8d853831-94e6-4ac7-a0c7-3e2795e9715b")
	if err != nil {
		logging.GetLogger().Error("Problem in AEAD backend generation", zap.Error(err))
	}
	masterKey := aead.NewKMSEnvelopeAEAD2(aead.AES128GCMKeyTemplate(), backend)

	for k, v := range keysetmap {
		kh1 := keyset.NewJSONReader(v)
		kh, err := keyset.Read(kh1, masterKey)
		if err != nil {
			logging.GetLogger().Error("Problem in AEAD backend generation", zap.Error(err))
		}
		DecryptedKeysetMap[k] = kh
	}
}

// GetKeyset gets a random keyset for encryption
func GetKeyset(DecryptedKeysetMap map[string]*keyset.Handle) []string {
	var KeysetArr = make([]string, len(DecryptedKeysetMap))
	for k := range DecryptedKeysetMap {
		KeysetArr = append(KeysetArr, k)
	}

	return KeysetArr
}

// SelectKeyset is used to randomly choose a keyset
func SelectKeyset(keysetArr []string) string {
	lengthOfArr := len(keysetArr)
	timeNow := time.Now().Unix()
	index := timeNow % int64(lengthOfArr)
	
	return keysetArr[int(index)]
}
