package keysetmanager

import (
	"encoding/json"
	"errors"
	"io/ioutil"
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
	fileName, err := ioutil.ReadFile("/Users/riddhimanparasar/tokenizer/keysetmap.json")
	if err != nil {
		logging.GetLogger().Error("Error encountered while reading the keyset file.", zap.Error(err))
		return nil, err
	}

	m := make(map[string]string)
	err = json.Unmarshal(fileName, &m)
	if err != nil {
		logging.GetLogger().Error("Error encountered while unmarshalling the keyset file.", zap.Error(err))
		return nil, err
	}
	keysetmap := make(map[string]*strings.Reader)

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
