package jsonparser

import (
	"encoding/json"
	"net/http"

	"bitbucket.org/pharmaeasyteam/goframework/logging"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models"
	"go.uber.org/zap"
	//"strconv"
)

// ParseData ...
func ParseData(req *http.Request) (string, string) {
	decoder := json.NewDecoder(req.Body)

	test := models.Param{}
	//err := decoder.DisallowUnknownFields()
	err := decoder.Decode(&test)
	if err != nil {
		logging.GetLogger().Error("Problem in input params", zap.Error(err))
	}
	return test.Data.Content, test.Source

}

// CheckAllParams ...
func CheckAllParams(req *http.Request) {
	query := req.URL.Query().Get("source")
	l
}
