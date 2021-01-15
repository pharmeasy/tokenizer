package kms

import (
	"bitbucket.org/pharmaeasyteam/goframework/logging"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/google/tink/go/core/registry"
	"github.com/google/tink/go/integration/awskms"
	"go.uber.org/zap"
)

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
	uriClient := "aws-kms://arn:<partition>:kms:<region>:[:path]"
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
