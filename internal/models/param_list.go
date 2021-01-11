package models

type Param struct {
	RequestID	int 	`json:"requestId"`
	Source		string		`json:"source"`
	Data		Data 		`json:"data"`
}

type Data struct {
	AttributeType	string	`json:"attributeType"`
	Content			string	`json:"content"`
}

