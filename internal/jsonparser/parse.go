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

	// This is the validation check. But its not working. Will have to work on it

	// if err := req.ParseForm(); err != nil {
	// 		logging.GetLogger().Error("Problem in input params", zap.Error(err))
	// }
	// requestId := req.FormValue("requestId")
	// source := req.FormValue("source")
	// level := req.FormValue("level")
	// if requestId == "" || source == "" || level == "" {
	// 		logging.GetLogger().Error("Problem in input params")
	// 		return nil
	// }

	decoder := json.NewDecoder(req.Body)
	test := encryption.Request{}
	err := decoder.Decode(&test)
	if err != nil {
		logging.GetLogger().Error("Problem in input params", zap.Error(err))
	}
	return test.Data

}
