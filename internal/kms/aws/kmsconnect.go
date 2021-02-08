package kms

import (
	//"bytes"
	//"fmt"
	//"io/ioutil"

	"math/rand"
	"os"
	"time"

	"bitbucket.org/pharmaeasyteam/goframework/logging"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/core/registry"

	//"github.com/google/tink/go/insecurecleartextkeyset"
	"github.com/google/tink/go/integration/awskms"
	"github.com/google/tink/go/keyset"

	"go.uber.org/zap"
)

var datakeyMap = map[string][]string{
	"arn:<partition>:kms:<region>:[:path1]": []string{"data_key1", "data_key2", "data_key3"},
	"arn:<partition>:kms:<region>:[:path2]": []string{"data_key1", "data_key2", "data_key3"},
}

var arnMap = map[int]string{
	0: "arn:<partition>:kms:<region>:[:path1]",
	1: "arn:<partition>:kms:<region>:[:path2]",
}

// CreateInstance ...
func CreateInstance() *kms.KMS {
	mySession, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1")},
	)

	if err != nil {
		logging.GetLogger().Error("Problem in session ", zap.Error(err))
	}

	svc := kms.New(mySession)
	return svc
}

// CreateCmkKey ...
func CreateCmkKey(svc *kms.KMS) *string {
	result, err := svc.CreateKey(&kms.CreateKeyInput{
		Tags: []*kms.Tag{
			{
				TagKey:   aws.String("CreatedBy"),
				TagValue: aws.String("ExampleUser"),
			},
		},
	})
	if err != nil {
		logging.GetLogger().Error("Problem in keyset generation", zap.Error(err))
	}

	return result.KeyMetadata.KeyId
}

// KmsClient ...
func KmsClient() (registry.KMSClient, error) {
	uriClient := "aws-kms://arn:aws:kms:ap-south-1:127603365779:key/8d853831-94e6-4ac7-a0c7-3e2795e9715b"
	return awskms.NewClient(uriClient)
}

// DataEncrypt ...
func DataEncrypt(data string, keyURI string) []byte {
	kmsClient, err := KmsClient()
	if err != nil {
		logging.GetLogger().Error("Problem in kms client generation", zap.Error(err))
	}
	//uriClient := "aws-kms://arn:<partition>:kms:<region>:[:path]"
	backend, err := kmsClient.GetAEAD(keyURI)
	if err != nil {
		logging.GetLogger().Error("Problem in AEAD backend generation", zap.Error(err))
	}
	ct, err := backend.Encrypt([]byte(data), []byte(keyURI))
	return ct
}

// DataDecrypt ...
func DataDecrypt(ct []byte, keyURI string) []byte {
	kmsClient, err := KmsClient()
	if err != nil {
		logging.GetLogger().Error("Problem in kms client generation", zap.Error(err))
	}
	//uriClient := "aws-kms://arn:<partition>:kms:<region>:[:path]"
	backend, err := kmsClient.GetAEAD(keyURI)
	if err != nil {
		logging.GetLogger().Error("Problem in AEAD backend generation", zap.Error(err))
	}
	pt, err := backend.Decrypt(ct, []byte(keyURI))
	return pt
}

// getARN returns the data key used for encryption
func getARN(dataKeyMap map[string][]string, arnMap map[int]string) string {
	timeNow := time.Now().Unix()
	partition := timeNow % int64(len(arnMap))
	keyArr := arnMap[int(partition)]
	datakeyArr := dataKeyMap[keyArr]
	index := rand.Intn(len(datakeyArr))
	return datakeyArr[index]
}

// func GetKeysets() (*keyset.Handle, *keyset.MemReaderWriter) {
// 	awsclient, err := awskms.NewClient("aws-kms://arn:aws:kms:ap-south-1:127603365779:key/8d853831-94e6-4ac7-a0c7-3e2795e9715b")
// 	if err != nil {
// 		logging.GetLogger().Error("Problem in AEAD backend generation", zap.Error(err))
// 	}
// 	registry.RegisterKMSClient(awsclient)
// 	backend, err := awsclient.GetAEAD("aws-kms://arn:aws:kms:ap-south-1:127603365779:key/8d853831-94e6-4ac7-a0c7-3e2795e9715b")
// 	if err != nil {
// 		logging.GetLogger().Error("Problem in AEAD backend generation", zap.Error(err))
// 	}

// 	kh1, err := keyset.NewHandle(aead.AES128GCMKeyTemplate())
// 	if err != nil {
// 		logging.GetLogger().Error("Problem in AEAD backend generation", zap.Error(err))
// 	}
// 	masterKey := aead.NewKMSEnvelopeAEAD2(aead.AES128GCMKeyTemplate(), backend)
// 	//file, _ := os.Open("/Users/riddhimanparasar/e-keyset.json")

// 	// memKeyset := keyset.NewJSONReader(file)
// 	// b, _ := ioutil.ReadAll(file)
// 	// fmt.Println(string(b))

// 	memKeyset := &keyset.MemReaderWriter{}

// 	if err := kh1.Write(memKeyset, masterKey); err != nil {
// 		logging.GetLogger().Error("Problem in AEAD backend generation", zap.Error(err))
// 	}
// 	kh, err := keyset.Read(memKeyset, masterKey)
// 	if err != nil {
// 		logging.GetLogger().Error("Problem in AEAD backend generation", zap.Error(err))
// 	}
// 	return kh, memKeyset
// }

func GetKeysets(filename string) *keyset.Handle {
	awsclient, err := awskms.NewClient("aws-kms://arn:aws:kms:ap-south-1:127603365779:key/8d853831-94e6-4ac7-a0c7-3e2795e9715b")
	if err != nil {
		logging.GetLogger().Error("Problem in AEAD backend generation", zap.Error(err))
	}
	registry.RegisterKMSClient(awsclient)
	backend, err := awsclient.GetAEAD("aws-kms://arn:aws:kms:ap-south-1:127603365779:key/8d853831-94e6-4ac7-a0c7-3e2795e9715b")
	if err != nil {
		logging.GetLogger().Error("Problem in AEAD backend generation", zap.Error(err))
	}

	file, _ := os.Open(filename)
	kh1 := keyset.NewJSONReader(file)
	masterKey := aead.NewKMSEnvelopeAEAD2(aead.AES128GCMKeyTemplate(), backend)
	kh, err := keyset.Read(kh1, masterKey)
	if err != nil {
		logging.GetLogger().Error("Problem in AEAD backend generation", zap.Error(err))
	}
	return kh

}

// func GetKeysets() *keyset.Handle {
// 	awsclient, err := awskms.NewClient("aws-kms://arn:aws:kms:ap-south-1:127603365779:key/8d853831-94e6-4ac7-a0c7-3e2795e9715b")
// 	if err != nil {
// 		logging.GetLogger().Error("Problem in AEAD backend generation", zap.Error(err))
// 	}
// 	registry.RegisterKMSClient(awsclient)
// 	// backend, err := awsclient.GetAEAD("aws-kms://arn:aws:kms:ap-south-1:127603365779:key/8d853831-94e6-4ac7-a0c7-3e2795e9715b")
// 	// if err != nil {
// 	// 	logging.GetLogger().Error("Problem in AEAD backend generation", zap.Error(err))
// 	// }
// 	// masterKey := aead.NewKMSEnvelopeAEAD2(aead.AES128GCMKeyTemplate(), backend)
// 	file, _ := os.Open("/Users/riddhimanparasar/keyset.cfg")
// 	// b, _ := ioutil.ReadFile("/Users/riddhimanparasar/private-keyset.cfg")
// 	// fmt.Print(string(b))
// 	// rdr := keyset.NewBinaryReader(bytes.NewReader(b))
// 	reader := keyset.NewBinaryReader(file)
// 	handle, err := testkeyset.Read(reader)
// 	if err != nil {
// 		logging.GetLogger().Error("Problem in AEAD backend generation", zap.Error(err))
// 	}
// 	// return kh.KeysetInfo().String()
// 	return handle
// }

// func GetKeysets() *keyset.Handle {
// 	keyURI := "aws-kms://arn:aws:kms:ap-south-1:127603365779:key/8d853831-94e6-4ac7-a0c7-3e2795e9715b"
// 	// Generate the keyset and storing in JSON
// 	kh1, err := keyset.NewHandle(aead.AES128GCMKeyTemplate())
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	buf, err := os.Create("./encrypted-keyset.json")
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	defer buf.Close()
// 	w := keyset.NewJSONWriter(buf)
// 	//added buffer to jsonwriter

// 	// Fetch the master key from a KMS.
// 	gcpClient, err := awskms.NewClient(keyURI)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	registry.RegisterKMSClient(gcpClient)

// 	backend, err := gcpClient.GetAEAD(keyURI)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	masterKey := aead.NewKMSEnvelopeAEAD2(aead.AES128GCMKeyTemplate(), backend)

// 	if err := kh1.Write(w, masterKey); err != nil {
// 		log.Fatal(err)
// 	}

// 	// Reading the Keyset

// 	jsonKeyset, err := os.Open("./encrypted-keyset.json") // For read access.
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	r := keyset.NewJSONReader(jsonKeyset)

// 	kh2, err := keyset.Read(r, masterKey)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if kh1.String() != kh2.String() {
// 		log.Println("key handlers are not equal")
// 	}
// 	return kh2

// }

// func GetKeysets() *keyset.Handle {
// 	keyURI := "aws-kms://arn:aws:kms:ap-south-1:127603365779:key/8d853831-94e6-4ac7-a0c7-3e2795e9715b"
// 	// Generate the keyset and storing in JSON
// 	// kh1, err := keyset.NewHandle(aead.AES128GCMKeyTemplate())
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// buf, err := os.Create("./encrypted-keyset.json")
// 	// if err != nil {
// 	// 	log.Println(err)
// 	// }
// 	// defer buf.Close()
// 	// w := keyset.NewJSONWriter(buf)
// 	//added buffer to jsonwriter

// 	// Fetch the master key from a KMS.
// 	gcpClient, err := awskms.NewClient(keyURI)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	registry.RegisterKMSClient(gcpClient)

// 	backend, err := gcpClient.GetAEAD(keyURI)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	masterKey := aead.NewKMSEnvelopeAEAD2(aead.AES128GCMKeyTemplate(), backend)

// 	// if err := kh1.Write(w, masterKey); err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// Reading the Keyset

// 	jsonKeyset, err := os.Open("./encrypted-keyset.json") // For read access.
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	r := keyset.NewJSONReader(jsonKeyset)

// 	kh2, err := keyset.Read(r, masterKey)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// if kh1.String() != kh2.String() {
// 	// 	log.Println("key handlers are not equal")
// 	// }
// 	return kh2

// }
