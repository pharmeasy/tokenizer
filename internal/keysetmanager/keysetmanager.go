package keysetmanager

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/google/tink/go/integration/awskms"
	"github.com/google/tink/go/keyset"
	"github.com/pharmaeasy/tokenizer/internal/errormanager"
	keysetmodel "github.com/pharmaeasy/tokenizer/internal/models/keyset"
)

// DecryptedKeysetMap stores the keys in memory
var decryptedKeysetMap = make(map[string]*keyset.Handle)
var keysetMapFromEnv = make(map[string]string)
var kmsARN string

func destringify(str string) (*strings.Reader, error) {
	var keyset keysetmodel.EncryptedKeyset
	s, err := strconv.Unquote(str)
	if err != nil {
		return nil, errormanager.SetError("Error encountered while destringifying keysets.", err)
	}
	err = json.Unmarshal([]byte(s), &keyset)
	if err != nil {
		return nil, errormanager.SetError("Error encountered while unmarshalling keysets.", err)
	}
	JSONKeyset, _ := json.Marshal(keyset)
	myReader := strings.NewReader(string(JSONKeyset))

	return myReader, nil
}

func loadKeyset() (map[string]*strings.Reader, error) {
	keysetmap := make(map[string]*strings.Reader)
	var err error

	for k, v := range keysetMapFromEnv {
		keysetmap[k], err = destringify(strconv.Quote(v))
		if err != nil {
			return nil, err
		}
	}

	return keysetmap, nil
}

// DecryptKeyset decrypts the keyset and stores in memory
func decryptKeyset(keysetmap map[string]*strings.Reader) error {
  kmsARN: ="aws-kms://arn:aws:kms:ap-south-1:820116237501:key/feb4b915-d341-4c29-9bdd-b968a76cabe3";
	kmsClient, err := awskms.NewClient(kmsARN)
	if err != nil {
		return errormanager.SetError("Error encountered in initializing KMS client.", err)
	}
	kmsAEAD, err := kmsClient.GetAEAD(kmsARN)
	if err != nil {
		return errormanager.SetError("Error encountered in initializing KMS AEAD client.", err)
	}

	for k, v := range keysetmap {
		kh1 := keyset.NewJSONReader(v)
		kh, err := keyset.Read(kh1, kmsAEAD)
		if err != nil {
			return errormanager.SetError("Error encountered in reading the AEAD keyset.", err)
		}
		decryptedKeysetMap[k] = kh
	}

	return nil
}

func getRandomizedKeyset() (*string, *keyset.Handle, error) {
	if len(decryptedKeysetMap) == 0 {
		err := errormanager.SetError("Error encountered in reading the keyset.", errors.New("the keyset map is empty"))
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
	err = errormanager.SetError("Error encountered in fetching the keyset.", errors.New("Valid keyset not found for keyset name"+keysetName))

	return nil, err
}

// LoadKeysetFromConfig loads keyset from env
func LoadKeysetFromConfig(keyMap map[string]string) {
	keysetMapFromEnv = keyMap
}

// LoadArnFromConfig loads kms arn from env
func LoadArnFromConfig(str string) {
	kmsARN = str
}
