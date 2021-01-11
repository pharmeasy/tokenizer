package jsonparser

import (
	"encoding/json"
	"net/http"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models"
)

func ParseData(req *http.Request) (string , string) {
	decoder := json.NewDecoder(req.Body)
	test := models.Param{}
	err := decoder.Decode(&test)
	if err != nil {
		panic(err)
	}
	return test.Data.Content , test.Source

}