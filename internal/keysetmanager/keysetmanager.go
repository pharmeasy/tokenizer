package keysetmanager

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/pharmaeasyteam/goframework/logging"
	keysetmodel "bitbucket.org/pharmaeasyteam/tokenizer/internal/models/keyset"
	"github.com/google/tink/go/integration/awskms"
	"github.com/google/tink/go/keyset"
	"go.uber.org/zap"
)

// DecryptedKeysetMap stores the keys in memory
var decryptedKeysetMap = make(map[string]*keyset.Handle)

func destringify(str string) (*strings.Reader, error) {
	var keyset keysetmodel.EncryptedKeyset
	s, err := strconv.Unquote(str)
	if err != nil {
		logging.GetLogger().Error("Error encountered in destrinify", zap.Error(err))
		return nil, err
	}
	err = json.Unmarshal([]byte(s), &keyset)
	if err != nil {
		logging.GetLogger().Error("Error encountered in Unmarshal", zap.Error(err))
		return nil, err
	}
	JSONKeyset, _ := json.Marshal(keyset)
	myReader := strings.NewReader(string(JSONKeyset))

	return myReader, nil
}

func loadKeyset(keymap map[string]string) (map[string]*strings.Reader, error) {

	// m := make(map[string]string)
	// m["ks1-1"] = "{ \"keysetInfo\": { \"primaryKeyId\": 1747494060, \"keyInfo\": [{ \"typeUrl\": \"type.googleapis.com/google.crypto.tink.AesGcmKey\", \"outputPrefixType\": \"TINK\", \"keyId\": 1747494060, \"status\": \"ENABLED\" }] }, \"encryptedKeyset\": \"AQICAHi8bFjdSVH5sM5Ii2lLbGeg14e5X59hjSZ0w450ooBMqAHIQkLlr629X7IiulKowkscAAAAzzCBzAYJKoZIhvcNAQcGoIG+MIG7AgEAMIG1BgkqhkiG9w0BBwEwHgYJYIZIAWUDBAEuMBEEDBGL1Y5zAn4FDnnWkgIBEICBh47nZmrjziKp9KwtDfQHkGqq1EX+tejM+ZxPdfbWe5xbjFgx+RebqOGCz34j4ek/QuNJNOjIFc/+eiK0IVn6d657uA4Km2VKOpCxrIaWqkAXVB7E22vCg23iIuZsfYiyLzOSD252PRJwE4L/TlpeFHNF4PmBH/Go5+tfhZj/WSDxCavqQQUgMw==\" }"
	// m["ks2-1"] = "{ \"keysetInfo\": { \"primaryKeyId\": 1747494060, \"keyInfo\": [{ \"typeUrl\": \"type.googleapis.com/google.crypto.tink.AesGcmKey\", \"outputPrefixType\": \"TINK\", \"keyId\": 1747494060, \"status\": \"ENABLED\" }] }, \"encryptedKeyset\": \"AQICAHi8bFjdSVH5sM5Ii2lLbGeg14e5X59hjSZ0w450ooBMqAHIQkLlr629X7IiulKowkscAAAAzzCBzAYJKoZIhvcNAQcGoIG+MIG7AgEAMIG1BgkqhkiG9w0BBwEwHgYJYIZIAWUDBAEuMBEEDBGL1Y5zAn4FDnnWkgIBEICBh47nZmrjziKp9KwtDfQHkGqq1EX+tejM+ZxPdfbWe5xbjFgx+RebqOGCz34j4ek/QuNJNOjIFc/+eiK0IVn6d657uA4Km2VKOpCxrIaWqkAXVB7E22vCg23iIuZsfYiyLzOSD252PRJwE4L/TlpeFHNF4PmBH/Go5+tfhZj/WSDxCavqQQUgMw==\" }"
	// m["ks3-1"] = "{ \"keysetInfo\": { \"primaryKeyId\": 1747494060, \"keyInfo\": [{ \"typeUrl\": \"type.googleapis.com/google.crypto.tink.AesGcmKey\", \"outputPrefixType\": \"TINK\", \"keyId\": 1747494060, \"status\": \"ENABLED\" }] }, \"encryptedKeyset\": \"AQICAHi8bFjdSVH5sM5Ii2lLbGeg14e5X59hjSZ0w450ooBMqAHIQkLlr629X7IiulKowkscAAAAzzCBzAYJKoZIhvcNAQcGoIG+MIG7AgEAMIG1BgkqhkiG9w0BBwEwHgYJYIZIAWUDBAEuMBEEDBGL1Y5zAn4FDnnWkgIBEICBh47nZmrjziKp9KwtDfQHkGqq1EX+tejM+ZxPdfbWe5xbjFgx+RebqOGCz34j4ek/QuNJNOjIFc/+eiK0IVn6d657uA4Km2VKOpCxrIaWqkAXVB7E22vCg23iIuZsfYiyLzOSD252PRJwE4L/TlpeFHNF4PmBH/Go5+tfhZj/WSDxCavqQQUgMw==\" }"
	// m["ks4-1"] = "{ \"keysetInfo\": { \"primaryKeyId\": 1747494060, \"keyInfo\": [{ \"typeUrl\": \"type.googleapis.com/google.crypto.tink.AesGcmKey\", \"outputPrefixType\": \"TINK\", \"keyId\": 1747494060, \"status\": \"ENABLED\" }] }, \"encryptedKeyset\": \"AQICAHi8bFjdSVH5sM5Ii2lLbGeg14e5X59hjSZ0w450ooBMqAHIQkLlr629X7IiulKowkscAAAAzzCBzAYJKoZIhvcNAQcGoIG+MIG7AgEAMIG1BgkqhkiG9w0BBwEwHgYJYIZIAWUDBAEuMBEEDBGL1Y5zAn4FDnnWkgIBEICBh47nZmrjziKp9KwtDfQHkGqq1EX+tejM+ZxPdfbWe5xbjFgx+RebqOGCz34j4ek/QuNJNOjIFc/+eiK0IVn6d657uA4Km2VKOpCxrIaWqkAXVB7E22vCg23iIuZsfYiyLzOSD252PRJwE4L/TlpeFHNF4PmBH/Go5+tfhZj/WSDxCavqQQUgMw==\" }"

	keysetmap := make(map[string]*strings.Reader)
	var err error

	for k, v := range keymap {
		keysetmap[k], err = destringify(strconv.Quote(v))
		if err != nil {
			logging.GetLogger().Error("Error encountered while destringifying the keyset file.", zap.Error(err))
			return nil, err
		}
	}

	return keysetmap, nil
}

// DecryptKeyset decrypts the keyset and stores in memory
func decryptKeyset(keysetmap map[string]*strings.Reader, keyURI string) error {

	//keyURI := "aws-kms://arn:aws:kms:ap-south-1:127603365779:key/8d853831-94e6-4ac7-a0c7-3e2795e9715b"

	kmsClient, err := awskms.NewClient(keyURI)
	if err != nil {
		logging.GetLogger().Error("Error encountered in initializing KMS client.", zap.Error(err))
		return err
	}
	kmsAEAD, err := kmsClient.GetAEAD(keyURI)
	if err != nil {
		logging.GetLogger().Error("Error encountered in initializing KMS client.", zap.Error(err))
		return err
	}

	for k, v := range keysetmap {
		kh1 := keyset.NewJSONReader(v)
		kh, err := keyset.Read(kh1, kmsAEAD)
		if err != nil {
			logging.GetLogger().Error("Error encountered in reading the keyset.", zap.Error(err))
			return err
		}
		decryptedKeysetMap[k] = kh
	}

	return nil
}

func getRandomizedKeyset() (*string, *keyset.Handle, error) {
	if len(decryptedKeysetMap) == 0 {
		err := errors.New("the keyset map is empty")
		logging.GetLogger().Error("Error encountered in reading the keyset.", zap.Error(err))
		return nil, nil, err
	}

	var keysetArr = make([]string, len(decryptedKeysetMap))

	i := 0
	for k := range decryptedKeysetMap {
		keysetArr[i] = k
		i++
	}

	lengthOfArr := len(keysetArr)
	timeNow := time.Now().Unix()
	index := timeNow % int64(lengthOfArr)

	keyName := keysetArr[int(index)]
	keyHandle := decryptedKeysetMap[keysetArr[int(index)]]

	return &keyName, keyHandle, nil
}

// InitKeysets initializes the keysets
func initKeysets(keymap map[string]string, keyURI string) error {
	if len(decryptedKeysetMap) != 0 {
		return nil
	}
	// read from source
	keysetmap, err := loadKeyset(keymap)
	if err != nil {
		return err
	}
	// decrypt & store the keysets in memory
	err = decryptKeyset(keysetmap, keyURI)
	if err != nil {
		return err
	}

	return nil
}

// GetKeysetHandlerForEncryption returns a random key handler for encryption process
func GetKeysetHandlerForEncryption(keymap map[string]string, keyURI string) (*string, *keyset.Handle, error) {
	err := initKeysets(keymap, keyURI)
	if err != nil {
		return nil, nil, err
	}

	keyName, keyHandle, err := getRandomizedKeyset()
	if err != nil {
		return nil, nil, err
	}

	return keyName, keyHandle, nil
}

// GetKeysetHandlerForDecryption returns a random key handler for decryption process
func GetKeysetHandlerForDecryption(keysetName string, keymap map[string]string, keyURI string) (*keyset.Handle, error) {
	err := initKeysets(keymap, keyURI)
	if err != nil {
		return nil, err
	}

	if kh, ok := decryptedKeysetMap[keysetName]; ok {
		return kh, nil
	}
	err = errors.New("Something went wrong while processing your request")
	logging.GetLogger().Error("Valid keyset not found for keyset name "+keysetName, zap.Error(err))

	return nil, err
}
