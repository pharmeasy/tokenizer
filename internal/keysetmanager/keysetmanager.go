package keysetmanager

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/pharmaeasyteam/goframework/logging"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/keysetmodel"
	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/core/registry"
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

func loadKeyset() (map[string]*strings.Reader, error) {

	m := make(map[string]string)
	m["ks1-1"] = "{\"encryptedKeyset\":\"AAAAqgECAgB4vGxY3UlR+bDOSItpS2xnoNeHuV+fYY0mdMOOdKKATKgBEnOuy9f/EZSEt//tlbJaHwAAAHAwbgYJKoZIhvcNAQcGoGEwXwIBADBaBgkqhkiG9w0BBwEwHgYJYIZIAWUDBAEuMBEEDGNy1Mv9AMbkI9ppdwIBEIAtgMwCo432GJ14Qg/j/cA9TtaDKzMOeUPByihWxwbYMobGrN6I4RtoR+fF8dmR/f9tHmriEKbPyTF571qPBta1gmYbWOZF8Vka1vL3+AXDdQM0wulg/eBi0N/Cy37NbOOUB2WdjQdfkp7xonDi9Upklcq4wsrfEr91CKy3zz2FfMw3IM9txrk12AlIpYNvDoDybkLjyrISfT3qIIUkidguO9QZpDno\",\"keysetInfo\":{\"primaryKeyId\":1938564032,\"keyInfo\":[{\"typeUrl\":\"type.googleapis.com/google.crypto.tink.AesGcmKey\",\"status\":\"ENABLED\",\"keyId\":1938564032,\"outputPrefixType\":\"TINK\"}]}}"
	m["ks2-1"] = "{\"encryptedKeyset\":\"AAAAqgECAgB4vGxY3UlR+bDOSItpS2xnoNeHuV+fYY0mdMOOdKKATKgBA+JtJT9l5GguSvb7iWZVVwAAAHAwbgYJKoZIhvcNAQcGoGEwXwIBADBaBgkqhkiG9w0BBwEwHgYJYIZIAWUDBAEuMBEEDL7x1AuuhfL0vSz6DQIBEIAtvpAZmAHYU0GBqYnfYg2t53TyAcY2BT1xJAJ6F/uAJf43JPsjlr0BOGQ5f21bjNmvI8VG0HWOL0bPYQoRtKMquBSufqA/HSoKCi6QJaQ3I8qjzod6qawdZ2469en5MjhFElSxt3oxnCL1xjTKtUmP0Bg3d/m4MwTr2PQmvZjn38WZ4a4WNIifHn2Xck0yEjP55n5kNvOMBOgB4mFjGdAU9fd80izF\",\"keysetInfo\":{\"primaryKeyId\":765849618,\"keyInfo\":[{\"typeUrl\":\"type.googleapis.com/google.crypto.tink.AesGcmKey\",\"status\":\"ENABLED\",\"keyId\":765849618,\"outputPrefixType\":\"TINK\"}]}}"
	m["ks3-1"] = "{\"encryptedKeyset\":\"AAAAqgECAgB4vGxY3UlR+bDOSItpS2xnoNeHuV+fYY0mdMOOdKKATKgBYMzySTxp/uTnvOz/SYubMgAAAHAwbgYJKoZIhvcNAQcGoGEwXwIBADBaBgkqhkiG9w0BBwEwHgYJYIZIAWUDBAEuMBEEDBn1Oj+zJbB6l6N2DgIBEIAtyxz0vDGBdEHXO9g6oOn6V1xzbQ8t2x2e+jcPZ1e0VXiVUlbkFdglD+xEo3fbOMREBcocaggF5toNNX84Dm1n5f0fGUH1xD8NsgM9O33kJ46niBxmw17cUoVi1jb4X5LRiKqbfGVCna4AbrTOuiAmgIgx+vQ9KgXarnFKyZV9pNMbP3wqdLpIrIGlKC45oLhHBLGbSAqYLODaiIfey5j3GBL3biwp\",\"keysetInfo\":{\"primaryKeyId\":2501366831,\"keyInfo\":[{\"typeUrl\":\"type.googleapis.com/google.crypto.tink.AesGcmKey\",\"status\":\"ENABLED\",\"keyId\":2501366831,\"outputPrefixType\":\"TINK\"}]}}"
	m["ks4-1"] = "{\"encryptedKeyset\":\"AAAAqgECAgB4vGxY3UlR+bDOSItpS2xnoNeHuV+fYY0mdMOOdKKATKgBb6/a0FiKQ6uegVNLH8N6sgAAAHAwbgYJKoZIhvcNAQcGoGEwXwIBADBaBgkqhkiG9w0BBwEwHgYJYIZIAWUDBAEuMBEEDAvd6QiQa0skdVyQvwIBEIAt0Iil9tQBX+1mpVH2dgCdRT6LwXMRKVvw3GMovClqhVMo4FBg2g5W6ylYtERDwHzKa86iA+ELPb67PAPM2yOyTuYB75owXZTHLw+d147/qMtdHw/VigUvOKU46BeRD1CaZY/JZ3SWWnUsl/1gO7aebvKKVMfpK3Ayj6RPS7XsLvu+84Z8XK2zDNVlDiqCiJokSq4OzS4TCGORTpwWLT5slBixGFqv\",\"keysetInfo\":{\"primaryKeyId\":3905263506,\"keyInfo\":[{\"typeUrl\":\"type.googleapis.com/google.crypto.tink.AesGcmKey\",\"status\":\"ENABLED\",\"keyId\":3905263506,\"outputPrefixType\":\"TINK\"}]}}"

	keysetmap := make(map[string]*strings.Reader)
	var err error

	for k, v := range m {
		keysetmap[k], err = destringify(strconv.Quote(v))
		if err != nil {
			logging.GetLogger().Error("Error encountered while destringifying the keyset file.", zap.Error(err))
			return nil, err
		}
	}

	return keysetmap, nil
}

// DecryptKeyset decrypts the keyset and stores in memory
func decryptKeyset(keysetmap map[string]*strings.Reader) error {

	awsclient, err := awskms.NewClient("aws-kms://arn:aws:kms:ap-south-1:127603365779:key/8d853831-94e6-4ac7-a0c7-3e2795e9715b")
	if err != nil {
		logging.GetLogger().Error("Error encountered in initializing KMS client.", zap.Error(err))
		return err
	}
	registry.RegisterKMSClient(awsclient)
	backend, err := awsclient.GetAEAD("aws-kms://arn:aws:kms:ap-south-1:127603365779:key/8d853831-94e6-4ac7-a0c7-3e2795e9715b")
	if err != nil {
		logging.GetLogger().Error("Error encountered in registering KMS client.", zap.Error(err))
		return err
	}
	masterKey := aead.NewKMSEnvelopeAEAD2(aead.AES128GCMKeyTemplate(), backend)

	for k, v := range keysetmap {
		kh1 := keyset.NewJSONReader(v)
		kh, err := keyset.Read(kh1, masterKey)
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
func initKeysets() error {
	if len(decryptedKeysetMap) != 0 {
		return nil
	}
	// read from source
	keysetmap, err := loadKeyset()
	if err != nil {
		return err
	}
	// decrypt & store the keysets in memory
	err = decryptKeyset(keysetmap)
	if err != nil {
		return err
	}

	return nil
}

// GetKeysetHandlerForEncryption returns a random key handler for encryption process
func GetKeysetHandlerForEncryption() (*string, *keyset.Handle, error) {
	err := initKeysets()
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
func GetKeysetHandlerForDecryption(keysetName string) (*keyset.Handle, error) {
	err := initKeysets()
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
