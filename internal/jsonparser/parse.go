package jsonparser

import (
	"encoding/json"
	"net/http"

	"bitbucket.org/pharmaeasyteam/goframework/logging"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/request/encryption"
	"go.uber.org/zap"
	//"strconv"
)

// ParseData ...
func ParseData(req *http.Request) []encryption.Data {

	decoder := json.NewDecoder(req.Body)
	test := encryption.Request{}
	err := decoder.Decode(&test)
	if err != nil {
		logging.GetLogger().Error("Problem in input params", zap.Error(err))
	}

	// input validations
	requestId := test.RequestID
	source := test.Source
	level := test.Level

	if requestId < 1 || source == "" || level < 1 {
		logging.GetLogger().Error("Problem in input params", zap.Error(err))
		return nil
	}
	return test.Data

}
