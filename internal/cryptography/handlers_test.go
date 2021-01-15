package cryptography

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const cryptographySuccessResponse = `
	"token" : fe6fe952-885f-483c-9a41-4278d2c4c63c
`

const cryptographyEmptyResponse = `
	"token" : 
`

func Test_DataEncrypt(t *testing.T) {

	x, _ := DataEncrypt("9422401444", "private_key")
	assert.Equal(t, x, �I�����g$=�PV�D�t���n�8%��joΆƲ�{�h))
	
}
